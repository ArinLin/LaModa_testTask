package core

const (
	// parsing resps
	InvalidIDCode   = "invalid_id"
	InvalidBodyCode = "invalid_request_body"
	IDRequiredCode  = "id_is_required"

	// goods resps
	GoodReceivedCode  = "good_received"
	GoodsReceivedCode = "goods_received"
	GoodCreatedCode   = "good_created"
	GoodUpdatedCode   = "good_updated"
	GoodDeletedCode   = "good_deleted"
	GoodNotFoundCode  = "good_not_found"

	// stocks resps
	StocksReceivedCode   = "stocks_received"
	StocksReservedCode   = "stocks_reserved"
	StocksReleasedCode   = "stocks_released"
	StockNotFoundCode    = "stock_not_found"
	NotEnoughAmountCode  = "not_enough_amount"
	NotEnoughReserveCode = "not_enough_reserve"

	// warehouses resps
	WarehouseReceivedCode  = "warehouse_received"
	WarehousesReceivedCode = "warehouses_received"
	WarehouseCreatedCode   = "warehouse_created"
	WarehouseUpdatedCode   = "warehouse_updated"
	WarehouseDeletedCode   = "warehouse_deleted"
	WarehouseNotFoundCode  = "warehouse_not_found"

	// general
	InternalErrorCode = "general_internal"
)
