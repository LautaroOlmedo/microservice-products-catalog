package domain

import (
	"errors"
	"time"
)

var ErrProductNotFound = errors.New("product not found")
var ErrInsufficientStock = errors.New("insufficient stock")

type Product struct {
	ID          string  `sql:"id" json:"id"`
	Name        string  `sql:"name" json:"name"`
	Description string  `sql:"description" json:"description"`
	Price       float64 `sql:"price" json:"price"`
	Stock       int     `sql:"stock" json:"stock"`
}

type Order struct {
	ID        string    `sql:"id" json:"id"`
	ProductID string    `sql:"product_id" json:"product_id"`
	Quantity  int       `sql:"quantity" json:"quantity"`
	Total     float64   `sql:"total" json:"total"`
	Date      time.Time `sql:"created_at" json:"created_at"`
}
