package apperrors

import "errors"

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInsufficientQuantity = errors.New("insufficient quantity")
	ErrSymbolNotAllowed     = errors.New("symbol not allowed")
)
