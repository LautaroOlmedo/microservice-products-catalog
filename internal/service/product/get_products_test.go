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

func TestGetProducts(t *testing.T) {

	mocksProducts := []domain.Product{
		{ID: uuid.New().String(), Name: "Gopher", Description: "Realistic replic for the Gopher animal", Price: 32.23, Stock: 50},
		{ID: uuid.New().String(), Name: "Rusty", Description: "Realistic replic for the Rusty animal", Price: 21.90, Stock: 10},
	}
	limit := 10

	dbError := errors.New("my sql connection failed")

	type testCase struct {
		testName         string
		setupMock        func(storage *mocks.MockStorageRepository)
		expectedProducts []domain.Product
		expectedError    error
	}

	testCases := []testCase{
		{
			testName: "Success - fetch products",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetProducts(gomock.Any(), limit).Return(mocksProducts, nil).Times(1)
			},
			expectedProducts: mocksProducts,
			expectedError:    nil,
		},
		{
			testName: "Success - Products table is empty",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetProducts(gomock.Any(), limit).Return([]domain.Product{}, nil).Times(1)
			},
			expectedProducts: []domain.Product{},
			expectedError:    nil,
		},
		{
			testName: "Failure - Database fails when call to GetProducts()",
			setupMock: func(storage *mocks.MockStorageRepository) {
				storage.EXPECT().GetProducts(gomock.Any(), limit).Return([]domain.Product{}, dbError).Times(1)
			},
			expectedError:    dbError,
			expectedProducts: []domain.Product{},
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
			_, err := service.GetProducts(context.Background(), limit)

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
