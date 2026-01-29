package my_sql

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	db := r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	if err := db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", product.ID).
		Updates(product).
		Error; err != nil {
		return err
	}

	fmt.Printf("[LOG] - Product with ID : %s updated correctly\n", product.ID)
	return nil
}
