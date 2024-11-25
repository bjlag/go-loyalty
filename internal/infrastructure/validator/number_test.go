package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "positive",
			number: "12345678903",
			want:   true,
		},
		{
			name:   "negative",
			number: "12345678902",
			want:   false,
		},
		{
			name:   "empty",
			number: "",
			want:   false,
		},
		{
			name:   "not_digit_in_number",
			number: "123456789A3",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckLuhn(tt.number)

			assert.Equal(t, tt.want, got)
		})
	}
}
