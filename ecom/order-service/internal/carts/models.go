package carts

import "time"

type NewCartItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=1"`
}

type Cart struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CartItem struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	CartID    string `db:"cart_id" json:"cart_id"`
	Price     int    `db:"price" json:"price"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type CartItemResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
}

type CartResponse struct {
	CartItems []CartItemResponse `json:"cart_items"`
	CartID    string             `json:"cart_id"`
}

type ProductServiceResponse struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
	PriceID   string `json:"price_id"`
	Price     int    `json:"price"`
}
