package writer

import (
	"context"
	"microservice-products-catalog/internal/domain"
)

//go:generate mockgen -source=write_handler.go -destination=././mocks/product_service_mock.go -package=mocks

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, productID string, quantity int) error
}

// WriteHandler depends on the interface, not concrete types
type WriteHandler struct {
	ProductService ProductService
	OrderService   OrderService
}

func NewWriteHandler(productService ProductService, orderService OrderService) *WriteHandler {
	return &WriteHandler{
		ProductService: productService,
		OrderService:   orderService,
	}
}
