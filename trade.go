package mt5client

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (p *Pool) Balance(login string, opType uint8, balance float64, comment string) (uint64, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	p.log.Debugf("MT5 BALANCE REQUEST: login=%s type=%d balance=%f", login, opType, balance)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandTradeBalance,
			Params: map[string]string{
				"LOGIN":   login,
				"TYPE":    fmt.Sprintf("%d", opType),
				"BALANCE": fmt.Sprintf("%f", balance),
				"COMMENT": comment,
			},
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		if resp.Err != nil {
			return 0, resp.Err
		}
		return resp.Response.(uint64), resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		break
	}

	return 0, errors.New("timeout expired")
}

func (c *MT5Client) balance(m *ClientControlMessage) {
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
		safeSend(m.Cb, &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s retcode error: %s", m.Cmd.Name, cmd.Params[MT5RetCode])})
		return
	}

	c.log.Debugf("#%d %s response: %+v", c.clientId, cmd.Name, cmd.Params)

	ticket, err := strconv.ParseUint(cmd.Params["TICKET"], 10, 64)
	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: ticket,
		Err:      err,
	})
}
