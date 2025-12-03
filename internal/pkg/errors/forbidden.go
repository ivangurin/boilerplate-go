package errors

import "errors"

type ErrForbidden struct {
	Msg string `json:"error"`
}

func NewForbiddenError(msg string) *ErrForbidden {
	return &ErrForbidden{
		Msg: msg,
	}
}

func (e *ErrForbidden) Error() string {
	return e.Msg
}

func IsErrForbidden(err error) bool {
	var errForbidden *ErrForbidden
	return errors.As(err, &errForbidden)
}
