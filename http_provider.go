package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewHttpProvider(endpoint string) *HttpProvider {
	return &HttpProvider{
		client:   &http.Client{},
		endpoint: endpoint,
		id:       0,
	}
}

type HttpProvider struct {
	client   *http.Client
	endpoint string
	id       int
}

func (p *HttpProvider) Call(method string, params ...interface{}) (*Response, error) {
	p.id++
	return p.request(&Request{
		ID:      p.id,
		Method:  method,
		JSONRPC: "2.0",
		Params:  params,
	})
}

func (p *HttpProvider) request(req *Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", p.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResponse *Response
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(&rpcResponse)

	if err != nil {
		if resp.StatusCode >= 400 {
			err = &HttpError{
				Code: resp.StatusCode,
				err:  fmt.Errorf("RPC call [%v] failed: %v", req.Method, err.Error()),
			}
		}
	}

	return rpcResponse, err
}
