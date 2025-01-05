package products

import (
	"time"
)

// Product struct represents the products table in the stores
type Product struct {
	ID          string    `json:"id"` // UUID
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"` // Price should be in text format
	Category    string    `json:"category"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"created_at"` // Timestamp of creation
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp of last update
}

// NewProduct struct represents the data required when creating a new product
type NewProduct struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required"`
	Price       string `json:"price" validate:"required"` // Price should be in text format
	Category    string `json:"category" validate:"required"`
	Stock       uint   `json:"stock" validate:"required,min=1"`
	// required: mandatory fields.
}

type StripeProduct struct {
	ID              string    `json:"id"`
	ProductId       string    `json:"product_id"`
	StripeProductId string    `json:"stripe_product_id"`
	PriceId         string    `json:"price_id"`
	Price           uint64    `json:"price"`
	CreatedAt       time.Time `json:"created_at"` // Timestamp of creation
	UpdatedAt       time.Time `json:"updated_at"` // Timestamp of last update
}

type ProductInfo struct {
	Id      string `json:"product_id"`
	Stock   uint   `json:"stock"`
	PriceId string `json:"price_id"`
	Price   int    `json:"price"`
}
