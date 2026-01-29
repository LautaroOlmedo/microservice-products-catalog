package my_sql

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) GetOrders(ctx context.Context) ([]domain.Order, error) {

	var orders []domain.Order

	db := r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Find(&orders).
		Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
