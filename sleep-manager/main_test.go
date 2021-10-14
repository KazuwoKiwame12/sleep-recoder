package main

import (
	"errors"
	"math"
	"sleep-manager/entity"
	mock_db "sleep-manager/mock/db"
	"sleep-manager/usecase"
	"sleep-manager/utility"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

const (
	userID                        string = "sampleID"
	bedinCommand                  string = "眠る"
	wakeCommand                   string = "起きた"
	helpCommand                   string = "説明"
	getFiveDaysCommand            string = "取得"
	getMonthCommand               string = "取得 2021 10"
	getMonthCommandNotYear        string = "取得 sample 10"
	getMonthCommandNotNumberMonth string = "取得 2021 sample"
	getMonthCommandNotMonth       string = "取得 2021 13"
	defaultCommand                string = "sample"
)

var (
	errorAny          error  = errors.New("any")
	listRecordMessage string = "****睡眠記録****\n平均睡眠時間: 6.8時間\n\n各日にちの睡眠記録\n【13日】: 😥 がんばれ!\n\t就寝: 1時1分\n\t起床: 8時21分\n\t睡眠時間: 7.3時間\n\n【14日】: 😥 がんばれ!\n\t就寝: 1時11分\n\t起床: 7時30分\n\t睡眠時間: 6.3時間\n\n"
)

func TestExecCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		setup func() input
		want  string
	}{
		{
			name: "bedin command system error message",
			setup: func() input {
				t.Setenv("IS_TEST", "test")
				now := utility.CreateDateWIthJst()
				usecase.TestNow = time.Date(2021, 10, 10, 23, 0, 0, 0, now.Location())
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().SaveBedinTime(gomock.Any(), userID).Return(errorAny)
				s := usecase.Setter{C: mockClient}
				return input{
					command: bedinCommand,
					userID:  userID,
					setter:  s,
				}
			},
			want: utility.MessageSystemError,
		},
		{
			name: "bedin command time error message",
			setup: func() input {
				t.Setenv("IS_TEST", "test")
				now := utility.CreateDateWIthJst()
				usecase.TestNow = time.Date(2021, 10, 10, 15, 0, 0, 0, now.Location())
				return input{
					command: bedinCommand,
					userID:  userID,
				}
			},
			want: utility.MessageBedinTimeError,
		},
		{
			name: "bedin command suucess message",
			setup: func() input {
				t.Setenv("IS_TEST", "test")
				now := utility.CreateDateWIthJst()
				usecase.TestNow = time.Date(2021, 10, 10, 11, 0, 0, 0, now.Location())
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().SaveBedinTime(gomock.Any(), userID).Return(nil)
				s := usecase.Setter{C: mockClient}
				return input{
					command: bedinCommand,
					userID:  userID,
					setter:  s,
				}
			},
			want: utility.MessageSuccessRecord,
		},
		{
			name: "wake command not sleep message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().SaveWakeTime(gomock.Any(), userID).Return(errorAny)
				s := usecase.Setter{C: mockClient}
				return input{
					command: wakeCommand,
					userID:  userID,
					setter:  s,
				}
			},
			want: utility.MessageNotSleep,
		},
		{
			name: "wake command suucess message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().SaveWakeTime(gomock.Any(), userID).Return(nil)
				s := usecase.Setter{C: mockClient}
				return input{
					command: wakeCommand,
					userID:  userID,
					setter:  s,
				}
			},
			want: utility.MessageSuccessRecord,
		},
		{
			name: "fivedays command system error message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInFivedays(gomock.Any(), userID).Return([]entity.SleepRecord{}, errorAny)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getFiveDaysCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: utility.MessageSystemError,
		},
		{
			name: "fivedays command not found message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInFivedays(gomock.Any(), userID).Return([]entity.SleepRecord{}, nil)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getFiveDaysCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: utility.MessageNotFound,
		},
		{
			name: "fivedays command suucess message",
			setup: func() input {
				now := utility.CreateDateWIthJst()
				date := time.Date(2021, 10, 13, 0, 0, 0, 0, now.Location())
				timeBs := []time.Time{date.Add(time.Hour).Add(time.Minute), date.Add(time.Hour).Add(11 * time.Minute)}
				timeWs := []time.Time{date.Add(8 * time.Hour).Add(21 * time.Minute), date.Add(7 * time.Hour).Add(30 * time.Minute)}
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInFivedays(gomock.Any(), userID).Return([]entity.SleepRecord{
					{
						Date:     date,
						UserID:   userID,
						TimeB:    timeBs[0].Unix(),
						TimeW:    timeWs[0].Unix(),
						Duration: math.Round(timeWs[0].Sub(timeBs[0]).Hours()*10) / 10,
					},
					{
						Date:     date.AddDate(0, 0, 1),
						UserID:   userID,
						TimeB:    timeBs[1].Unix(),
						TimeW:    timeWs[1].Unix(),
						Duration: math.Round(timeWs[1].Sub(timeBs[1]).Hours()*10) / 10,
					},
				}, nil)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getFiveDaysCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: listRecordMessage,
		},
		{
			name: "month command not exist year message",
			setup: func() input {
				return input{
					command: getMonthCommandNotYear,
					userID:  userID,
				}
			},
			want: utility.MessageNotExistYear,
		},
		{
			name: "month command not exist month message__not number input",
			setup: func() input {
				return input{
					command: getMonthCommandNotNumberMonth,
					userID:  userID,
				}
			},
			want: utility.MessageNotExistMonth,
		},
		{
			name: "month command not exist month message__not month input",
			setup: func() input {
				return input{
					command: getMonthCommandNotMonth,
					userID:  userID,
				}
			},
			want: utility.MessageNotExistMonth,
		},
		{
			name: "month command system error message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInMonth(gomock.Any(), gomock.Any(), userID).Return([]entity.SleepRecord{}, errorAny)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getMonthCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: utility.MessageSystemError,
		},
		{
			name: "month command not found message",
			setup: func() input {
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInMonth(gomock.Any(), gomock.Any(), userID).Return([]entity.SleepRecord{}, nil)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getMonthCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: utility.MessageNotFound,
		},
		{
			name: "month command suucess message",
			setup: func() input {
				now := utility.CreateDateWIthJst()
				date := time.Date(2021, 10, 13, 0, 0, 0, 0, now.Location())
				timeBs := []time.Time{date.Add(time.Hour).Add(time.Minute), date.Add(time.Hour).Add(11 * time.Minute)}
				timeWs := []time.Time{date.Add(8 * time.Hour).Add(21 * time.Minute), date.Add(7 * time.Hour).Add(30 * time.Minute)}
				mockClient := mock_db.NewMockClient(ctrl)
				mockClient.EXPECT().ListInFivedays(gomock.Any(), userID).Return([]entity.SleepRecord{
					{
						Date:     date,
						UserID:   userID,
						TimeB:    timeBs[0].Unix(),
						TimeW:    timeWs[0].Unix(),
						Duration: math.Round(timeWs[0].Sub(timeBs[0]).Hours()*10) / 10,
					},
					{
						Date:     date.AddDate(0, 0, 1),
						UserID:   userID,
						TimeB:    timeBs[1].Unix(),
						TimeW:    timeWs[1].Unix(),
						Duration: math.Round(timeWs[1].Sub(timeBs[1]).Hours()*10) / 10,
					},
				}, nil)
				g := usecase.Getter{C: mockClient}
				return input{
					command: getFiveDaysCommand,
					userID:  userID,
					getter:  g,
				}
			},
			want: listRecordMessage,
		},
		{
			name: "help command suucess message",
			setup: func() input {
				return input{
					command: helpCommand,
					userID:  userID,
				}
			},
			want: utility.MessageHelp,
		},
		{
			name: "default command suucess message",
			setup: func() input {
				return input{
					command: defaultCommand,
					userID:  userID,
				}
			},
			want: utility.MessageDefault,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := test.setup()
			result := execCommand(i)
			if result != test.want {
				t.Errorf("unmatched error: result is %s, want is %s", result, test.want)
			}
		})
	}
}
