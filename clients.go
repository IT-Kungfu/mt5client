package mt5client

import (
	"fmt"
)

func (p *Pool) GetClientIds(group string) {
	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandClientGetIds,
			Params: map[string]string{"GROUP": group},
		},
		Cb: p.cb,
	}
}

func (c *MT5Client) getClients(m *ClientControlMessage) {
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

	m.Cb <- &ClientResponse{
		Cmd:      m.Cmd,
		Response: &ClientsResponse{},
		Err:      err,
	}
}
