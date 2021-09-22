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

func (c *Client) SaveSleepTime(now time.Time, userID string) error {
	sr := SleepRecord{
		Date:   c.createDate(now.Year(), int(now.Month()), c.getCorrectDayWithHour(now.Day(), now.Hour())),
		UserID: userID,
		TimeS:  now,
	}
	return c.Table.Put(sr).Run()
}

func (c *Client) getCorrectDayWithHour(day, hour int) int {
	const (
		TimeStartSleep      = 22
		TimeLatestHourInDay = 23
	)

	if TimeLatestHourInDay >= hour && hour >= TimeStartSleep {
		day += 1
	}
	return day
}

func (c *Client) createDate(y, m, d int) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, jst)
}
