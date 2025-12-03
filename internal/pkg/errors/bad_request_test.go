package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBadRequestError(t *testing.T) {
	err := errors.New("тест1")
	require.False(t, IsErrBadRequest(err))

	err = NewBadRequestError("тест2")
	require.Equal(t, "тест2", err.Error())
	require.True(t, IsErrBadRequest(err))
}
