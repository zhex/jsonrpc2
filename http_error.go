package jsonrpc2

type HttpError struct {
	Code int
	err  error
}

func (e *HttpError) Error() string {
	return e.err.Error()
}
