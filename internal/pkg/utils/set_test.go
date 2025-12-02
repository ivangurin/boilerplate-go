package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run(".Add", func(t *testing.T) {
		obj := Set[int64]{}

		obj.Add(42)

		_, ok := obj[42]
		assert.True(t, ok, "It should add value as map key.")
	})

	t.Run(".Has", func(t *testing.T) {
		obj := Set[int64]{1: {}, 2: {}}

		assert.True(t, obj.Has(1), "It should return true for existing element.")
		assert.False(t, obj.Has(0), "It should return false for absent element.")
	})

	t.Run(".ToSlice", func(t *testing.T) {
		obj := Set[int64]{1: {}, 2: {}}
		expected := []int64{1, 2}

		actual := obj.ToSlice()

		assert.ElementsMatch(t, expected, actual, "It should return map keys as slice.")
	})

	t.Run(".Delete", func(t *testing.T) {
		obj := Set[int64]{1: {}, 2: {}, 3: {}}
		expected := []int64{1, 3}

		obj.Delete(2)
		obj.Delete(123)

		assert.Len(t, obj, 2)
		assert.False(t, obj.Has(2))
		assert.ElementsMatch(t, expected, obj.ToSlice(), "It should return map keys as slice.")
	})
}

func TestToSet(t *testing.T) {
	type args struct {
		collection []int
		selector   func(arg int) int
	}
	tests := []struct {
		name string
		args args
		want Set[int]
	}{
		{
			name: "Test 1",
			args: args{
				collection: []int{1, 2, 2, 3, 3, 3},
				selector:   func(arg int) int { return arg },
			},
			want: Set[int]{1: {}, 2: {}, 3: {}},
		},
		{
			name: "Test 2",
			args: args{
				collection: []int{1, 1, 1, 1, 1, 1},
				selector:   func(arg int) int { return arg },
			},
			want: Set[int]{1: {}},
		},
		{
			name: "Test 3",
			args: args{
				collection: []int{-2, -1, 0, 1, 2},
				selector:   func(arg int) int { return arg * arg }, // square of a number
			},
			want: Set[int]{0: {}, 1: {}, 4: {}},
		},
		{
			name: "Test 4",
			args: args{
				collection: []int{},
				selector:   func(arg int) int { return arg },
			},
			want: Set[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSet(tt.args.collection, tt.args.selector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
