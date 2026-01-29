package my_sql

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) DeleteProduct(ctx context.Context, id string) error {
	var db = r.db

	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	result := db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Product{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	fmt.Printf("[LOG] - Product with ID : %s deeted correctly\n", id)
	return nil

}
