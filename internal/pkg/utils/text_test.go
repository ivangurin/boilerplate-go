package utils

import "testing"

func TestToText(t *testing.T) {
	tests := []struct {
		name   string
		val    *string
		expect string
	}{
		{"nil pointer", nil, ""},
		{"non-nil pointer", Ptr("42"), "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToText(tt.val); got != tt.expect {
				t.Errorf("ToText() = %v, want %v", got, tt.expect)
			}
		})
	}
}
