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
	TimeB    int64     `json:"tim_bedin"`
	TimeW    int64     `json:"tim_wake"`
	Duration float64   `json:"duration"`
}

type Evaluation int

const (
	SuperBad Evaluation = iota
	Bad
	Good
	Perfect
)

type ResponseContent struct {
	Date     string `json:"date"`
	TimeB    string `json:"time_bedin"`
	TimeW    string `json:"time_wake"`
	Duration string `json:"duration"`
	Eval     string `json:"evaluation"`
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

func (c *Client) Get(past time.Time, userID string) (ResponseContents, error) {
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
		bedinTime := time.Unix(r.TimeB, 0).In(jst)
		wakeTime := time.Unix(r.TimeW, 0).In(jst)
		sleepTime := wakeTime.Sub(bedinTime).Hours()
		rc := ResponseContent{
			Date:     fmt.Sprintf("%d日", r.Date.Day()),
			TimeB:    fmt.Sprintf("%d時%d分", bedinTime.Hour(), bedinTime.Minute()),
			TimeW:    fmt.Sprintf("%d時%d分", wakeTime.Hour(), wakeTime.Minute()),
			Duration: fmt.Sprintf("%.1f時間", sleepTime),
			Eval:     c.evaluateSleep(bedinTime, wakeTime),
		}
		rcs.Record[i] = rc
		rcs.Avg += sleepTime
	}
	rcs.Avg = math.Round((rcs.Avg/float64(len(sr)))*10) / 10
	return rcs, nil
}

func (c *Client) evaluateSleep(bedin, wake time.Time) string {
	const (
		TimeLatestHourInDay = 23
		TimeEarliestBedin   = 22
		TimeAlwaysAwake     = 12
		TimeHavetoBedin     = 1
		TimeWantWake        = 7
		IdealSleepDuration  = 7
	)

	var evaluation Evaluation = Perfect

	coordinateHour := func(h int) int {
		if h >= TimeEarliestBedin {
			return h - 24
		}
		return h
	}
	hourB := coordinateHour(bedin.Hour())
	hourW := wake.Hour()

	if TimeHavetoBedin <= hourB && hourB < TimeEarliestBedin {
		evaluation -= 1
	}
	if hourW > TimeWantWake {
		evaluation -= 1
	}
	if hourW-hourB != IdealSleepDuration {
		evaluation -= 1
	}

	var result string = ""
	switch evaluation {
	case SuperBad:
		result = "😱 0:伸び代しかない!"
	case Bad:
		result = "😥 1:がんばれ!"
	case Good:
		result = "😁 2:良いね!"
	case Perfect:
		result = "🤩 3:完璧!"
	default:
		result = "🤩 エラー!"
	}
	return result
}
