package entity_test

import (
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
	tests := []struct {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.ConvertToResponse()
			if result != test.want {
				t.Errorf("error: result = %v, want = %v", result, test.want)
			}
		})
	}
}

func TestGetLineMessage(t *testing.T) {
	instance := entity.ResponseContents{
		Record: []entity.ResponseContent{
			{
				Date:     "13日",
				TimeB:    "1時1分",
				TimeW:    "8時21分",
				Duration: "7.3時間",
				Eval:     "😥 1:がんばれ!",
			},
			{
				Date:     "14日",
				TimeB:    "1時11分",
				TimeW:    "7時30分",
				Duration: "6.3時間",
				Eval:     "😥 1:がんばれ!",
			},
		},
		Avg: 6.8,
	}
	head := "****睡眠記録****\n平均睡眠時間: 6.8時間\n\n各日にちの睡眠記録\n"
	bodyOne := "【13日】: 😥 1:がんばれ!\n\t就寝: 1時1分\n\t起床: 8時21分\n\t睡眠時間: 7.3時間\n\n"
	bodyTwo := "【14日】: 😥 1:がんばれ!\n\t就寝: 1時11分\n\t起床: 7時30分\n\t睡眠時間: 6.3時間\n\n"
	want := head + bodyOne + bodyTwo
	result := instance.GetLineMessage()
	if result != want {
		t.Errorf("unmatched error:\nresult is\n%s\nwant is\n%s", result, want)
	}
}
