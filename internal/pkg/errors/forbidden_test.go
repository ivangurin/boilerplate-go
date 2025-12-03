package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForbiddenError(t *testing.T) {
	err := errors.New("тест1")
	require.False(t, IsErrForbidden(err))

	err = NewForbiddenError("тест2")
	require.Equal(t, "тест2", err.Error())
	require.True(t, IsErrForbidden(err))
}
