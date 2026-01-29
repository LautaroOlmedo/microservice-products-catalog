package my_sql

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) CreateOrder(ctx context.Context, order domain.Order) error {
	db := r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	if err := db.WithContext(ctx).Create(&order).Error; err != nil {
		return err
	}

	fmt.Printf("[LOG] - Order with ID : %s saved correctly\n", order.ID)
	return nil
}
