package product

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) CreateProduct(ctx context.Context, product domain.Product) error {
	// TODO [technical debate] validate if already exists
	// TODO another approach if capture duplicate_key error and return elegant message to user
	return s.TransactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return s.Storage.SaveProduct(txCtx, &product)
	})
}
