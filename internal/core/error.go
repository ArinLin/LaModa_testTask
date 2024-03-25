package core

import "errors"

var (
	ErrGoodNotFound      = errors.New("good not found")
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrStockNotFound     = errors.New("stock not found")
	ErrNotEnoughAmount   = errors.New("not enough amount in stock")
	ErrNotEnoughReserve  = errors.New("not enough reserve in stock")
	ErrInvalidLogLevel   = errors.New("invalid log level")
)
