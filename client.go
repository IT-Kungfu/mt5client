package mt5client

import (
	"context"
	"fmt"
	"github.com/IT-Kungfu/logger"
	"net"
	"sync"
	"time"
)

type MT5Client struct {
	cfg       *Config
	log       *logger.Logger
	conn      *net.TCPConn
	connMux   sync.Mutex
	controlCh chan *ClientControlMessage
	clientId  int
}

func NewMT5Client(ctx context.Context) (*MT5Client, error) {
	services := ctx.Value("services").(map[string]interface{})
	c := &MT5Client{
		cfg:       services["mt5cfg"].(*Config),
		log:       services["log"].(*logger.Logger),
		controlCh: services["controlCh"].(chan *ClientControlMessage),
		clientId:  services["clientId"].(int),
	}

	err := c.init()
	if err != nil {
		return nil, err
	}

	err = c.auth()
	if err != nil {
		return nil, err
	}

	go c.ping()
	go c.loop()

	return c, nil
}

func (c *MT5Client) loop() {
	for {
		select {
		case m := <-c.controlCh:
			if c.commandHandler(m) {
				return
			}
		case <-time.After(time.Duration(c.cfg.MT5PingTimeout) * time.Second):
		}
	}
}

func (c *MT5Client) commandHandler(m *ClientControlMessage) bool {
	switch m.Cmd.Name {
	case MT5CommandQuit:
		c.quit(m)
		return true
	case MT5CommandOrderGetTotal:
		c.getOrdersTotal(m)
	case MT5CommandOrderGetHistoryTotal:
		c.getOrdersTotal(m)
	case MT5CommandOrderGetPage:
		c.getOrders(m)
	case MT5CommandOrderGetBatch:
		c.getOrders(m)
	case MT5CommandOrderGetHistoryPage:
		c.getOrders(m)
	case MT5CommandOrderGetHistoryBatch:
		c.getOrders(m)
	case MT5CommandDealGetTotal:
		c.getDealsTotal(m)
	case MT5CommandDealGetPage:
		c.getDeals(m)
	case MT5CommandDealGetBatch:
		c.getDeals(m)
	case MT5CommandDealDelete:
		c.deleteDeals(m)
	case MT5CommandPositionGetTotal:
		c.getPositionsTotal(m)
	case MT5CommandPositionGetPage:
		c.getPositionsTotal(m)
	case MT5CommandPositionGetBatch:
		c.getPositions(m)
	case MT5CommandPositionDelete:
		c.deletePositions(m)
	case MT5CommandClientGetIds:
		c.getClients(m)
	case MT5CommandUserGet:
		c.getUser(m)
	case MT5CommandUserGetBatch:
		c.getUsersBatch(m)
	case MT5CommandUserAdd:
		c.addUpdateUser(m)
	case MT5CommandUserUpdate:
		c.addUpdateUser(m)
	case MT5CommandUserDelete:
		c.deleteUser(m)
	case MT5CommandUserAccountGetBatch:
		c.getUserAccounts(m)
	case MT5CommandTradeBalance:
		c.balance(m)
	case MT5CommandTickGetHistory:
		c.getTickHistory(m)
	case MT5CommandChartGet:
		c.getChart(m)
	case MT5CommandDealerSend:
		c.sendDealer(m)
	}
	return false
}

func (c *MT5Client) init() error {
	c.log.Infof("Connecting to MT5 server %s:%d", c.cfg.MT5Host, c.cfg.MT5Port)

	mt5Addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.cfg.MT5Host, c.cfg.MT5Port))
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp", nil, mt5Addr)
	if err != nil {
		return err
	}

	_, err = c.conn.Write([]byte("MT5WEBAPI"))
	if err != nil {
		return fmt.Errorf("error init MT5WEBAPI %v", err)
	}

	return nil
}

func (c *MT5Client) ping() {
	var err error
	for {
		time.Sleep(time.Duration(c.cfg.MT5PingTimeout) * time.Second)
		request, _ := makePacket("", 0, 0)
		c.connMux.Lock()
		if c.conn != nil {
			_, err = c.conn.Write(request)
		}
		c.connMux.Unlock()
		if err != nil {
			c.log.Errorf("#%d ping request failed %v", c.clientId, err)
			_ = c.reconnect()
		}
	}
}

func safeSend(ch chan *ClientResponse, value *ClientResponse) (success bool) {
	defer func() {
		if recover() != nil {
			success = false
		}
	}()
	ch <- value
	return true
}
