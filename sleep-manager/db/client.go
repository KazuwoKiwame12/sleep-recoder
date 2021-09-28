package db

import (
	"sleep-manager/entity"
	"sleep-manager/utility"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SleepRecordClient struct {
	Table dynamo.Table
}

func NewSleepRecordClient(tableName string, session *session.Session, config *aws.Config) *SleepRecordClient {
	db := dynamo.New(session, config)
	return &SleepRecordClient{
		Table: db.Table(tableName),
	}
}

func (s *SleepRecordClient) SaveWakeTime(now time.Time, userID string) error {
	targetDate := utility.CreateStartDate(now.Year(), now.Month(), now.Day())
	var sr entity.SleepRecord
	if err := s.Table.Get("Date", targetDate).Range("UserID", dynamo.Equal, userID).One(&sr); err != nil {
		return err
	}

	bedinTime := time.Unix(sr.TimeB, 0)
	diff := now.Sub(bedinTime).Hours()
	sr.Duration = diff
	sr.AdjustDuration()
	sr.TimeW = now.Unix()
	err := s.Table.Put(sr).Run()
	return err
}

func (s *SleepRecordClient) SaveBedinTime(now time.Time, userID string) error {
	sr := entity.SleepRecord{
		Date:   utility.CreateStartDate(now.Year(), now.Month(), utility.GetCorrectDayWithHour(now.Day(), now.Hour())),
		UserID: userID,
		TimeB:  now.Unix(),
	}
	return s.Table.Put(sr).Run()
}

func (s *SleepRecordClient) ListInFivedays(now time.Time, userID string) ([]entity.SleepRecord, error) {
	sr := make([]entity.SleepRecord, 0, 5)
	if err := s.Table.Scan().Filter("'UserID' = ? AND 'TimeW' > ?", userID, now.AddDate(0, 0, -4).Unix()).All(&sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *SleepRecordClient) ListInMonth(year int, month time.Month, userID string) ([]entity.SleepRecord, error) {
	from := utility.CreateStartDate(year, month, 1)
	to := from.AddDate(0, 1, 0)

	sr := make([]entity.SleepRecord, 0, 31)
	if err := s.Table.Scan().Filter("'UserID' = ? AND 'TimeW' > ? AND 'TimeW' <= ?", userID, from.Unix(), to.Unix()).All(&sr); err != nil {
		return nil, err
	}
	return sr, nil
}
