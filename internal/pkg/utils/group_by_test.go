package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToGroupBy(t *testing.T) {
	type TestData struct {
		ID      int
		GroupID int
		Value   string
	}

	type TestCase struct {
		Name     string
		Data     []TestData
		Expected map[int][]TestData
	}

	testCases := []TestCase{
		{
			Name: "one row with uniq key",
			Data: []TestData{
				{ID: 1, GroupID: 1, Value: "A"},
			},
			Expected: map[int][]TestData{
				1: {
					{ID: 1, GroupID: 1, Value: "A"},
				},
			},
		},
		{
			Name: "two rows with uniq key",
			Data: []TestData{
				{ID: 1, GroupID: 1, Value: "A"},
				{ID: 2, GroupID: 2, Value: "B"},
			},
			Expected: map[int][]TestData{
				1: {
					{ID: 1, GroupID: 1, Value: "A"},
				},
				2: {
					{ID: 2, GroupID: 2, Value: "B"},
				},
			},
		},
		{
			Name: "three rows with common key",
			Data: []TestData{
				{ID: 1, GroupID: 1, Value: "A"},
				{ID: 2, GroupID: 2, Value: "B"},
				{ID: 3, GroupID: 2, Value: "C"},
			},
			Expected: map[int][]TestData{
				1: {
					{ID: 1, GroupID: 1, Value: "A"},
				},
				2: {
					{ID: 2, GroupID: 2, Value: "B"},
					{ID: 3, GroupID: 2, Value: "C"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got := ToGroupBy(tc.Data, func(td TestData) (int, TestData) {
				return td.GroupID, td
			})
			for expectedKey, expectedValue := range tc.Expected {
				gotValue, exists := got[expectedKey]
				require.True(t, exists)
				require.Equal(t, expectedValue, gotValue)
			}
		})
	}
}
