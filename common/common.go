package common

type OrderRequest struct {
	UserID   uint64     `json:"user_id"`
	OrderID  uint64     `json:"order_id"`
	ItemInfo []ItemInfo `json:"item_info"`
}

type ItemInfo struct {
	ItemID   uint64 `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type ItemPriceInfo struct {
	ItemID uint64  `json:"item_id"`
	Price  float64 `json:"price"`
}

type OrderCreationResponse struct {
	OrderId uint64 `json:"order_id"`
}

type OrderQueueRequest struct {
	OrderRequest OrderRequest
	OrderId      uint64 `json:"order_id"`
}
