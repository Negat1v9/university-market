package httpresponse

import "fmt"

type HttpError struct {
	Code int
	Msg  string
	// field that is filled in to transmit the error to the top for logging (optional)
	LogError error
}

func NewError(code int, msg string) *HttpError {
	return &HttpError{
		Code: code,
		Msg:  msg,
	}
}

func ServerError() *HttpError {
	r := &HttpError{
		Code: 500,
		Msg:  "server error, try again later",
	}
	return r
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("error: %s, code %d", e.Msg, e.Code)
}
