package utility_test

import (
	"notify/utility"
	"testing"
)

func TestGetCorrectDayWithHour(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			day  int
			hour int
		}
		want int
	}{
		{
			name: "plus one day",
			input: struct {
				day  int
				hour int
			}{
				day:  5,
				hour: 22,
			},
			want: 6,
		},
		{
			name: "no change",
			input: struct {
				day  int
				hour int
			}{
				day:  5,
				hour: 0,
			},
			want: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utility.GetCorrectDayWithHour(test.input.day, test.input.hour)
			if result != test.want {
				t.Errorf("unmatched error: result is %d, want is %d", result, test.want)
			}
		})
	}
}
