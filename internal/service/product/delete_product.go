package product

import (
	"context"
)

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	return s.TransactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return s.Storage.DeleteProduct(ctx, id)
	})
}
