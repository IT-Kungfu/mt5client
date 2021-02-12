package mt5client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (p *Pool) GetPositionsTotal(login string) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandPositionGetTotal,
			Params: map[string]string{
				"LOGIN": login,
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetPositionsPage(login string, offset, total int) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandPositionGetPage,
			Params: map[string]string{
				"LOGIN":  login,
				"OFFSET": fmt.Sprintf("%d", offset),
				"TOTAL":  fmt.Sprintf("%d", total),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetPositionsBatch(login, group, ticket []string, symbol string) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandPositionGetBatch,
			Params: map[string]string{
				"LOGIN":  strings.Join(login, ","),
				"GROUP":  strings.Join(group, ","),
				"TICKET": strings.Join(ticket, ","),
				"SYMBOL": symbol,
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) DeletePositions(positions []uint64) error {
	cb := make(chan *ClientResponse)
	defer close(cb)

	tickets := make([]string, 0, len(positions))
	for _, d := range positions {
		tickets = append(tickets, fmt.Sprintf("%d", d))
	}

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandPositionDelete,
			Params: map[string]string{
				"TICKET": strings.Join(tickets, ","),
			},
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		return resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		break
	}

	return errors.New("timeout expired")
}

func (p *Pool) ClosePosition(position *Position) (*DealerUpdates, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	var pType int
	if position.Action == DealActionBuy {
		pType = DealActionSell
	} else {
		pType = DealActionBuy
	}

	params := struct {
		Position string `json:"Position"`
		Action   string `json:"Action"`
		Login    string `json:"Login"`
		Symbol   string `json:"Symbol"`
		Volume   string `json:"Volume"`
		TypeFill string `json:"TypeFill"`
		Type     string `json:"Type"`
		Comment  string `json:"Comment"`
	}{
		Position: fmt.Sprintf("%d", position.Position),
		Action:   fmt.Sprintf("%d", TradeActionDealerFirst),
		Login:    position.Login,
		Symbol:   position.Symbol,
		Volume:   fmt.Sprintf("%d", position.Volume),
		TypeFill: fmt.Sprintf("%d", OrderFillingFillFOK),
		Type:     fmt.Sprintf("%d", pType),
		Comment:  position.Comment,
	}
	payload, _ := json.Marshal(params)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:    MT5CommandDealerSend,
			Payload: string(payload),
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		if resp.Err != nil {
			return nil, resp.Err
		}
		return resp.Response.(*DealerUpdates), nil
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second * 2):
		break
	}

	return nil, errors.New("timeout expired")
}

func (c *MT5Client) getPositionsTotal(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	c.log.Debugf("#%d %s header: %s body: %s", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body)

	cmd, err := c.sendRequest(req)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s retcode error: %s", m.Cmd.Name, cmd.Params[MT5RetCode])}
		return
	}

	c.log.Debugf("%s response: %+v", cmd.Name, cmd.Params)

	total, err := strconv.Atoi(cmd.Params["TOTAL"])
	m.Cb <- &ClientResponse{
		Cmd:      m.Cmd,
		Response: &PositionsTotalResponse{Total: total},
		Err:      err,
	}
}

func (c *MT5Client) getPositions(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	c.log.Debugf("#%d %s header: %s body: %s", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body)

	cmd, err := c.sendRequest(req)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s retcode error: %s", m.Cmd.Name, cmd.Params[MT5RetCode])}
		return
	}

	c.log.Debugf("#%d %s response: %+v", c.clientId, cmd.Name, cmd.Params)

	positions := make([]Position, 0, 100)
	if err = json.Unmarshal([]byte(cmd.Payload), &positions); err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s response unmarshal error: %v", m.Cmd.Name, err)}
		return
	}

	if len(positions) > 0 {
		c.log.Debugf("#%d Positions received: %d", c.clientId, len(positions))
	}

	m.Cb <- &ClientResponse{
		Cmd:      m.Cmd,
		Response: &PositionsResponse{Positions: positions},
		Err:      err,
	}
}

func (c *MT5Client) deletePositions(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		safeSend(m.Cb, &ClientResponse{Cmd: m.Cmd, Err: err})
		return
	}

	c.log.Debugf("#%d %s header: %s body: %s", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body)

	cmd, err := c.sendRequest(req)
	if err != nil {
		safeSend(m.Cb, &ClientResponse{Cmd: m.Cmd, Err: err})
		return
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		safeSend(m.Cb, &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s retcode error: %s (%s)", m.Cmd.Name, cmd.Params[MT5RetCode], req)})
		return
	}

	c.log.Debugf("%s response: %+v", cmd.Name, cmd.Params)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: cmd,
		Err:      err,
	})
}
