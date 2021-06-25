package model

type BlockchainResponse struct {
	Symbol           string  `json:"symbol"`
	Price_24h        float64 `json:"price_24h"`
	Volume_24h       float64 `json:"volume_24h"`
	Last_trade_price float64 `json:"last_trade_price"`
}
