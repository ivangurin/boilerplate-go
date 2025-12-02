package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPtr(t *testing.T) {
	val := 42
	ptr := Ptr(val)
	require.NotNil(t, ptr)
	require.Equal(t, val, *ptr)
}

func TestDePtr(t *testing.T) {
	ptr := Ptr(42)
	require.NotNil(t, ptr)
	require.Equal(t, 42, DePtr(ptr))

	var valInt *int
	require.Equal(t, 0, DePtr(valInt))

	var valStr *string
	require.Equal(t, "", DePtr(valStr))

	var valBool *bool
	require.Equal(t, false, DePtr(valBool))
}
