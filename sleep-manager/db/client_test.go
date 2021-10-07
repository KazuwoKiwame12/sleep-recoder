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

func TestSaveWakeTime(t *testing.T) {
	c := getClient()
	wakeTime := utility.CreateDateWIthJst()
	bedinTime := wakeTime.Add(-8 * time.Hour)
	userID := "sample"
	date := utility.CreateStartDate(wakeTime.Year(), wakeTime.Month(), utility.GetCorrectDayWithHour(wakeTime.Day(), wakeTime.Hour()))
	if err := c.Table.Put(
		entity.SleepRecord{
			Date:   date,
			UserID: userID,
			TimeB:  bedinTime.Unix(),
		},
	).Run(); err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if err := c.Table.Delete("Date", date).Range("UserID", userID).Run(); err != nil {
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
			name: "failed to get targetdate data",
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
			date: date,
			input: struct {
				now    time.Time
				userID string
			}{
				now:    wakeTime,
				userID: userID,
			},
			want: entity.SleepRecord{
				Date:     date,
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
				if test.want == empty {
					return
				}
				t.Error(err)
			}
			var sr entity.SleepRecord
			if err := c.Table.Get("Date", test.date).Range("UserID", dynamo.Equal, test.input.userID).One(&sr); err != nil {
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
	bedinTime := utility.CreateDateWIthJst()
	userID := "sample"
	date := utility.CreateStartDate(bedinTime.Year(), bedinTime.Month(), utility.GetCorrectDayWithHour(bedinTime.Day(), bedinTime.Hour()))
	t.Cleanup(func() {
		if err := c.Table.Delete("Date", date).Range("UserID", userID).Run(); err != nil {
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
			date: date,
			input: struct {
				now    time.Time
				userID string
			}{
				now:    bedinTime,
				userID: userID,
			},
			want: entity.SleepRecord{
				Date:   date,
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
			if err := c.Table.Get("Date", test.date).Range("UserID", dynamo.Equal, test.input.userID).One(&sr); err != nil {
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
	now := utility.CreateDateWIthJst()
	userID := "sample"
	var sixDays []time.Time = make([]time.Time, 6)
	var want []entity.SleepRecord = make([]entity.SleepRecord, 5)
	for i := 0; i < len(sixDays); i++ {
		date := now.AddDate(0, 0, -1*i)
		sixDays[i] = utility.CreateStartDate(date.Year(), date.Month(), utility.GetCorrectDayWithHour(date.Day(), date.Hour()))
		item := entity.SleepRecord{
			Date:     sixDays[i],
			UserID:   userID,
			TimeB:    date.Add(-7 * time.Hour).Unix(),
			TimeW:    date.Unix(),
			Duration: date.Sub(date.Add(-7 * time.Hour)).Hours(),
		}
		item.AdjustDuration()
		if err := c.Table.Put(item).Run(); err != nil {
			t.Error(err)
		}
		if i != len(sixDays)-1 {
			want[i] = item
		}
	}

	t.Cleanup(func() {
		for _, date := range sixDays {
			if err := c.Table.Delete("Date", date).Range("UserID", userID).Run(); err != nil {
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
				now:    now,
				userID: userID,
			},
			want: want,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srs, err := c.ListInFivedays(test.input.now, test.input.userID)
			if err != nil {
				t.Error(err)
			}
			for i, sr := range srs {
				if !isSameSleepRecord(sr, test.want[i], t) {
					t.Errorf("error: get-data is %v, but want is %v", sr, test.want)
				}
			}
		})
	}
}

func TestListInMonth(t *testing.T) {
	c := getClient()
	now := time.Now()
	from := utility.CreateStartDate(now.Year(), now.Month(), 1)
	to := from.AddDate(0, 1, 0)
	userID := "sample"
	var want []entity.SleepRecord = make([]entity.SleepRecord, int(to.Sub(from).Hours())/24)
	for i := from; i.Before(to); i.AddDate(0, 0, 1) {
		item := entity.SleepRecord{
			Date:     i,
			UserID:   userID,
			TimeB:    i.Add(-7 * time.Hour).Unix(),
			TimeW:    i.Unix(),
			Duration: i.Sub(i.Add(-7 * time.Hour)).Hours(),
		}
		item.AdjustDuration()
		if err := c.Table.Put(item).Run(); err != nil {
			t.Error(err)
		}
		want[i.Day()-1] = item
	}

	t.Cleanup(func() {
		for i := from; i.Before(to); i.AddDate(0, 0, 1) {
			if err := c.Table.Delete("Date", i).Range("UserID", userID).Run(); err != nil {
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
				now:    now,
				userID: userID,
			},
			want: want,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srs, err := c.ListInMonth(test.input.now.Year(), test.input.now.Month(), test.input.userID)
			if err != nil {
				t.Error(err)
			}
			for i, sr := range srs {
				if !isSameSleepRecord(sr, test.want[i], t) {
					t.Errorf("error: get-data is %v, but want is %v", sr, test.want)
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
