package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewHttpProvider(endpoint string) *HttpProvider {
	return &HttpProvider{
		Client:   &http.Client{},
		Endpoint: endpoint,
		ID:       0,
	}
}

type HttpProvider struct {
	Client   *http.Client
	Endpoint string
	ID       int
}

func (p *HttpProvider) Call(method string, params ...interface{}) (*Response, error) {
	p.ID++
	return p.request(&Request{
		ID:      p.ID,
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
	fmt.Println(string(body))
	r, err := http.NewRequest("POST", p.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	resp, err := p.Client.Do(r)
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
