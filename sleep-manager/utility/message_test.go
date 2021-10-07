package utility_test

import (
	"sleep-manager/utility"
	"testing"
)

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  utility.Command
	}{
		{
			name:  "bedin command",
			input: "眠る",
			want:  utility.CommandBedin,
		},
		{
			name:  "wake command",
			input: "起きた",
			want:  utility.CommandWake,
		},
		{
			name:  "fivedays command",
			input: "取得",
			want:  utility.CommandFiveDays,
		},
		{
			name:  "month command",
			input: "取得 2021 11",
			want:  utility.CommandMonth,
		},
		{
			name:  "default command",
			input: "ホゲホゲ",
			want:  utility.CommandDefault,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utility.ValidateCommand(test.input)
			if result != test.want {
				t.Errorf("unmatched error: result is %v, want is %v", result, test.want)
			}
		})
	}
}
