package db_test

import (
	"notify/db"
	"notify/entity"
	"notify/utility"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	testDate time.Time = utility.CreateStartDate(2021, 10, 15) //2021年10月15日0時0分0秒
)

func TestListInWeek(t *testing.T) {
	setAWSCredentials(t)

	c := getClient()
	userIDs := []string{"1", "2"}
	var (
		weekDays      []time.Time = make([]time.Time, 7)
		wantUserIDOne             = make(entity.SleepRecords, 7)
		wantUserIDTwo             = make(entity.SleepRecords, 7)
		wants                     = []entity.SleepRecords{wantUserIDOne, wantUserIDTwo}
	)

	for i := 0; i < len(weekDays); i++ {
		date := testDate.AddDate(0, 0, -6+(1*i))
		weekDays[i] = utility.CreateStartDate(date.Year(), date.Month(), utility.GetCorrectDayWithHour(date.Day(), date.Hour()))
		for j, id := range userIDs {
			item := entity.SleepRecord{
				Date:     weekDays[i],
				UserID:   id,
				TimeB:    date.Unix(),
				TimeW:    date.Add(7 * time.Hour).Unix(),
				Duration: (7 * time.Hour).Hours(),
			}
			item.AdjustDuration()
			if err := c.Table.Put(item).Run(); err != nil {
				t.Error(err)
			}
			wants[j][i] = item
		}
	}

	t.Cleanup(func() {
		for _, date := range weekDays {
			for _, id := range userIDs {
				if err := c.Table.Delete("UserID", id).Range("Date", date).Run(); err != nil {
					t.Error(err)
				}
			}
		}
	})

	tests := []struct {
		name  string
		input struct {
			now time.Time
		}
		want []entity.SleepRecords
	}{
		{
			name: "success",
			input: struct {
				now time.Time
			}{
				now: testDate,
			},
			want: wants,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results, err := c.ListInWeekForAllUser(test.input.now)
			if err != nil {
				t.Error(err)
			}
			if len(results) == 0 {
				t.Error("no results")
			}
			// userID=1
			for i, result := range results[:7] {
				if !isSameSleepRecord(result, test.want[0][i], t) {
					t.Errorf("unmatched error for userID=1: result[%d] is %v, want[%d] is %v", i, result, i, test.want[0][i])
				}
			}
			// userID=2
			for i, result := range results[7:] {
				if !isSameSleepRecord(result, test.want[1][i], t) {
					t.Errorf("unmatched error for userID=2: result[%d] is %v, want[%d] is %v", i, result, i, test.want[1][i])
				}
			}
		})
	}
}

func getClient() *db.SleepRecordClient {
	endPoint := "http://localhost:8000"
	tableName := "SleepRecord"
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3").WithEndpoint(endPoint)
	return db.NewSleepRecordClient(tableName, sess, config)
}

func setAWSCredentials(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "hoge")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "hoge")
	t.Setenv("AWS_DEFAULT_REGION", "ap-northeast-3")
}

func isSameSleepRecord(a, b entity.SleepRecord, t *testing.T) bool {
	if !a.Date.Equal(b.Date) {
		t.Log("Date diff")
		return false
	}
	if a.UserID != b.UserID {
		t.Log("UserID diff")
		return false
	}
	if a.TimeB != b.TimeB {
		t.Log("TimeB diff")
		return false
	}
	if a.TimeW != b.TimeW {
		t.Log("TimeW diff")
		return false
	}
	if a.Duration != b.Duration {
		t.Log("Duration diff")
		return false
	}
	return true
}
