package order

import (
	"context"
	"github.com/google/uuid"
	"microservice-products-catalog/internal/domain"
	"time"
)

func (s *Service) CreateOrder(
	ctx context.Context,
	productID string,
	quantity int,
) error {

	return s.TransactionManager.WithTransaction(ctx, func(txCtx context.Context) error {

		product, err := s.ProductService.GetProductByID(txCtx, productID)
		if err != nil {
			return err
		}

		if product.Stock < quantity {
			return domain.ErrInsufficientStock
		}

		product.Stock -= quantity

		if err := s.ProductService.SaveProduct(txCtx, product); err != nil {
			return err
		}
		order := &domain.Order{
			ID:        uuid.New().String(),
			ProductID: productID,
			Quantity:  quantity,
			Total:     product.Price * float64(quantity),
			Date:      time.Now(),
		}

		return s.Storage.CreateOrder(txCtx, *order)
	})
}
