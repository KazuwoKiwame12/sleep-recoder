package usecase_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"sleep-manager/entity"
	mock_db "sleep-manager/mock/db"
	"sleep-manager/usecase"
	"sleep-manager/utility"

	"github.com/golang/mock/gomock"
)

func TestListRecordsInFiveDays(t *testing.T) {
	userID := "sample"
	nowJst := utility.CreateDateWIthJst()
	wakeTime := time.Date(nowJst.Year(), nowJst.Month(), nowJst.Day(), 7, 0, 0, 0, nowJst.Location())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		setup func() *mock_db.MockClient
		want  entity.ResponseContents
	}{
		{
			name: "success",
			setup: func() *mock_db.MockClient {
				mock := mock_db.NewMockClient(ctrl)
				mock.EXPECT().ListInFivedays(gomock.Any(), userID).Return(
					[]entity.SleepRecord{
						{
							Date:     utility.CreateStartDate(wakeTime.Year(), wakeTime.Month(), wakeTime.Day()),
							TimeB:    wakeTime.Add(-8 * time.Hour).Unix(),
							TimeW:    wakeTime.Unix(),
							Duration: wakeTime.Sub(wakeTime.Add(-8 * time.Hour)).Hours(),
						},
						{
							Date:     utility.CreateStartDate(wakeTime.Year(), wakeTime.Month(), wakeTime.Day()).AddDate(0, 0, -1),
							TimeB:    wakeTime.AddDate(0, 0, -1).Add(-8 * time.Hour).Unix(),
							TimeW:    wakeTime.AddDate(0, 0, -1).Unix(),
							Duration: wakeTime.AddDate(0, 0, -1).Sub(wakeTime.AddDate(0, 0, -1).Add(-8 * time.Hour)).Hours(),
						},
					},
					nil,
				)
				return mock
			},
			want: entity.ResponseContents{
				Record: []entity.ResponseContent{
					{
						Date:     fmt.Sprintf("%d日", wakeTime.Day()),
						TimeB:    fmt.Sprintf("%d時%d分", wakeTime.Add(-8*time.Hour).Hour(), wakeTime.Add(-8*time.Hour).Minute()),
						TimeW:    fmt.Sprintf("%d時%d分", wakeTime.Hour(), wakeTime.Minute()),
						Duration: "8.0時間",
						Eval:     "😁 2:良いね!",
					},
					{
						Date:     fmt.Sprintf("%d日", wakeTime.AddDate(0, 0, -1).Day()),
						TimeB:    fmt.Sprintf("%d時%d分", wakeTime.AddDate(0, 0, -1).Add(-8*time.Hour).Hour(), wakeTime.AddDate(0, 0, -1).Add(-8*time.Hour).Minute()),
						TimeW:    fmt.Sprintf("%d時%d分", wakeTime.AddDate(0, 0, -1).Hour(), wakeTime.AddDate(0, 0, -1).Minute()),
						Duration: "8.0時間",
						Eval:     "😁 2:良いね!",
					},
				},
				Avg: 8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.setup()
			g := &usecase.Getter{C: mock}
			result, err := g.ListRecordsInFiveDays(userID)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("unmatched error: result is %+v, want is %+v", result, test.want)
			}
		})
	}
}

