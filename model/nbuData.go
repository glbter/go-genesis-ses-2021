package model

type NbuData struct {
	Id     int     `json:"r030"`
	Txt    string  `json:"txt"`
	Rate   float64 `json:"rate"`
	Symbol string  `json:"cc"`
	Date   string  `json:"exchangedate"`
}
