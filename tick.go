package mt5client

import (
	"fmt"
	"time"
)

func (p *Pool) GetTickHistory(symbol string, from, to int64, data string) error {
	cb := make(chan *ClientResponse)
	defer close(cb)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandTickGetHistory,
			Params: map[string]string{
				"SYMBOL": symbol,
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
				"DATA":   data,
			},
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		if resp.Err != nil {
			return resp.Err
		}
		return resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		break
	}

	return nil
}

func (p *Pool) GetChart(symbol string, from, to int64, data string) error {
	cb := make(chan *ClientResponse)
	defer close(cb)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name: MT5CommandChartGet,
			Params: map[string]string{
				"SYMBOL": symbol,
				"FROM":   fmt.Sprintf("%d", from),
				"TO":     fmt.Sprintf("%d", to),
				"DATA":   data,
			},
		},
		Cb: cb,
	}

	select {
	case resp := <-cb:
		if resp.Err != nil {
			return resp.Err
		}
		return resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		break
	}

	return nil
}

func (c *MT5Client) getTickHistory(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	c.log.Infof("#%d %s header: %s body: %s", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body)

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

	fmt.Println(cmd.Payload)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: cmd.Payload,
		Err:      err,
	})
}

func (c *MT5Client) getChart(m *ClientControlMessage) {
	req, body, err := c.makeRequest(m.Cmd)
	if err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
		return
	}

	c.log.Infof("#%d %s header: %s body: %s", c.clientId, m.Cmd.Name, req[:MT5HeaderLength], body)

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

	fmt.Println(cmd.Payload)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: cmd.Payload,
		Err:      err,
	})
}
