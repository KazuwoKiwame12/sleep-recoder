package database

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeS    time.Time `json:"tim_s"`
	TimeW    time.Time `json:"tim_w"`
	Duration float64   `json:"duration"`
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

func (c *Client) SaveWakeTime(userID string) error {
	now := time.Now().In(time.UTC)
	targetDate := createDate(now.Year(), int(now.Month()), now.Day())
	var sr SleepRecord
	if err := c.Table.Get("Date", targetDate).Range("UserID", dynamo.Equal, userID).One(&sr); err != nil {
		return err
	}

	diff := now.Sub(sr.TimeS).Hours()
	sr.Duration = diff
	sr.TimeW = now
	err := c.Table.Put(sr).Run()
	return err
}

func createDate(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}
