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

	data := []struct {
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

	for i, d := range data {
		t.Run(fmt.Sprintf("%d: %s", i, d.name), func(t *testing.T) {
			if err := c.SaveWakeTime(d.input.now, d.input.userID); err != nil {
				empty := entity.SleepRecord{}
				if d.want == empty {
					return
				}
				t.Error(err)
			}
			var sr entity.SleepRecord
			if err := c.Table.Get("Date", d.date).Range("UserID", dynamo.Equal, d.input.userID).One(&sr); err != nil {
				t.Error(err)
			}
			if !isSameSleepRecord(sr, d.want, t) {
				t.Errorf("error: get-data is %v, but want is %v", sr, d.want)
			}
		})
	}

	if err := c.Table.Delete("Date", date).Range("UserID", userID).Run(); err != nil {
		t.Error(err)
	}
}

func TestSaveBedinTime(t *testing.T) {
	c := getClient()
	bedinTime := utility.CreateDateWIthJst()
	userID := "sample"
	date := utility.CreateStartDate(bedinTime.Year(), bedinTime.Month(), utility.GetCorrectDayWithHour(bedinTime.Day(), bedinTime.Hour()))

	data := []struct {
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

	for i, d := range data {
		t.Run(fmt.Sprintf("%d: %s", i, d.name), func(t *testing.T) {
			if err := c.SaveBedinTime(d.input.now, d.input.userID); err != nil {
				t.Error(err)
			}
			var sr entity.SleepRecord
			if err := c.Table.Get("Date", d.date).Range("UserID", dynamo.Equal, d.input.userID).One(&sr); err != nil {
				t.Error(err)
			}
			if !isSameSleepRecord(sr, d.want, t) {
				t.Errorf("error: get-data is %v, but want is %v", sr, d.want)
			}
		})
	}

	if err := c.Table.Delete("Date", date).Range("UserID", userID).Run(); err != nil {
		t.Error(err)
	}
}

func TestListInFivedays(t *testing.T) {

}

func TestListInMonth(t *testing.T) {

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
