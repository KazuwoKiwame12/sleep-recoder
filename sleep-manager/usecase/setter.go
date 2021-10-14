package usecase

import (
	"os"
	"sleep-manager/db"
	"sleep-manager/utility"
	"time"
)

var TestNow time.Time // テスト用の時刻変数

type Setter struct {
	C db.Client
}

func (s *Setter) SaveBedinTime(userID string) error {
	now := utility.CreateDateWIthJst()
	if len(os.Getenv("IS_TEST")) != 0 {
		now = TestNow
	}
	if err := utility.ValidateBedintime(now.Hour()); err != nil {
		return err
	}
	return s.C.SaveBedinTime(now, userID)
}

func (s *Setter) SaveWakeTime(userID string) error {
	return s.C.SaveWakeTime(utility.CreateDateWIthJst(), userID)
}
