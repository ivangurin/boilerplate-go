package pwd

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	password := gofakeit.Word()
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.True(t, CheckPasswordHash(password, hash))
}
