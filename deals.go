package mt5client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (p *Pool) GetDealsTotal(login string, from, to int64) (int, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandDealGetTotal,
			Params: map[string]string{
				"LOGIN": login,
				"FROM":  fmt.Sprintf("%d", from),
				"TO":    fmt.Sprintf("%d", to),
			},
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		if resp.Err != nil {
			return 0, resp.Err
		}
		return resp.Response.(*DealsTotalResponse).Total, resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		break
	}

	return 0, errors.New("timeout expired")
}

func (p *Pool) GetDealsPage(login string, from, to int64, offset, total int) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandDealGetPage,
			Params: map[string]string{
				"LOGIN":  login,
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
				"OFFSET": fmt.Sprintf("%d", offset),
				"TOTAL":  fmt.Sprintf("%d", total),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetDealsBatch(login string, group, ticket []string, from, to int64) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandDealGetBatch,
			Params: map[string]string{
				"LOGIN":  login,
				"GROUP":  strings.Join(group, ","),
				"TICKET": strings.Join(ticket, ","),
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) DeleteDeals(deals []uint64) error {
	cb := make(chan *ClientResponse)
	defer close(cb)

	tickets := make([]string, 0, len(deals))
	for _, d := range deals {
		tickets = append(tickets, fmt.Sprintf("%d", d))
	}

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandDealDelete,
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

func (c *MT5Client) getDealsTotal(m *ClientControlMessage) {
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

	total, err := strconv.Atoi(cmd.Params["TOTAL"])
	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: &DealsTotalResponse{Total: total},
		Err:      err,
	})
}

func (c *MT5Client) getDeals(m *ClientControlMessage) {
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
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s retcode error: %s (%s)", m.Cmd.Name, cmd.Params[MT5RetCode], req)}
		return
	}

	c.log.Debugf("#%d %s response: %+v", c.clientId, cmd.Name, cmd.Params)

	deals := make([]Deal, 0, 100)
	err = json.Unmarshal([]byte(cmd.Payload), &deals)

	if len(deals) > 0 {
		c.log.Debugf("#%d Deals received: %d", c.clientId, len(deals))
	}

	m.Cb <- &ClientResponse{
		Cmd: m.Cmd,
		Response: &DealsResponse{
			Deals: deals,
		},
		Err: err,
	}
}

func (c *MT5Client) deleteDeals(m *ClientControlMessage) {
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
