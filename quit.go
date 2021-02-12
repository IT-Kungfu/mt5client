package mt5client

import "fmt"

func (c *MT5Client) Quit() {
	defer func() {
		close(c.controlCh)
	}()
	cb := make(chan *ClientResponse)
	c.controlCh <- &ClientControlMessage{
		Cmd: &MT5Command{Name: MT5CommandQuit},
		Cb:  cb,
	}
	<-cb
}

func (c *MT5Client) quit(m *ClientControlMessage) {
	var err error
	defer func() {
		m.Cb <- &ClientResponse{Cmd: m.Cmd, Err: err}
	}()
	if c.conn != nil {
		c.log.Debugf("MT5Client #%d quit", c.clientId)
		body := fmt.Sprintf("%s%s", MT5CommandQuit, MT5PacketSeparator)
		request, _ := makePacket(body, 0, 0)
		_, err = c.conn.Write(request)
		if err != nil {
			return
		}
		err = c.conn.Close()
	}
}
