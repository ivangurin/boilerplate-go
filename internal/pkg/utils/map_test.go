package utils

import (
	"fmt"
	"maps"
	"reflect"
	"testing"
)

func TestToMapWithIntKeysAndStringValues(t *testing.T) {
	list := []int{1, 2, 3, 4, 5}
	result := ToMap(list, func(arg int) (int, string) {
		return arg, fmt.Sprintf("Number %d", arg)
	})
	expected := Map[int, string]{
		1: "Number 1",
		2: "Number 2",
		3: "Number 3",
		4: "Number 4",
		5: "Number 5",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestToMapWithEmptySlice(t *testing.T) {
	var list []int
	result := ToMap(list, func(arg int) (int, string) {
		return arg, fmt.Sprintf("Number %d", arg)
	})
	if len(result) != 0 {
		t.Errorf("Expected an empty map, but got %v", result)
	}
}

func TestToMapByField(t *testing.T) {
	type TestData struct {
		ID    int
		Value string
	}

	testCases := []struct {
		name       string
		collection []TestData
		selector   func(arg TestData) int
		want       map[int]TestData
	}{
		{
			name: "Test with non-empty slice",
			collection: []TestData{
				{ID: 1, Value: "First"},
				{ID: 2, Value: "Second"},
				{ID: 3, Value: "Third"},
			},
			selector: func(arg TestData) int {
				return arg.ID
			},
			want: Map[int, TestData]{
				1: {ID: 1, Value: "First"},
				2: {ID: 2, Value: "Second"},
				3: {ID: 3, Value: "Third"},
			},
		},
		{
			name:       "Test with empty slice",
			collection: []TestData{},
			selector: func(arg TestData) int {
				return arg.ID
			},
			want: Map[int, TestData]{},
		},
		{
			name: "Test with duplicate keys",
			collection: []TestData{
				{ID: 1, Value: "First"},
				{ID: 1, Value: "Overwritten"},
				{ID: 2, Value: "Second"},
			},
			selector: func(arg TestData) int {
				return arg.ID
			},
			want: Map[int, TestData]{
				1: {ID: 1, Value: "Overwritten"},
				2: {ID: 2, Value: "Second"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToMapByField(tc.collection, tc.selector)
			if !maps.Equal(got, tc.want) {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}
