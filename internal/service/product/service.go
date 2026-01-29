package product

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

//go:generate mockgen -source=service.go -destination=././mocks/product_repository_mock.go -package=mocks
type StorageRepository interface {
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	GetProducts(ctx context.Context, limit int) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id string) error
	SaveProduct(ctx context.Context, product *domain.Product) error
}

//go:generate mockgen -source=service.go -destination=././mocks/product_repository_mock.go -package=mocks

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// Service depends on the interface, not concrete types.
type Service struct {
	Storage            StorageRepository
	TransactionManager TransactionManager
}

func NewService(storage StorageRepository, transactionManager TransactionManager) *Service {
	return &Service{
		Storage:            storage,
		TransactionManager: transactionManager,
	}
}
