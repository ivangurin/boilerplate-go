package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimitOffset(t *testing.T) {
	type TestCase struct {
		Name     string
		Limit    *int
		Offset   *int
		Expect   int
		Contains []int
	}

	data := []int{1, 2, 3, 4, 5}

	testCases := []TestCase{
		{
			Name:     "No limit and no offset",
			Limit:    nil,
			Offset:   nil,
			Expect:   5,
			Contains: []int{1, 2, 3, 4, 5},
		},
		{
			Name:     "limit 1, no offset",
			Limit:    Ptr(1),
			Offset:   nil,
			Expect:   1,
			Contains: []int{1},
		},
		{
			Name:     "limit 2, no offset",
			Limit:    Ptr(2),
			Offset:   nil,
			Expect:   2,
			Contains: []int{1, 2},
		},
		{
			Name:     "limit 2, offset 1",
			Limit:    Ptr(2),
			Offset:   Ptr(1),
			Expect:   2,
			Contains: []int{2, 3},
		},
		{
			Name:     "limit 2, offset 2",
			Limit:    Ptr(2),
			Offset:   Ptr(2),
			Expect:   2,
			Contains: []int{3, 4},
		},
		{
			Name:     "limit 2, offset 4",
			Limit:    Ptr(2),
			Offset:   Ptr(4),
			Expect:   1,
			Contains: []int{5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := LimitOffset(data, tc.Limit, tc.Offset)
			require.Len(t, result, tc.Expect)
			for _, v := range tc.Contains {
				require.Contains(t, result, v)
			}
		})
	}
}
