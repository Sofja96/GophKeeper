package shared

import "testing"

func TestValidateLuhn(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "Valid Luhn number",
			number: "4532015112830366",
			want:   true,
		},
		{
			name:   "Invalid Luhn number",
			number: "4532015112830367",
			want:   false,
		},
		{
			name:   "Valid Luhn number with odd length",
			number: "79927398713",
			want:   true,
		},
		{
			name:   "Invalid Luhn number with odd length",
			number: "79927398714",
			want:   false,
		},
		{
			name:   "Empty string",
			number: "",
			want:   false,
		},
		{
			name:   "Non-numeric characters",
			number: "4532a15112830366",
			want:   false,
		},
		{
			name:   "Single digit valid",
			number: "0",
			want:   true,
		},
		{
			name:   "Single digit invalid",
			number: "1",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLuhn(tt.number); got != tt.want {
				t.Errorf("ValidateLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
