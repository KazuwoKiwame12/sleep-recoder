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

func (c *Client) SaveWakeTime(now time.Time, userID string) error {
	targetDate := c.createDate(now.Year(), int(now.Month()), now.Day())
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

func (c *Client) createDate(y, m, d int) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, jst)
}
