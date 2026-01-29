package order

import (
	"context"
	"fmt"
	"microservice-products-catalog/internal/domain"
)

func (s *Service) GetOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := s.Storage.GetOrders(ctx)
	if err != nil {
		return []domain.Order{}, fmt.Errorf("get orders error: %w", err)
	}
	return orders, nil

}
