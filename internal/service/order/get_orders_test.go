package order_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/internal/domain"
	"microservice-products-catalog/internal/service/order"
	"microservice-products-catalog/internal/service/order/mocks"
	"testing"
	"time"
)

func TestGetProducts(t *testing.T) {

	mocksOrders := []domain.Order{
		{ID: uuid.New().String(), ProductID: uuid.New().String(), Quantity: 5, Total: 32.23, Date: time.Now()},
		{ID: uuid.New().String(), ProductID: uuid.New().String(), Quantity: 7, Total: 21.90, Date: time.Now()},
	}

	dbError := errors.New("my sql connection failed")

	type testCase struct {
		testName         string
		setupMock        func(storage *mocks.MockStorageRepository)
		expectedProducts []domain.Order
		expectedError    error
	}

	testCases := []testCase{
		{
			testName: "Success - Fetch Orders",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetOrders(gomock.Any()).Return(mocksOrders, nil).Times(1)
			},
			expectedProducts: mocksOrders,
			expectedError:    nil,
		},
		{
			testName: "Success - Products table is empty",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().
					GetOrders(gomock.Any()).
					Return([]domain.Order{}, nil).
					Times(1)
			},
			expectedProducts: []domain.Order{},
			expectedError:    nil,
		},
		{
			testName: "Failure - Database fails when call to GetOrders()",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().
					GetOrders(gomock.Any()).
					Return([]domain.Order{}, dbError).
					Times(1)
			},
			expectedProducts: []domain.Order{},
			expectedError:    dbError,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransaction := mocks.NewMockTransactionManager(ctrl)
			mockStorage := mocks.NewMockStorageRepository(ctrl)
			productService := mocks.NewMockProductService(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockStorage)
			}

			service := order.NewService(mockStorage, mockTransaction, productService)

			// Act
			_, err := service.GetOrders(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))

			} else {
				assert.NoError(t, err)
			}

		})

	}
}
