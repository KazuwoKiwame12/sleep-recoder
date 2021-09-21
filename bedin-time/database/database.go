package database

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const TimeStartSleep = 22

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeS    time.Time `json:"tim_s"`
	TimeW    time.Time `json:"tim_w"`
	Duration float32   `json:"duration"`
}

type Client struct {
	Table dynamo.Table
}

func NewClient(session *session.Session, config *aws.Config) *Client {
	db := dynamo.New(session, config)
	return &Client{
		Table: db.Table(os.Getenv("DYNAMODB_TABLE_NAME")),
	}
}

func (c *Client) SaveSleepTime(userID string) error {
	now := time.Now()
	var day int = now.Day()
	if now.Hour() >= TimeStartSleep {
		day += 1
	}
	sr := SleepRecord{
		Date:   createDate(now.Year(), int(now.Month()), day),
		UserID: userID,
		TimeS:  now,
	}
	err := c.Table.Put(sr).Run()
	return err
}

func createDate(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}
