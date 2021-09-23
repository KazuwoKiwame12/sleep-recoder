package database

import (
	"bedin-time/utility"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeB    int64     `json:"time_bedin"`
	TimeA    int64     `json:"time_awake"`
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

func (c *Client) Save(now time.Time, userID string) error {
	sr := SleepRecord{
		Date:   utility.CreateDate(now.Year(), int(now.Month()), utility.GetCorrectDayWithHour(now.Day(), now.Hour())),
		UserID: userID,
		TimeB:  now.Unix(),
	}
	return c.Table.Put(sr).Run()
}
