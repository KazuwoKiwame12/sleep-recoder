package db_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sleep-manager/db"
	"sleep-manager/entity"
	"sleep-manager/utility"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
)

var (
	testDate time.Time = utility.CreateStartDate(2021, 10, 15) //2021年10月15日0時0分0秒

)

func TestSaveWakeTime(t *testing.T) {
	c := getClient()
	bedinTime := testDate
	wakeTime := bedinTime.Add(8 * time.Hour)
	userID := "sample"
	if err := c.Table.Put(
		entity.SleepRecord{
			Date:   testDate,
			UserID: userID,
			TimeB:  bedinTime.Unix(),
		},
	).Run(); err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if err := c.Table.Delete("UserID", userID).Range("Date", testDate).Run(); err != nil {
			t.Error(err)
		}
	})

	tests := []struct {
		name  string
		date  time.Time
		input struct {
			now    time.Time
			userID string
		}
		want entity.SleepRecord
	}{
		{
			name: "failed to get bedinTime of testDate",
			input: struct {
				now    time.Time
				userID string
			}{
				now:    wakeTime.AddDate(0, 0, -2),
				userID: userID,
			},
		},
		{
			name: "success",
			date: testDate,
			input: struct {
				now    time.Time
				userID string
			}{
				now:    wakeTime,
				userID: userID,
			},
			want: entity.SleepRecord{
				Date:     testDate,
				UserID:   userID,
				TimeB:    bedinTime.Unix(),
				TimeW:    wakeTime.Unix(),
				Duration: wakeTime.Sub(bedinTime).Hours(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := c.SaveWakeTime(test.input.now, test.input.userID); err != nil {
				empty := entity.SleepRecord{}
				if test.want == empty { // 睡眠時刻が取得できないケースのテスト
					return
				}
				t.Error(err)
			}
			var sr entity.SleepRecord
			if err := c.Table.Get("UserID", test.input.userID).Range("Date", dynamo.Equal, test.date).One(&sr); err != nil {
				t.Error(err)
			}
			if !isSameSleepRecord(sr, test.want, t) {
				t.Errorf("error: get-data is %v, but want is %v", sr, test.want)
			}
		})
	}
}

func TestSaveBedinTime(t *testing.T) {
	c := getClient()
	userID := "sample"
	bedinTime := testDate.Add(time.Hour)
	t.Cleanup(func() {
		if err := c.Table.Delete("UserID", userID).Range("Date", testDate).Run(); err != nil {
			t.Error(err)
		}
	})

	tests := []struct {
		name  string
		date  time.Time
		input struct {
			now    time.Time
			userID string
		}
		want entity.SleepRecord
	}{
		{
			name: "success",
			date: testDate,
			input: struct {
				now    time.Time
				userID string
			}{
				now:    bedinTime,
				userID: userID,
			},
			want: entity.SleepRecord{
				Date:   testDate,
				UserID: userID,
				TimeB:  bedinTime.Unix(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := c.SaveBedinTime(test.input.now, test.input.userID); err != nil {
				t.Error(err)
			}
			var sr entity.SleepRecord
			if err := c.Table.Get("UserID", test.input.userID).Range("Date", dynamo.Equal, test.date).One(&sr); err != nil {
				t.Error(err)
			}
			if !isSameSleepRecord(sr, test.want, t) {
				t.Errorf("error: get-data is %v, but want is %v", sr, test.want)
			}
		})
	}
}

func TestListInFivedays(t *testing.T) {
	c := getClient()
	userID := "sample"
	var (
		fiveDays []time.Time          = make([]time.Time, 5)
		want     []entity.SleepRecord = make([]entity.SleepRecord, 5) //降順データ
	)
	for i := 0; i < len(fiveDays); i++ {
		date := testDate.AddDate(0, 0, -1*i)
		fiveDays[i] = utility.CreateStartDate(date.Year(), date.Month(), utility.GetCorrectDayWithHour(date.Day(), date.Hour()))
		item := entity.SleepRecord{
			Date:     fiveDays[i],
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
		for _, date := range fiveDays {
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
			results, err := c.ListInFivedays(test.input.now, test.input.userID)
			if err != nil {
				t.Error(err)
			}
			for i, result := range results {
				if !isSameSleepRecord(result, test.want[i], t) {
					t.Errorf("unmatched error: result[%d] is %v, want[%d] is %v", i, result, i, test.want[i])
				}
			}
		})
	}
}

func TestListInMonth(t *testing.T) {
	c := getClient()
	from := utility.CreateStartDate(testDate.Year(), testDate.Month(), 1)
	to := from.AddDate(0, 1, 0)
	userID := "sample"
	var (
		thirtyDays []time.Time          = make([]time.Time, 30)
		want       []entity.SleepRecord = make([]entity.SleepRecord, 30) //降順データ
	)
	for i := 0; i < len(thirtyDays); i++ {
		date := to.AddDate(0, 0, -1*i)
		thirtyDays[i] = utility.CreateStartDate(date.Year(), date.Month(), utility.GetCorrectDayWithHour(date.Day(), date.Hour()))
		item := entity.SleepRecord{
			Date:     thirtyDays[i],
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
		for _, date := range thirtyDays {
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
			results, err := c.ListInMonth(test.input.now.Year(), test.input.now.Month(), test.input.userID)
			if err != nil {
				t.Error(err)
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
	err := godotenv.Load(fmt.Sprintf("%s/%s", filepath.Dir("../.env"), ".env"))
	if err != nil {
		log.Fatal("error: loading .env file")
	}
	endPoint := os.Getenv("DYNAMODB_ENDPOINT")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
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
