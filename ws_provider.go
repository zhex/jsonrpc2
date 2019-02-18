package jsonrpc2

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func NewWsProvider(endpoint string) *WsProvider {
	client, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &WsProvider{
		client:    client,
		endpoint:  endpoint,
		id:        0,
		msgChan:   make(map[int](chan *Response)),
		ctx:       ctx,
		ctxCancel: cancel,
	}

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
				resp := &Response{}
				err := client.ReadJSON(resp)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if c, ok := p.msgChan[resp.ID]; ok {
					c <- resp
				}
			}
		}
	}()

	return p
}

type WsProvider struct {
	client    *websocket.Conn
	endpoint  string
	id        int
	msgChan   map[int](chan *Response)
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (p *WsProvider) Call(method string, params ...interface{}) (*Response, error) {
	p.id++
	req := &Request{
		ID:      p.id,
		Method:  method,
		JSONRPC: "2.0",
		Params:  params,
	}
	p.msgChan[p.id] = make(chan *Response, 1)
	err := p.client.WriteJSON(req)
	if err != nil {
		return nil, err
	}

	resp := <-p.msgChan[p.id]
	delete(p.msgChan, p.id)
	return resp, nil
}

func (p *WsProvider) Subscribe(method string, callback func(*Response)) error {
	// todo
	return nil
}

func (p *WsProvider) Close() {
	p.ctxCancel()
	for key := range p.msgChan {
		delete(p.msgChan, key)
	}
	p.client.Close()
}
