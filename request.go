package jsonrpc2

import "encoding/json"

type Request struct {
	ID     int
	Method string
	Params *json.RawMessage
}
