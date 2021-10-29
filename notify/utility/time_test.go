package utility_test

import (
	"notify/utility"
	"testing"
	"time"
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

func Test_GetDiffOfDays(t *testing.T) {
	type input struct {
		a time.Time
		b time.Time
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name  string
		input input
		want  int
	}{
		{
			name: "check 2018/10/3 10:00:00 and 2018/10/1 21:00:00",
			input: input{
				a: time.Date(2018, time.October, 3, 10, 0, 0, 0, jst),
				b: time.Date(2018, time.October, 1, 21, 0, 0, 0, jst),
			},
			want: 2,
		},
		{
			name: "check 2018/10/1 21:00:00 and 2018/10/3 10:00:00",
			input: input{
				a: time.Date(2018, time.October, 1, 21, 0, 0, 0, jst),
				b: time.Date(2018, time.October, 3, 10, 0, 0, 0, jst),
			},
			want: 2,
		},
		{
			name: "check 2018/10/1 21:00:00 and 2018/9/29 10:00:00",
			input: input{
				a: time.Date(2018, time.October, 1, 21, 0, 0, 0, jst),
				b: time.Date(2018, time.September, 29, 10, 0, 0, 0, jst),
			},
			want: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utility.GetDiffOfDays(test.input.a, test.input.b)
			if result != test.want {
				t.Errorf("unmatched error: result is %d, want is %d", result, test.want)
			}
		})
	}
}
