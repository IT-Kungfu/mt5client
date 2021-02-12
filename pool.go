package mt5client

import (
	"context"
	"github.com/IT-Kungfu/logger"
)

type ClientControlMessage struct {
	Cmd *MT5Command
	Cb  chan *ClientResponse
}

type ClientResponse struct {
	Cmd      *MT5Command
	Response interface{}
	Err      error
	ClientId int
}

type Client struct {
	handler   *MT5Client
	controlCh chan *ClientControlMessage
}

type Pool struct {
	cfg        *Config
	log        *logger.Logger
	clients    []*Client
	cb         chan *ClientResponse
	poolSize   int
	nextClient int
}

func NewMT5ClientPool(ctx context.Context, poolSize int) (*Pool, error) {
	services := ctx.Value("services").(map[string]interface{})
	cb := make(chan *ClientResponse, poolSize)
	pool := &Pool{
		cfg:      services["mt5cfg"].(*Config),
		log:      services["log"].(*logger.Logger),
		clients:  make([]*Client, 0, poolSize),
		poolSize: poolSize,
		cb:       cb,
	}

	for i := 0; i < poolSize; i++ {
		controlCh := make(chan *ClientControlMessage)
		services["controlCh"] = controlCh
		services["clientId"] = i
		clientCtx := context.WithValue(ctx, "services", services)

		h, err := NewMT5Client(clientCtx)
		if err != nil {
			return nil, err
		}
		c := &Client{
			handler:   h,
			controlCh: controlCh,
		}
		pool.clients = append(pool.clients, c)
	}

	return pool, nil
}

func (p *Pool) Response() chan *ClientResponse {
	return p.cb
}

func (p *Pool) getClient() *Client {
	c := p.clients[p.nextClient]
	if p.nextClient+1 < p.poolSize {
		p.nextClient++
	} else {
		p.nextClient = 0
	}
	return c
}

func (p *Pool) Close() {
	for _, c := range p.clients {
		c.handler.Quit()
	}
}
