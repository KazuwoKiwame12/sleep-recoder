package usecase

import (
	"fmt"
	"math"
	"sleep-manager/db"
	"sleep-manager/entity"
	"sleep-manager/utility"
	"time"
)

type Getter struct {
	C db.Client
}

func (g *Getter) ListRecordsInFiveDays(userID string) (entity.ResponseContents, error) {
	nowJst := utility.CreateDateWIthJst()
	srs, err := g.C.ListInFivedays(nowJst, userID)
	if err != nil {
		return entity.ResponseContents{}, err
	}

	rcs := entity.ResponseContents{
		Record: make([]entity.ResponseContent, len(srs)),
	}
	for i, sr := range srs {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		bedinTime := time.Unix(sr.TimeB, 0).In(jst)
		wakeTime := time.Unix(sr.TimeW, 0).In(jst)
		sleepTime := wakeTime.Sub(bedinTime).Hours()
		rc := entity.ResponseContent{
			Date:     fmt.Sprintf("%d日", sr.Date.Day()),
			TimeB:    fmt.Sprintf("%d時%d分", bedinTime.Hour(), bedinTime.Minute()),
			TimeW:    fmt.Sprintf("%d時%d分", wakeTime.Hour(), wakeTime.Minute()),
			Duration: fmt.Sprintf("%.1f時間", sleepTime),
			Eval:     g.evaluateSleep(bedinTime, wakeTime),
		}
		rcs.Record[i] = rc
		rcs.Avg += sleepTime
	}
	rcs.Avg = math.Round((rcs.Avg/float64(len(srs)))*10) / 10
	return rcs, nil
}

func (g *Getter) ListRecordsInMonth(year int, month time.Month, userID string) (entity.ResponseContents, error) {
	srs, err := g.C.ListInMonth(year, month, userID)
	if err != nil {
		return entity.ResponseContents{}, err
	}

	rcs := entity.ResponseContents{
		Record: make([]entity.ResponseContent, len(srs)),
	}
	for i, sr := range srs {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		bedinTime := time.Unix(sr.TimeB, 0).In(jst)
		wakeTime := time.Unix(sr.TimeW, 0).In(jst)
		sleepTime := wakeTime.Sub(bedinTime).Hours()
		rc := entity.ResponseContent{
			Date:     fmt.Sprintf("%d日", sr.Date.Day()),
			TimeB:    fmt.Sprintf("%d時%d分", bedinTime.Hour(), bedinTime.Minute()),
			TimeW:    fmt.Sprintf("%d時%d分", wakeTime.Hour(), wakeTime.Minute()),
			Duration: fmt.Sprintf("%.1f時間", sleepTime),
			Eval:     g.evaluateSleep(bedinTime, wakeTime),
		}
		rcs.Record[i] = rc
		rcs.Avg += sleepTime
	}
	rcs.Avg = math.Round((rcs.Avg/float64(len(srs)))*10) / 10
	return rcs, nil
}

func (g *Getter) evaluateSleep(bedin, wake time.Time) string {
	const (
		TimeLatestHourInDay = 23
		TimeEarliestBedin   = 22
		TimeAlwaysAwake     = 12
		TimeHavetoBedin     = 1
		TimeWantWake        = 7
		IdealSleepDuration  = 7
	)

	var evaluation entity.Evaluation = entity.Perfect

	coordinateHour := func(h int) int {
		if h >= TimeEarliestBedin {
			return h - 24
		}
		return h
	}
	hourB := coordinateHour(bedin.Hour())
	hourW := wake.Hour()

	if TimeHavetoBedin <= hourB && hourB < TimeEarliestBedin {
		evaluation -= 1
	}
	if hourW > TimeWantWake {
		evaluation -= 1
	}
	if hourW-hourB != IdealSleepDuration {
		evaluation -= 1
	}

	return evaluation.ConvertToResponse()
}
