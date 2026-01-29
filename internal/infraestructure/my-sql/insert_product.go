package my_sql

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) SaveProduct(ctx context.Context, product *domain.Product) error {
	db := r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	result := db.
		WithContext(ctx).
		Save(product)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	fmt.Printf("[LOG] - Product with ID : %s saved correctly\n", product.ID)
	return nil
}