func TestListRecordsInMonth(t *testing.T) {
	userID := "sample"
	nowJst := utility.CreateDateWIthJst()
	wakeTime := time.Date(nowJst.Year(), nowJst.Month(), nowJst.Day(), 7, 0, 0, 0, nowJst.Location())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		setup func() *mock_db.MockClient
		want  entity.ResponseContents
	}{
		{
			name: "success",
			setup: func() *mock_db.MockClient {
				mock := mock_db.NewMockClient(ctrl)
				mock.EXPECT().ListInMonth(gomock.Any(), gomock.Any(), userID).Return(
					[]entity.SleepRecord{
						{
							Date:     utility.CreateStartDate(wakeTime.Year(), wakeTime.Month(), wakeTime.Day()),
							TimeB:    wakeTime.Add(-8 * time.Hour).Unix(),
							TimeW:    wakeTime.Unix(),
							Duration: wakeTime.Sub(wakeTime.Add(-8 * time.Hour)).Hours(),
						},
						{
							Date:     utility.CreateStartDate(wakeTime.Year(), wakeTime.Month(), wakeTime.Day()).AddDate(0, 0, -1),
							TimeB:    wakeTime.AddDate(0, 0, -1).Add(-8 * time.Hour).Unix(),
							TimeW:    wakeTime.AddDate(0, 0, -1).Unix(),
							Duration: wakeTime.AddDate(0, 0, -1).Sub(wakeTime.AddDate(0, 0, -1).Add(-8 * time.Hour)).Hours(),
						},
					},
					nil,
				)
				return mock
			},
			want: entity.ResponseContents{
				Record: []entity.ResponseContent{
					{
						Date:     fmt.Sprintf("%d日", wakeTime.Day()),
						TimeB:    fmt.Sprintf("%d時%d分", wakeTime.Add(-8*time.Hour).Hour(), wakeTime.Add(-8*time.Hour).Minute()),
						TimeW:    fmt.Sprintf("%d時%d分", wakeTime.Hour(), wakeTime.Minute()),
						Duration: "8.0時間",
						Eval:     "😁 2:良いね!",
					},
					{
						Date:     fmt.Sprintf("%d日", wakeTime.AddDate(0, 0, -1).Day()),
						TimeB:    fmt.Sprintf("%d時%d分", wakeTime.AddDate(0, 0, -1).Add(-8*time.Hour).Hour(), wakeTime.AddDate(0, 0, -1).Add(-8*time.Hour).Minute()),
						TimeW:    fmt.Sprintf("%d時%d分", wakeTime.AddDate(0, 0, -1).Hour(), wakeTime.AddDate(0, 0, -1).Minute()),
						Duration: "8.0時間",
						Eval:     "😁 2:良いね!",
					},
				},
				Avg: 8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.setup()
			g := &usecase.Getter{C: mock}
			result, err := g.ListRecordsInMonth(nowJst.Year(), nowJst.Month(), userID)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("unmatched error: result is %+v, want is %+v", result, test.want)
			}
		})
	}
}

func TestEvaluateSleep(t *testing.T) {
	type input struct {
		bedinTime time.Time
		wakeTime  time.Time
	}

	tests := []struct {
		name  string
		input struct {
			bedinTime time.Time
			wakeTime  time.Time
		}
		want entity.Evaluation
	}{
		{
			name: "perfect case",
			input: input{
				bedinTime: createDateWithHour(7).Add(-7 * time.Hour),
				wakeTime:  createDateWithHour(7),
			},
			want: entity.Perfect,
		},
		{
			name: "good case",
			input: input{
				bedinTime: createDateWithHour(7).Add(-8 * time.Hour),
				wakeTime:  createDateWithHour(7),
			},
			want: entity.Good,
		},
		{
			name: "bad case 1",
			input: input{
				bedinTime: createDateWithHour(1),
				wakeTime:  createDateWithHour(8),
			},
			want: entity.Bad,
		},
		{
			name: "bad case 2",
			input: input{
				bedinTime: createDateWithHour(0),
				wakeTime:  createDateWithHour(9),
			},
			want: entity.Bad,
		},
		{
			name: "super bad case",
			input: input{
				bedinTime: createDateWithHour(1),
				wakeTime:  createDateWithHour(10),
			},
			want: entity.SuperBad,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := usecase.ExEvaluateSleep(nil, test.input.bedinTime, test.input.wakeTime)
			if result != test.want.ConvertToResponse() {
				t.Errorf("unmatched error: result is %s, want is %s", result, test.want.ConvertToResponse())
			}
		})
	}
}

func createDateWithHour(hour int) time.Time {
	n := utility.CreateDateWIthJst()
	return time.Date(n.Year(), n.Month(), n.Day(), hour, 0, 0, 0, n.Location())
}
