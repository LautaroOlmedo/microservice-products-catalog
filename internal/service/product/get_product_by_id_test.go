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

func TestGetProductByID(t *testing.T) {
	mockProduct := &domain.Product{
		ID: uuid.New().String(), Name: "Gopher", Description: "Realistic replic for the Gopher animal", Price: 32.23, Stock: 50,
	}
	notFoundID := uuid.New().String()

	type testCase struct {
		testName         string
		input            string
		setupMock        func(storage *mocks.MockStorageRepository)
		expectedProducts *domain.Product
		expectedError    error
	}

	testCases := []testCase{
		{
			testName: "Success - Fetch Product",
			input:    mockProduct.ID,
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetProductByID(gomock.Any(), mockProduct.ID).Return(mockProduct, nil).Times(1)
			},
			expectedProducts: mockProduct,
			expectedError:    nil,
		},
		{
			testName: "Failure - Product Not Found",
			input:    notFoundID,
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetProductByID(gomock.Any(), notFoundID).Return(mockProduct, domain.ErrProductNotFound).Times(1)
			},
			expectedProducts: nil,
			expectedError:    domain.ErrProductNotFound,
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
				tc.setupMock(mockStorage)
			}

			service := product.NewService(mockStorage, mockTransaction)

			// Act
			_, err := service.GetProductByID(context.Background(), tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
