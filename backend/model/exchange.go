package model

type Quote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type FinnhubQuote struct {
	CurrentPrice float64 `json:"c"` // c = Current price
	Change       float64 `json:"d"` // d = Change
	Timestamp    int64   `json:"t"` // t = Timestamp
}

type TradeRequest struct {
	UserID int64  `json:"user_id"`
	Symbol string `json:"symbol"`
	Count  int    `json:"count"`
}
