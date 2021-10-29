package entity

type ResponseContents map[string][]PlotData

type PlotData struct {
	TimeB int64
	TimeW int64
}
