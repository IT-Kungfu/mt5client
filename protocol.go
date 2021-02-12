package mt5client

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"math/rand"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	MT5HeaderLength                = 9
	MT5PacketSeparator             = "\r\n"
	MT5CommandSeparator            = "|"
	MT5ParamSeparator              = "="
	MT5RetCode                     = "RETCODE"
	MT5RetCodeSuccess              = "0 Done"
	MT5CommandAuthStart            = "AUTH_START"
	MT5CommandAuthAnswer           = "AUTH_ANSWER"
	MT5CommandQuit                 = "QUIT"
	MT5CommandOrderGetTotal        = "ORDER_GET_TOTAL"
	MT5CommandOrderGetPage         = "ORDER_GET_PAGE"
	MT5CommandOrderGetBatch        = "ORDER_GET_BATCH"
	MT5CommandOrderGetHistoryTotal = "HISTORY_GET_TOTAL"
	MT5CommandOrderGetHistoryPage  = "HISTORY_GET_PAGE"
	MT5CommandOrderGetHistoryBatch = "HISTORY_GET_BATCH"
	MT5CommandDealGetTotal         = "DEAL_GET_TOTAL"
	MT5CommandDealGetPage          = "DEAL_GET_PAGE"
	MT5CommandDealGetBatch         = "DEAL_GET_BATCH"
	MT5CommandDealDelete           = "DEAL_DELETE"
	MT5CommandPositionGetTotal     = "POSITION_GET_TOTAL"
	MT5CommandPositionGetPage      = "POSITION_GET_PAGE"
	MT5CommandPositionGetBatch     = "POSITION_GET_BATCH"
	MT5CommandPositionDelete       = "POSITION_DELETE"
	MT5CommandClientGetIds         = "CLIENT_IDS"
	MT5CommandUserGet              = "USER_GET"
	MT5CommandUserGetBatch         = "USER_GET_BATCH"
	MT5CommandUserAdd              = "USER_ADD"
	MT5CommandUserUpdate           = "USER_UPDATE"
	MT5CommandUserDelete           = "USER_DELETE"
	MT5CommandUserAccountGetBatch  = "USER_ACCOUNT_GET_BATCH"
	MT5CommandTradeBalance         = "TRADE_BALANCE"
	MT5CommandTickGetHistory       = "TICK_HISTORY_GET"
	MT5CommandChartGet             = "CHART_GET"
	MT5CommandDealerSend           = "DEALER_SEND"
	MT5CommandDealerUpdates        = "DEALER_UPDATES"
)

type MT5Header struct {
	bodyLen      int
	packetNumber int
	flag         uint8
}

type MT5Command struct {
	Name    string
	Params  map[string]string
	Payload string
}

var (
	utf16 = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
)

func (c *MT5Client) makeRequest(cmd *MT5Command) ([]byte, string, error) {
	body := cmd.Name + MT5CommandSeparator
	for k, v := range cmd.Params {
		if v != "" {
			body += k + MT5ParamSeparator + v + MT5CommandSeparator
		}
	}
	body += MT5PacketSeparator
	body += cmd.Payload
	req, err := makePacket(body, 0, 0)
	return req, body, err
}

func makePacket(body string, packetNumber uint16, flag uint8) ([]byte, error) {
	encBody, err := utf16.NewEncoder().String(body)
	if err != nil {
		return nil, err
	}

	p := []byte(fmt.Sprintf("%04x", len(encBody)))
	p = append(p, []byte(fmt.Sprintf("%04x", packetNumber))...)
	p = append(p, []byte(fmt.Sprintf("%x", flag))...)
	p = append(p, encBody...)

	return p, nil
}

func (c *MT5Client) reconnect() error {
	c.connMux.Lock()
	defer c.connMux.Unlock()

	_ = c.conn.Close()
	c.conn = nil

	for {
		if err := c.init(); err != nil {
			c.log.Errorf("Reconnect error %v", err)
			time.Sleep(time.Second)
		} else {
			c.log.Infof("Reconnect successfully")
			break
		}
	}

	return c.auth()
}

