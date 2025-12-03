package errors

import "errors"

type ErrUnauthorized struct {
	Msg string `json:"error"`
}

func NewUnauthorizedError(msg string) *ErrUnauthorized {
	return &ErrUnauthorized{
		Msg: msg,
	}
}

func (e *ErrUnauthorized) Error() string {
	return e.Msg
}

func IsErrUnauthorized(err error) bool {
	var errUnauthorized *ErrUnauthorized
	return errors.As(err, &errUnauthorized)
}
