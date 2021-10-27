package db

import (
	"notify/entity"
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

func (s *SleepRecordClient) ListInWeek(now time.Time, userID string) ([]entity.SleepRecord, error) {
	sr := make([]entity.SleepRecord, 0, 7)
	if err := s.Table.Get("UserID", userID).Filter("'TimeW' > ?", now.AddDate(0, 0, -6).Unix()).Order(dynamo.Descending).All(&sr); err != nil {
		return nil, err
	}
	return sr, nil
}
