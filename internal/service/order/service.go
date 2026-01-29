package order

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

//go:generate mockgen -source=service.go -destination=././mocks/order_repository_mock.go -package=mocks

type ProductService interface {
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	SaveProduct(ctx context.Context, product *domain.Product) error
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type StorageRepository interface {
	CreateOrder(ctx context.Context, order domain.Order) error
	GetOrders(ctx context.Context) ([]domain.Order, error)
}

type Service struct {
	Storage            StorageRepository
	TransactionManager TransactionManager
	ProductService     ProductService
}

func NewService(storageRepository StorageRepository, transactionManager TransactionManager, productService ProductService) *Service {
	return &Service{
		Storage:            storageRepository,
		TransactionManager: transactionManager,
		ProductService:     productService,
	}
}
