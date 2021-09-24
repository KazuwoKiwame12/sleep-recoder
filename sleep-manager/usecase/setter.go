package usecase

import (
	"sleep-manager/db"
	"sleep-manager/utility"
)

type Setter struct {
	C db.Client
}

func (s *Setter) SaveBedinTime(userID string) error {
	return s.C.SaveBedinTime(utility.CreateDateWIthJst(), userID)
}

func (s *Setter) SaveWakeTime(userID string) error {
	return s.C.SaveWakeTime(utility.CreateDateWIthJst(), userID)
}
