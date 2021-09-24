package entity

type Evaluation int

func (e Evaluation) ConvertToResponse() string {
	switch e {
	case SuperBad:
		return "😱 0:伸び代しかない!"
	case Bad:
		return "😥 1:がんばれ!"
	case Good:
		return "😁 2:良いね!"
	case Perfect:
		return "🤩 3:完璧!"
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
