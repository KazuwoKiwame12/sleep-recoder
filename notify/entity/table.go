package entity

import (
	"math"
	"sort"
	"time"
)

type SleepRecord struct {
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	TimeB    int64     `json:"time_bedin"`
	TimeW    int64     `json:"time_wake"`
	Duration float64   `json:"duration"`
}

func (s *SleepRecord) AdjustDuration() {
	s.Duration = math.Round(s.Duration*10) / 10
}

type SleepRecords []SleepRecord

func (ss SleepRecords) RetrieveUserIDs() []string {
	set := map[string]struct{}{}
	for _, s := range ss {
		if _, ok := set[s.UserID]; ok {
			continue
		}
		set[s.UserID] = struct{}{}
	}
	ids := make([]string, len(set))
	i := 0
	for id, _ := range set {
		ids[i] = id
		i++
	}
	sort.SliceStable(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}
