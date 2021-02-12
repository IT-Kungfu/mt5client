package mt5client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (p *Pool) GetOrdersTotal(login string) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandOrderGetTotal,
			Params: map[string]string{"LOGIN": login},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetOrdersHistoryTotal(login string, from, to int64) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandOrderGetHistoryTotal,
			Params: map[string]string{
				"LOGIN": login,
				"FROM":  fmt.Sprintf("%d", from),
				"TO":    fmt.Sprintf("%d", to),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetOrdersPage(login string, offset, total int) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandOrderGetPage,
			Params: map[string]string{
				"LOGIN":  login,
				"OFFSET": fmt.Sprintf("%d", offset),
				"TOTAL":  fmt.Sprintf("%d", total),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetOrdersBatch(login, group, ticket []string) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandOrderGetBatch,
			Params: map[string]string{
				"LOGIN":  strings.Join(login, ","),
				"GROUP":  strings.Join(group, ","),
				"TICKET": strings.Join(ticket, ","),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetOrdersHistoryPage(login string, from, to int64, offset, total int) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandOrderGetHistoryPage,
			Params: map[string]string{
				"LOGIN":  login,
				"OFFSET": fmt.Sprintf("%d", offset),
				"TOTAL":  fmt.Sprintf("%d", total),
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
			},
		},
		Cb: p.cb,
	}
}

func (p *Pool) GetOrdersHistoryBatch(login, group, ticket []string, from, to int64) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandOrderGetHistoryBatch,
			Params: map[string]string{
				"LOGIN":  strings.Join(login, ","),
				"GROUP":  strings.Join(group, ","),
				"TICKET": strings.Join(ticket, ","),
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
			},
		},
		Cb: p.cb,
	}
}

func (c *MT5Client) getOrdersTotal(m *ClientControlMessage) {
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
		Response: &OrdersTotalResponse{Total: total},
		Err:      err,
	}
}

func (c *MT5Client) getOrders(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	c.log.Debugf("#%d %s header: %s body: %s...", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body[:32])

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

	orders := make([]Order, 0, 100)
	err = json.Unmarshal([]byte(cmd.Payload), &orders)

	c.log.Infof("#%d Orders received: %d", c.clientId, len(orders))

	m.Cb <- &ClientResponse{
		Cmd:      m.Cmd,
		Response: &OrdersResponse{Orders: orders},
		Err:      err,
		ClientId: c.clientId,
	}
}
