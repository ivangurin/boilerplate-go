package errors

import "errors"

type ErrNotFound struct {
	Msg string `json:"error"`
}

func NewNotFoundError(msg string) *ErrNotFound {
	return &ErrNotFound{
		Msg: msg,
	}
}

func (e *ErrNotFound) Error() string {
	return e.Msg
}

func IsErrNotFound(err error) bool {
	var errBadRequest *ErrNotFound
	return errors.As(err, &errBadRequest)
}
