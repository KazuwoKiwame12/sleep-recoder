package entity

import (
	"fmt"
	"strings"
)

type Evaluation int

func (e Evaluation) ConvertToResponse() string {
	switch e {
	case SuperBad:
		return "😱 伸び代しかない!"
	case Bad:
		return "😥 がんばれ!"
	case Good:
		return "😁 良いね!"
	case Perfect:
		return "🤩 完璧!"
	default:
		return "🤩 エラー!"
	}
}

const (
	SuperBad Evaluation = iota
	Bad
	Good
	Perfect
)

type ResponseContent struct {
	Date     string `json:"date"`
	TimeB    string `json:"time_bedin"`
	TimeW    string `json:"time_wake"`
	Duration string `json:"duration"`
	Eval     string `json:"evaluation"`
}

type ResponseContents struct {
	Record []ResponseContent `json:"record"`
	Avg    float64           `json:"average"`
}

func (r ResponseContents) GetLineMessage() string {
	var head string = fmt.Sprintf("****睡眠記録****\n平均睡眠時間: %g時間\n\n各日にちの睡眠記録\n", r.Avg)
	strList := make([]string, len(r.Record))
	for i, record := range r.Record {
		str := fmt.Sprintf("【%s】: %s\n\t就寝: %s\n\t起床: %s\n\t睡眠時間: %s\n\n",
			record.Date, record.Eval,
			record.TimeB,
			record.TimeW,
			record.Duration,
		)
		strList[i] = str
	}
	body := strings.Join(strList, "")
	return head + body
}
