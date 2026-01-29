package my_sql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"microservice-products-catalog/internal/domain"
)

func (r *Repository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {

	db := r.db
	if tx, ok := GetTx(ctx); ok {
		db = tx
	}

	var product domain.Product

	err := db.
		WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}). // This block the row while transaction is executing
		Where("id = ?", id).
		First(&product).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProductNotFound
	}

	return &product, err
}
