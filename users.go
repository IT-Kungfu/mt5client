package mt5client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MT5UserRequest struct {
	Name             string
	Phone            string
	Group            string
	Login            string
	PasswordMain     string
	PasswordInvestor string
	Rights           string
	Email            string
	Leverage         string
}

func (p *Pool) GetUser(login string) (*User, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := map[string]string{
		"LOGIN": login,
	}

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserGet,
			Params: params,
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		if resp.Err != nil {
			return nil, resp.Err
		} else {
			return resp.Response.(*User), resp.Err
		}
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		return nil, errors.New("timeout expired")
	}
}

func (p *Pool) GetUsersBatch(login []string) ([]*User, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := map[string]string{
		"LOGIN": strings.Join(login, ","),
	}

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserGetBatch,
			Params: params,
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		if resp.Err != nil {
			return nil, resp.Err
		} else {
			return resp.Response.([]*User), resp.Err
		}
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		return nil, errors.New("timeout expired")
	}
}

func (p *Pool) AddUser(req *MT5UserRequest) (*ClientResponse, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := map[string]string{
		"NAME":          req.Name,
		"PHONE":         req.Phone,
		"GROUP":         req.Group,
		"LOGIN":         req.Login,
		"PASS_MAIN":     req.PasswordMain,
		"PASS_INVESTOR": req.PasswordInvestor,
		"RIGHTS":        req.Rights,
		"EMAIL":         p.prepareEmail(req.Email),
		"LEVERAGE":      req.Leverage,
	}

	p.log.Debugf("MT5 ADD USER REQUEST: %+v", params)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserAdd,
			Params: params,
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		return resp, resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		return nil, errors.New("timeout expired")
	}
}

func (p *Pool) UpdateUser(req *MT5UserRequest) (*ClientResponse, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := map[string]string{
		"NAME":          req.Name,
		"PHONE":         req.Phone,
		"GROUP":         req.Group,
		"LOGIN":         req.Login,
		"PASS_MAIN":     req.PasswordMain,
		"PASS_INVESTOR": req.PasswordInvestor,
		"RIGHTS":        req.Rights,
		"EMAIL":         p.prepareEmail(req.Email),
		"LEVERAGE":      req.Leverage,
	}

	p.log.Debugf("MT5 UPDATE USER REQUEST: %+v", params)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserUpdate,
			Params: params,
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		return resp, resp.Err
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		return nil, errors.New("timeout expired")
	}
}

func (p *Pool) DeleteUser(login string, timeout int) (*ClientResponse, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserDelete,
			Params: map[string]string{"LOGIN": login},
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		fmt.Println(resp)
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("timeout expired")
	}

	return resp, nil
}

func (p *Pool) GetUserAccounts(login []string) (map[uint64]*UserAccount, error) {
	cb := make(chan *ClientResponse)
	defer close(cb)

	params := map[string]string{
		"LOGIN": strings.Join(login, ","),
	}

	p.getClient().controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{
			Name:   MT5CommandUserAccountGetBatch,
			Params: params,
		},
		Cb: cb,
	}

	resp := &ClientResponse{}

	select {
	case resp = <-cb:
		if resp.Err != nil {
			return nil, resp.Err
		} else {
			respUsers := resp.Response.([]*UserAccount)
			users := make(map[uint64]*UserAccount, len(respUsers))
			for _, u := range respUsers {
				id, _ := strconv.ParseUint(u.Login, 10, 64)
				users[id] = u
			}
			return users, resp.Err
		}
	case <-time.After(time.Duration(p.cfg.MT5RequestTimeout) * time.Second):
		return nil, errors.New("timeout expired")
	}
}

func (c *MT5Client) addUpdateUser(m *ClientControlMessage) {
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
		safeSend(m.Cb, &ClientResponse{Cmd: m.Cmd, Err: errors.New(cmd.Params[MT5RetCode])})
		return
	}

	c.log.Debugf("#%d %s response: %+v", c.clientId, cmd.Name, cmd.Params)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: cmd.Payload,
		Err:      err,
		ClientId: c.clientId,
	})
}

func (c *MT5Client) deleteUser(m *ClientControlMessage) {
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

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: nil,
		Err:      err,
		ClientId: c.clientId,
	})
}

func (c *MT5Client) getUser(m *ClientControlMessage) {
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

	user := &User{}
	err = json.Unmarshal([]byte(cmd.Payload), user)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: user,
		Err:      err,
		ClientId: c.clientId,
	})
}

func (c *MT5Client) getUsersBatch(m *ClientControlMessage) {
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

	users := make([]*User, 0, len(m.Cmd.Params["LOGIN"]))
	err = json.Unmarshal([]byte(cmd.Payload), &users)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: users,
		Err:      err,
		ClientId: c.clientId,
	})
}

func (c *MT5Client) getUserAccounts(m *ClientControlMessage) {
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

	users := make([]*UserAccount, 0, len(m.Cmd.Params["LOGIN"]))
	err = json.Unmarshal([]byte(cmd.Payload), &users)

	safeSend(m.Cb, &ClientResponse{
		Cmd:      m.Cmd,
		Response: users,
		Err:      err,
		ClientId: c.clientId,
	})
}

func (p *Pool) prepareEmail(email string) string {
	mask := func(value string) string {
		m := ""
		for i, v := range value {
			if (i != 0 && i != len(value)-1) || (i != 0 && len(value) <= 2) {
				m += "*"
			} else {
				m += fmt.Sprintf("%c", v)
			}
		}
		return m
	}
	semail := strings.Split(email, "@")
	domain := strings.Split(semail[1], ".")
	return mask(semail[0]) + "@" + mask(domain[0]) + "." + domain[1]
}
