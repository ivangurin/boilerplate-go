package utils

import "testing"

func TestJoinInt(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		sep    string
		result string
	}{
		{
			name:   "empty slice",
			slice:  []int{},
			sep:    ",",
			result: "",
		},
		{
			name:   "single element",
			slice:  []int{1},
			sep:    ",",
			result: "1",
		},
		{
			name:   "multiple elements",
			slice:  []int{1, 2, 3},
			sep:    ",",
			result: "1,2,3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinInt(tt.slice, tt.sep); got != tt.result {
				t.Errorf("JoinInt() = %v, want %v", got, tt.result)
			}
		})
	}
}

func TestHideText(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		visibleChars int
		expected     string
	}{
		{
			name:         "short text",
			input:        "hi",
			visibleChars: 3,
			expected:     "hi",
		},
		{
			name:         "exact length",
			input:        "hello",
			visibleChars: 5,
			expected:     "hello",
		},
		{
			name:         "long string",
			input:        "hello world",
			visibleChars: 3,
			expected:     "hel********",
		},
		{
			name:         "zero visible chars",
			input:        "hello",
			visibleChars: 0,
			expected:     "*****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HideText(tt.input, tt.visibleChars); got != tt.expected {
				t.Errorf("HideText() = %v, want %v", got, tt.expected)
			}
		})
	}
}