func (c *MT5Client) sendRequest(request []byte) (*MT5Command, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("no connection")
	}
	_, err := c.conn.Write(request)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			c.log.Errorf("#%d write failed, trying reconnect %v", c.clientId, err)
			_ = c.reconnect()
			return nil, err
		}
		return nil, fmt.Errorf("#%d write request failed %v", c.clientId, err)
	}

	buffer := make([]byte, 0)
	body := make([]byte, 0)

	for {
		header, err := c.readHeader()
		if err != nil {
			c.log.Errorf("#%d read header failed %v", c.clientId, err)
			_ = c.reconnect()
			return nil, err
		}

		//c.log.Debugf("Response header %+v", header)

		if header.bodyLen == 0 {
			c.log.Debugf("#%d PING packet. Header: %+v", c.clientId, header)
			continue
		}

		if header.flag == 0x01 {
			chunk, err := c.readBody(header.bodyLen, []byte{})
			if err != nil {
				return nil, err
			}
			buffer = append(buffer, chunk...)
		} else {
			body, err = c.readBody(header.bodyLen, buffer)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return parseBody(body)
}

func (c *MT5Client) readBody(size int, appendBuffer []byte) ([]byte, error) {
	buffer := make([]byte, 0)
	for len(buffer) < size {
		packet := make([]byte, size-len(buffer))
		n, err := c.conn.Read(packet)
		if err != nil {
			return nil, fmt.Errorf("read body failed: %v", err)
		}
		buffer = append(buffer, packet[:n]...)
	}
	return append(appendBuffer, buffer...), nil
}

func (c *MT5Client) readHeader() (*MT5Header, error) {
	buffer := make([]byte, 0)
	for len(buffer) < MT5HeaderLength {
		packet := make([]byte, MT5HeaderLength-len(buffer))
		n, err := c.conn.Read(packet)
		if err != nil {
			return nil, fmt.Errorf("read header failed: %v", err)
		}
		buffer = append(buffer, packet[:n]...)
	}
	return parseHeader(buffer)
}

func parseHeader(header []byte) (*MT5Header, error) {
	bodyLen, err := strconv.ParseInt(string(header[:4]), 16, 32)
	if err != nil || bodyLen > 0xffff {
		return nil, fmt.Errorf("wrong body length: %v (%s)", err, header[:4])
	}

	packetNumber, err := strconv.ParseInt(string(header[4:8]), 16, 32)
	if err != nil || packetNumber > 0xffff {
		return nil, fmt.Errorf("wrong packet number: %v (%s)", err, header[4:8])
	}

	flag, err := strconv.ParseInt(string(header[8]), 16, 8)
	if err != nil || flag > 0xf {
		return nil, fmt.Errorf("wrong flag: %v (%v)", err, header[8])
	}

	return &MT5Header{bodyLen: int(bodyLen), packetNumber: int(packetNumber), flag: uint8(flag)}, nil
}

func parseBody(b []byte) (*MT5Command, error) {
	body, err := utf16.NewDecoder().Bytes(b)
	if err != nil {
		return nil, fmt.Errorf("parse body error: %v", err)
	}

	cmd := &MT5Command{
		Params: make(map[string]string, 5),
	}

	pIdx := strings.Index(string(body), MT5PacketSeparator)
	if pIdx == -1 {
		return nil, errors.New("error parsing body")
	}
	command := string(body[:pIdx])
	cmd.Payload = string(body[pIdx+len(MT5PacketSeparator):])

	for i, c := range strings.Split(command, MT5CommandSeparator) {
		if i == 0 {
			cmd.Name = c
		} else {
			p := strings.Split(c, MT5ParamSeparator)
			if len(p) == 2 {
				cmd.Params[p[0]] = p[1]
			}
		}
	}

	return cmd, nil
}

func makeRandomString() []byte {
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, 16)
	rand.Read(str)
	return str
}
