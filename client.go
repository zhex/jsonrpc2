package jsonrpc2

import (
	"net/url"
)

type Provider interface {
	Call(method string, param ...interface{}) (*Response, error)
}

type Request struct {
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type Response struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func New(endpoint string) (Provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "ws" || u.Scheme == "wss" {
		return nil, nil
	} else {
		return NewHttpProvider(u.String()), nil
	}
}
