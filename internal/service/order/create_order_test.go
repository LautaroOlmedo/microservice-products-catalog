package order_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/internal/domain"
	"microservice-products-catalog/internal/service/order"
	"microservice-products-catalog/internal/service/order/mocks"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	product := &domain.Product{
		ID:          "076e76d6-fc3e-4f95-a024-1b4984e76060",
		Name:        "Gopher",
		Description: "Realistic replic for the Gopher animal",
		Price:       65.42,
		Stock:       50,
	}

	type testCase struct {
		testName      string
		productID     string
		quantity      int
		setupMock     func(mockStorage *mocks.MockStorageRepository, mockProductService *mocks.MockProductService, mockTxManager *mocks.MockTransactionManager)
		expectedError error
	}

	testCases := []testCase{
		{
			testName:  "Success - create order",
			productID: product.ID,
			quantity:  5,
			setupMock: func(mockStorage *mocks.MockStorageRepository, mockProductService *mocks.MockProductService, mockTxManager *mocks.MockTransactionManager) {
				mockTxManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).Times(1)

				mockProductService.EXPECT().
					GetProductByID(gomock.Any(), product.ID).
					Return(product, nil).Times(1)

				mockProductService.EXPECT().
					SaveProduct(gomock.Any(), product).
					Return(nil).Times(1)

				mockStorage.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			testName:  "Failure - Product not Found",
			productID: "non-existent",
			quantity:  5,
			setupMock: func(mockStorage *mocks.MockStorageRepository, mockProductService *mocks.MockProductService, mockTxManager *mocks.MockTransactionManager) {
				mockTxManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).Times(1)

				mockProductService.EXPECT().
					GetProductByID(gomock.Any(), "non-existent").
					Return(nil, domain.ErrProductNotFound).Times(1)
			},
			expectedError: domain.ErrProductNotFound,
		},
		{
			testName:  "Failure - Insufficient Stock",
			productID: product.ID,
			quantity:  100,
			setupMock: func(mockStorage *mocks.MockStorageRepository, mockProductService *mocks.MockProductService, mockTxManager *mocks.MockTransactionManager) {
				mockTxManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).Times(1)

				mockProductService.EXPECT().
					GetProductByID(gomock.Any(), product.ID).
					Return(product, nil).Times(1)
			},
			expectedError: domain.ErrInsufficientStock,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockStorageRepository(ctrl)
			productServiceMock := mocks.NewMockProductService(ctrl)
			txManagerMock := mocks.NewMockTransactionManager(ctrl)

			if tc.setupMock != nil {
				tc.setupMock(mockStorage, productServiceMock, txManagerMock)
			}

			service := order.NewService(mockStorage, txManagerMock, productServiceMock)

			err := service.CreateOrder(context.Background(), tc.productID, tc.quantity)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestCreateOrder_Concurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorageRepository(ctrl)
	mockProductService := mocks.NewMockProductService(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	service := order.NewService(mockStorage, mockTxManager, mockProductService)

	const (
		initialStock = 50
		orderQty     = 6
		goroutines   = 10
	)

	productID := "product-123"

	// Shared stock
	var stock atomic.Int32
	stock.Store(initialStock)

	// Tx manager always executed the function
	mockTxManager.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	mockProductService.EXPECT().
		GetProductByID(gomock.Any(), productID).
		DoAndReturn(func(ctx context.Context, id string) (*domain.Product, error) {
			return &domain.Product{
				ID:    productID,
				Stock: int(stock.Load()),
				Price: 10,
			}, nil
		}).
		AnyTimes()

	// SaveProduct decrement the stock periodically and atomic way
	mockProductService.EXPECT().
		SaveProduct(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *domain.Product) error {
			for {
				current := stock.Load()
				if current < int32(orderQty) {
					return domain.ErrInsufficientStock
				}
				if stock.CompareAndSwap(current, current-int32(orderQty)) {
					return nil
				}
			}
		}).
		AnyTimes()

	mockStorage.EXPECT().
		CreateOrder(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	var wg sync.WaitGroup
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := service.CreateOrder(context.Background(), productID, orderQty)
			errs <- err
		}()
	}

	wg.Wait()
	close(errs)

	var success, insufficient int

	for err := range errs {
		if err == nil {
			success++
			continue
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			insufficient++
			continue
		}
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 8, success)
	assert.Equal(t, 2, insufficient)
	assert.Equal(t, int32(2), stock.Load()) // 50 - (8 * 6) = 2
}
