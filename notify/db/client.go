package db

import (
	"notify/entity"
	"sort"
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

func (s *SleepRecordClient) ListInWeekForAllUser(now time.Time) (entity.SleepRecords, error) {
	srs := entity.SleepRecords{}
	if err := s.Table.Scan().Filter("'TimeW' > ?", now.AddDate(0, 0, -6).Unix()).All(&srs); err != nil {
		return nil, err
	}
	sort.SliceStable(srs, func(i, j int) bool { return srs[i].TimeW < srs[j].TimeW })
	sort.SliceStable(srs, func(i, j int) bool { return srs[i].UserID < srs[j].UserID })
	return srs, nil
}
