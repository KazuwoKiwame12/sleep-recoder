//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE
package db

import (
	"notify/entity"
	"time"
)

type Client interface {
	ListInWeek(time.Time, string) (entity.SleepRecords, error)
}
