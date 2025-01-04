package main

type OrderPaidEvent struct {
	OrderId         string         `json:"order_id"` // UUID
	ProductQuantity map[string]int `json:"product_quantity"`
	UserId          string         `json:"user_id"`
	CreatedAt       time.Time      `json:"created_at"` // Timestamp of creation
}
