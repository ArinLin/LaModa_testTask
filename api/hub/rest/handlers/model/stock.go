package model

type Stock struct {
	GoodID   int    `json:"good_id" example:"1"`
	GoodName string `json:"good_name" example:"Hat"`
	Amount   int    `json:"amount" example:"10"`
	Reserved int    `json:"reserved" example:"10"`
}

type GetStocksResponse struct {
	Status  string  `json:"status" example:"OK"`
	MsgCode string  `json:"msg_code" example:"stocks_received"`
	Data    []Stock `json:"data"`
}

type ReserveStockResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"stock_reserved"`
}

type ReleaseStockResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"stock_released"`
}
