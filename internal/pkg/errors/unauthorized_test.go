package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnauthorizedError(t *testing.T) {
	err := errors.New("тест1")
	require.False(t, IsErrUnauthorized(err))

	err = NewUnauthorizedError("тест2")
	require.Equal(t, "тест2", err.Error())
	require.True(t, IsErrUnauthorized(err))
}
