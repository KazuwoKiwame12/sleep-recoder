package utility

import "time"

func CreateDate(y, m, d int) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, jst)
}
