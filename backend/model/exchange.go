package model

import "errors"

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

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInsufficientQuantity = errors.New("insufficient quantity")
var ErrSymbolNotAllowed = errors.New("symbol not allowed")
