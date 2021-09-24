package entity

import "time"

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeB    int64     `json:"time_bedin"`
	TimeW    int64     `json:"time_wake"`
	Duration float64   `json:"duration"`
}
