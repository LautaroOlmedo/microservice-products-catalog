package product

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.TransactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.Storage.GetProductByID(ctx, product.ID)
		if err != nil {
			return err
		}
		exists = product
		return s.Storage.UpdateProduct(ctx, exists)
	})
}
