package database

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeS    int64     `json:"tim_s"`
	TimeW    int64     `json:"tim_w"`
	Duration float64   `json:"duration"`
}

type Evaluation int

const (
	SUPERBAD Evaluation = iota
	BAD
	GOOD
	PERFECT
)

type ResponseContent struct {
	Date     string     `json:"date"`
	TimeS    string     `json:"time_s"`
	TimeW    string     `json:"time_w"`
	Duration string     `json:"duration"`
	Eval     Evaluation `json:"evaluation"`
}

type ResponseContents struct {
	Record []ResponseContent `json:"record"`
	Avg    float64           `json:"average"`
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

func (c *Client) GetFIveDaysRecord(past time.Time, userID string) (ResponseContents, error) {
	sr := make([]SleepRecord, 0, 5)
	if err := c.Table.Scan().Filter("'UserID' = ? AND 'TimeW' > ?", userID, past.Unix()).All(&sr); err != nil {
		return ResponseContents{}, err
	}

	rcs := ResponseContents{
		Record: make([]ResponseContent, len(sr)),
	}
	for i, r := range sr {
		if len(r.UserID) == 0 {
			break
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		timeS := time.Unix(r.TimeS, 0).In(jst)
		timeW := time.Unix(r.TimeW, 0).In(jst)
		sleepTime := timeW.Sub(timeS).Hours()
		rc := ResponseContent{
			Date:     fmt.Sprintf("%d日", r.Date.Day()),
			TimeS:    fmt.Sprintf("%d時%d分", timeS.Hour(), timeS.Minute()),
			TimeW:    fmt.Sprintf("%d時%d分", timeW.Hour(), timeW.Minute()),
			Duration: fmt.Sprintf("%.1f時間", sleepTime),
			Eval:     c.evaluateSleep(timeS, timeW),
		}
		rcs.Record[i] = rc
		rcs.Avg += sleepTime
	}
	rcs.Avg = math.Round((rcs.Avg/float64(len(sr)))*10) / 10
	return rcs, nil
}

func (c *Client) evaluateSleep(bedin, awake time.Time) Evaluation {
	var negativePoint Evaluation = 0

	coordinateHour := func(h int) int {
		if h >= 22 {
			return h - 24
		}
		return h
	}
	hourS := coordinateHour(bedin.Hour())
	hourW := awake.Hour()

	if 22 > hourS && hourS >= 1 {
		negativePoint += 1
	}
	if hourW > 7 {
		negativePoint += 1
	}
	if hourW-hourS != 7 {
		negativePoint += 1
	}

	return PERFECT - negativePoint
}
