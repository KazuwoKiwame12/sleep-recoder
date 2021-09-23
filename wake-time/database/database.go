package database

import (
	"os"
	"time"
	"wake-time/utility"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeB    int64     `json:"tim_bedin"`
	TimeW    int64     `json:"tim_wake"`
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
	targetDate := utility.CreateDate(now.Year(), int(now.Month()), now.Day())
	var sr SleepRecord
	if err := c.Table.Get("Date", targetDate).Range("UserID", dynamo.Equal, userID).One(&sr); err != nil {
		return err
	}

	bedinTime := time.Unix(sr.TimeB, 0)
	diff := now.Sub(bedinTime).Hours()
	sr.Duration = diff
	sr.TimeW = now.Unix()
	err := c.Table.Put(sr).Run()
	return err
}
