package product_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/internal/domain"
	"microservice-products-catalog/internal/service/product"
	"microservice-products-catalog/internal/service/product/mocks"
	"testing"
)

func TestDeleteProduct(t *testing.T) {
	productID := uuid.New().String()

	type testCase struct {
		testName      string
		input         string
		setupMock     func(storage *mocks.MockStorageRepository, txManager *mocks.MockTransactionManager)
		expectedError error
	}

	testCases := []testCase{
		{
			testName: "Success - delete product successfully",
			input:    productID,
			setupMock: func(storage *mocks.MockStorageRepository, txManager *mocks.MockTransactionManager) {
				txManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Times(1)
				storage.EXPECT().DeleteProduct(gomock.Any(), productID).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			testName: "Failure - Product Not Found",
			input:    "5b18f387-2e53-48af-9bc7-345ca1906e56",
			setupMock: func(storage *mocks.MockStorageRepository, txManager *mocks.MockTransactionManager) {
				txManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Times(1)
				storage.EXPECT().DeleteProduct(gomock.Any(), "5b18f387-2e53-48af-9bc7-345ca1906e56").Return(domain.ErrProductNotFound).Times(1)
			},
			expectedError: domain.ErrProductNotFound,
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
			if tc.setupMock != nil {
				tc.setupMock(mockStorage, mockTransaction)
			}

			service := product.NewService(mockStorage, mockTransaction)

			// Act
			err := service.DeleteProduct(context.Background(), tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
