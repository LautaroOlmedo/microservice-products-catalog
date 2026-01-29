package product

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.Storage.GetProductByID(ctx, id)
}
