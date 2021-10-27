package db_test

import (
	"notify/db"
	"notify/entity"
	"notify/utility"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	testDate time.Time = utility.CreateStartDate(2021, 10, 15) //2021年10月15日0時0分0秒
)

func TestListInWeek(t *testing.T) {
	c := getClient()
	userID := "sample"
	var (
		weekDays []time.Time          = make([]time.Time, 5)
		want     []entity.SleepRecord = make([]entity.SleepRecord, 5) //降順データ
	)
	for i := 0; i < len(weekDays); i++ {
		date := testDate.AddDate(0, 0, -1*i)
		weekDays[i] = utility.CreateStartDate(date.Year(), date.Month(), utility.GetCorrectDayWithHour(date.Day(), date.Hour()))
		item := entity.SleepRecord{
			Date:     weekDays[i],
			UserID:   userID,
			TimeB:    date.Unix(),
			TimeW:    date.Add(7 * time.Hour).Unix(),
			Duration: (7 * time.Hour).Hours(),
		}
		item.AdjustDuration()
		if err := c.Table.Put(item).Run(); err != nil {
			t.Error(err)
		}
		want[i] = item
	}

	t.Cleanup(func() {
		for _, date := range weekDays {
			if err := c.Table.Delete("UserID", userID).Range("Date", date).Run(); err != nil {
				t.Error(err)
			}
		}
	})

	tests := []struct {
		name  string
		input struct {
			now    time.Time
			userID string
		}
		want []entity.SleepRecord
	}{
		{
			name: "success",
			input: struct {
				now    time.Time
				userID string
			}{
				now:    testDate,
				userID: userID,
			},
			want: want,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results, err := c.ListInWeek(test.input.now, test.input.userID)
			if err != nil {
				t.Error(err)
			}
			if len(results) == 0 {
				t.Error("no results")
			}
			for i, result := range results {
				if !isSameSleepRecord(result, test.want[i], t) {
					t.Errorf("unmatched error: result[%d] is %v, want[%d] is %v", i, result, i, test.want[i])
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
