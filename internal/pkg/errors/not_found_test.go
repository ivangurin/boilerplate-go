package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundError(t *testing.T) {
	err := errors.New("тест1")
	require.False(t, IsErrNotFound(err))

	err = NewNotFoundError("тест2")
	require.Equal(t, "тест2", err.Error())
	require.True(t, IsErrNotFound(err))
}
