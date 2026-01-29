package product

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) SaveProduct(ctx context.Context, product *domain.Product) error {
	return s.Storage.SaveProduct(ctx, product)
}
