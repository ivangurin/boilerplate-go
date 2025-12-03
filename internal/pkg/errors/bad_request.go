package errors

import "errors"

type ErrBadRequest struct {
	Msg string `json:"error"`
}

func NewBadRequestError(msg string) *ErrBadRequest {
	return &ErrBadRequest{
		Msg: msg,
	}
}

func (e *ErrBadRequest) Error() string {
	return e.Msg
}

func IsErrBadRequest(err error) bool {
	var errBadRequest *ErrBadRequest
	return errors.As(err, &errBadRequest)
}
