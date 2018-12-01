package jsonrpc2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpProvider_Call(t *testing.T) {
	client := NewHttpProvider("https://gurujsonrpc.appspot.com/guru")
	resp, err := client.Call("guru.test", "Guru")
	if err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, "Hello Guru!", resp.Result)
	}
}
