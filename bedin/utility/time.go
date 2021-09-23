package utility

import (
	"errors"
	"time"
)

const (
	TimeLatestHourInDay = 23
	TimeEarliestBedin   = 22
	TimeAlwaysAwake     = 12
)

func ValidateBedintime(hour int) error {
	if TimeEarliestBedin > hour && hour >= TimeAlwaysAwake {
		return errors.New("TimeError: you don't sleep in p.m 0 ~ p.m 9")
	}
	return nil
}

func GetCorrectDayWithHour(day, hour int) int {
	if TimeLatestHourInDay >= hour && hour >= TimeEarliestBedin {
		day += 1
	}
	return day
}

func CreateDate(y, m, d int) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, jst)
}
