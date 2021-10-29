package entity_test

import (
	"notify/entity"
	"testing"
)

func TestAdjustDuration(t *testing.T) {
	input := 9.85
	want := 9.9
	sr := entity.SleepRecord{Duration: input}
	sr.AdjustDuration()
	if sr.Duration != want {
		t.Errorf("error: result = %v, want = %v", sr.Duration, want)
	}
}

func Test_RetrieveUserIDs(t *testing.T) {
	tests := []struct {
		name  string
		input entity.SleepRecords
		want  []string
	}{
		{
			name: "no duplicate",
			input: entity.SleepRecords{
				{UserID: "3"},
				{UserID: "1"},
				{UserID: "2"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "duplicate",
			input: entity.SleepRecords{
				{UserID: "2"},
				{UserID: "1"},
				{UserID: "2"},
			},
			want: []string{"1", "2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.RetrieveUserIDs()
			for i := 0; i < len(result); i++ {
				if result[i] != test.want[i] {
					t.Errorf("unmatched error: result[%d] is %v, want[%d] is %v", i, result[i], i, test.want[i])
				}
			}
		})
	}
}
