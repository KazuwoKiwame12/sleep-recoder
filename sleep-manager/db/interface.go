//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE
package db

import (
	"sleep-manager/entity"
	"time"
)

type Client interface {
	SaveBedinTime(time.Time, string) error
	SaveWakeTime(time.Time, string) error
	ListInFivedays(time.Time, string) ([]entity.SleepRecord, error)
	ListInMonth(int, time.Month, string) ([]entity.SleepRecord, error)
}
