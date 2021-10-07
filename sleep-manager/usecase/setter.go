package usecase

import (
	"sleep-manager/db"
	"sleep-manager/utility"
)

type Setter struct {
	C db.Client
}

func (s *Setter) SaveBedinTime(userID string) error {
	now := utility.CreateDateWIthJst()
	if err := utility.ValidateBedintime(now.Hour()); err != nil {
		return err
	}
	return s.C.SaveBedinTime(now, userID)
}

func (s *Setter) SaveWakeTime(userID string) error {
	return s.C.SaveWakeTime(utility.CreateDateWIthJst(), userID)
}
