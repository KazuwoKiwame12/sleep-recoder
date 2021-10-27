package utility

import (
	"time"
)

const (
	TimeLatestHourInDay = 23
	TimeEarliestBedin   = 22
	TimeAlwaysAwake     = 12
)

func CreateStartDate(y int, m time.Month, d int) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(y, m, d, 0, 0, 0, 0, jst)
}

func CreateDateWIthJst() time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Now().In(jst)
}

func GetCorrectDayWithHour(day, hour int) int {
	if TimeLatestHourInDay >= hour && hour >= TimeEarliestBedin {
		day += 1
	}
	return day
}
