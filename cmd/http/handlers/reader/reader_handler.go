package reader

import (
	"context"
	"microservice-products-catalog/cmd/http/auth"
	"microservice-products-catalog/internal/domain"
)

//go:generate mockgen -source=reader_handler.go -destination=./mocks/reader_handler_mocks.go -package=mocks

type TokenGenerator interface {
	Generate(ctx context.Context, claims auth.TokenClaims) (string, error)
}

type ProductService interface {
	GetProducts(ctx context.Context, limit int) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
}

type OrderService interface {
	GetOrders(ctx context.Context) ([]domain.Order, error)
}

type ReaderHandler struct {
	ProductService ProductService
	OrderService   OrderService
	TokenGenerator TokenGenerator
}

func NewReaderHandler(productService ProductService, orderService OrderService, tokenGenerator TokenGenerator) *ReaderHandler {
	return &ReaderHandler{
		ProductService: productService,
		OrderService:   orderService,
		TokenGenerator: tokenGenerator,
	}
}
