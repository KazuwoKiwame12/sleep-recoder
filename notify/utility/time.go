package utility

import (
	"math"
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

func CreateDateWithUnix(unix int64) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Unix(unix, 0).In(jst)
}

func GetDiffOfDays(a, b time.Time) int {
	aStartDate := CreateStartDate(a.Year(), a.Month(), a.Day())
	bStartDate := CreateStartDate(b.Year(), b.Month(), b.Day())
	return int(math.Abs(aStartDate.Sub(bStartDate).Hours() / 24.0))
}

func GetCorrectDayWithHour(day, hour int) int {
	if TimeLatestHourInDay >= hour && hour >= TimeEarliestBedin {
		day += 1
	}
	return day
}
