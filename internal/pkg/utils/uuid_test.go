package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	uuid := UUID()
	require.NotEmpty(t, uuid)
}
