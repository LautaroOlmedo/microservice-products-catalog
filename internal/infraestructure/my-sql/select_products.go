package my_sql

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) GetProducts(ctx context.Context, limit int) ([]domain.Product, error) {

	var products []domain.Product

	db := r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Limit(limit).
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}
	return products, nil
}
