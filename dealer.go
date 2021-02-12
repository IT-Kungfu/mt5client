package mt5client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (p *Pool) CreateEmptyDeal(login string, comment string) (*DealerUpdates, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := struct {
		Action  string `json:"Action"`
		Login   string `json:"Login"`
		Type    string `json:"Type"`
		Comment string `json:"Comment"`
	}{
		Action:  fmt.Sprintf("%d", TradeActionDealerFirst),
		Login:   login,
		Type:    fmt.Sprintf("%d", DealActionBalance),
		Comment: comment,
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

func (c *MT5Client) sendDealer(m *ClientControlMessage) {
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

	c.log.Debugf("%s response: %+v (%s)", cmd.Name, cmd.Params, cmd.Payload)

	type Resp struct {
		Id string `json:"id"`
	}

	resp := &Resp{}
	if err = json.Unmarshal([]byte(cmd.Payload), resp); err != nil {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: fmt.Errorf("%s response unmarshal error: %v", m.Cmd.Name, err)}
		return
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////

	updResp := c.getDealerUpdates(&MT5Command{
		Name:   MT5CommandDealerUpdates,
		Params: map[string]string{"ID": resp.Id},
	})

	m.Cb <- &ClientResponse{
		Cmd:      m.Cmd,
		Response: updResp.Response,
		Err:      updResp.Err,
	}
}

func (c *MT5Client) getDealerUpdates(m *MT5Command) *ClientResponse {
	req, body, err := c.makeRequest(m)
	if err != nil {
		return &ClientResponse{Cmd: m, Err: err}
	}

	c.log.Debugf("#%d %s header: %s body: %s", c.clientId, m.Name, req[:MT5HeaderLength], body)

	cmd, err := c.sendRequest(req)
	if err != nil {
		return &ClientResponse{Cmd: m, Err: err}
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		return &ClientResponse{Cmd: m, Err: fmt.Errorf("%s retcode error: %s", m.Name, cmd.Params[MT5RetCode])}
	}

	c.log.Debugf("%s response: %+v (%s)", cmd.Name, cmd.Params, cmd.Payload)

	resp := make(map[string][]*DealerUpdates, 1)
	if err = json.Unmarshal([]byte(cmd.Payload), &resp); err != nil {
		return &ClientResponse{Cmd: m, Err: fmt.Errorf("%s response unmarshal error: %v", m.Name, err)}
	}

	if _, ok := resp[m.Params["ID"]]; !ok {
		return &ClientResponse{Cmd: m, Err: fmt.Errorf("%s request id not found", m.Name)}
	}

	result := &DealerUpdatesResult{}
	answer := &DealerUpdatesAnswer{}
	for _, r := range resp[m.Params["ID"]] {
		if r.Result != nil {
			result = r.Result
		} else if r.Answer != nil {
			answer = r.Answer
		}
	}

	return &ClientResponse{
		Cmd: m,
		Response: &DealerUpdates{
			Result: result,
			Answer: answer,
		},
		Err: err,
	}
}
