package entity_test

import (
	"fmt"
	"sleep-manager/entity"
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

func TestConvertToResponse(t *testing.T) {
	data := []struct {
		name  string
		input entity.Evaluation
		want  string
	}{
		{
			name:  "SuperBad",
			input: entity.SuperBad,
			want:  "😱 0:伸び代しかない!",
		},
		{
			name:  "Bad",
			input: entity.Bad,
			want:  "😥 1:がんばれ!",
		},
		{
			name:  "Good",
			input: entity.Good,
			want:  "😁 2:良いね!",
		},
		{
			name:  "Perfect",
			input: entity.Perfect,
			want:  "🤩 3:完璧!",
		},
		{
			name:  "Error",
			input: entity.Evaluation(100),
			want:  "🤩 エラー!",
		},
	}

	for i, d := range data {
		t.Run(fmt.Sprintf("%d: %s", i, d.name), func(t *testing.T) {
			result := d.input.ConvertToResponse()
			if result != d.want {
				t.Errorf("error: result = %v, want = %v", result, d.want)
			}
		})
	}
}
