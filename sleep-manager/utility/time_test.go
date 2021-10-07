package utility_test

import (
	"sleep-manager/utility"
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

func TestValidateBedintime(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{
			name:    "failed to sleep in p.m 0 ~ p.m 9",
			input:   13,
			wantErr: true,
		},
		{
			name:  "success",
			input: 22,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utility.ValidateBedintime(test.input)
			result := err != nil
			if result != test.wantErr {
				t.Errorf("unmatched error: result is %v, want is %v", result, test.wantErr)
			}
		})
	}
}
