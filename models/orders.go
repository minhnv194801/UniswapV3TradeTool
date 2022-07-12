package models

type Order struct {
	Id               int64   `json:"id"`
	SubIds           []int64 `json:"sub_ids"`
	Pubkey           string  `json:"pubkey"`
	TradeType        int64   `json:"trade_type"`
	Coin             string  `json:"coin"`
	Currency         string  `json:"currency"`
	Type             int64   `json:"type"`
	Price            float64 `json:"price"`
	Quantity         float64 `json:"quantity"`
	ExecutedQuantity float64 `json:"executed_quantity"`
	StopPrice        float64 `json:"stop_price"`
	Status           int64   `json:"status"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
	StopLimitPrice   float64 `json:"stop_limit_price"`
}

type CreateOrderReq struct {
	Coin      string  `json:"coin"`
	Currency  string  `json:"currency"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	TradeType int64   `json:"trade_type"`
	Type      int64   `json:"type"`
}
