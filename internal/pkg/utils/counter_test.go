package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	counter := Counter[string]{}
	counter.Add("a")
	counter.Add("b")
	counter.Add("b")
	counter.Add("c")
	counter.Add("c")
	counter.Add("c")

	require.True(t, counter.Has("a"))
	require.True(t, counter.Has("b"))
	require.True(t, counter.Has("c"))
	require.False(t, counter.Has("d"))
	require.Len(t, counter.ToSlice(), 3)
	require.Equal(t, 1, counter.Count("a"))
	require.Equal(t, 2, counter.Count("b"))
	require.Equal(t, 3, counter.Count("c"))
	require.Equal(t, 0, counter.Count("d"))
	require.Len(t, counter.GT(1), 2)
}

func TestToCounter(t *testing.T) {
	type item struct {
		id   int
		name string
	}

	items := []item{
		{id: 1, name: "a"},
		{id: 2, name: "b"},
		{id: 3, name: "b"},
		{id: 4, name: "c"},
		{id: 5, name: "c"},
		{id: 5, name: "c"},
	}

	counter := ToCounter(items, func(i item) string {
		return i.name
	})

	require.Equal(t, 1, counter["a"])
	require.Equal(t, 2, counter["b"])
	require.Equal(t, 3, counter["c"])
	require.Equal(t, 0, counter["d"])
}
