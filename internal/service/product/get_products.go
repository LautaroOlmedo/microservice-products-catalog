package product

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) GetProducts(ctx context.Context, limit int) ([]domain.Product, error) {
	products, err := s.Storage.GetProducts(ctx, limit)
	if err != nil {
		return []domain.Product{}, fmt.Errorf("error fetching products: %w", err)
	}
	// TODO [technical debate] throw a goroutine for cache products
	return products, nil
}
