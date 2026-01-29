package product_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/internal/domain"
	"microservice-products-catalog/internal/service/product"
	"microservice-products-catalog/internal/service/product/mocks"
	"testing"
)

func TestCreate(t *testing.T) {
	productInput := &domain.Product{
		ID:          uuid.New().String(),
		Name:        "Gopher",
		Description: "Realistic replic for the Gopher animal",
		Price:       65.42,
		Stock:       50,
	}

	// Define a reusable database error
	dbError := errors.New("database constraint violation: product already exists")

	type testCase struct {
		testName      string
		input         domain.Product
		setupMock     func(storage *mocks.MockStorageRepository, txManager *mocks.MockTransactionManager)
		expectedError error
	}

	testCases := []testCase{
		{
			testName: "Success - Create Product",
			input:    *productInput,
			setupMock: func(
				storage *mocks.MockStorageRepository,
				txManager *mocks.MockTransactionManager,
			) {
				txManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Times(1)

				storage.EXPECT().
					SaveProduct(gomock.Any(), productInput).
					Return(nil).Times(1)
			},
			expectedError: nil,
		},

		{
			testName: "Failure - Error from storage layer",
			input:    *productInput,
			setupMock: func(storage *mocks.MockStorageRepository, txManager *mocks.MockTransactionManager) {
				txManager.EXPECT().
					WithTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Times(1)
				storage.EXPECT().SaveProduct(gomock.Any(), productInput).Return(dbError).Times(1)
			},
			expectedError: dbError,
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
			err := service.CreateProduct(context.Background(), tc.input)

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
