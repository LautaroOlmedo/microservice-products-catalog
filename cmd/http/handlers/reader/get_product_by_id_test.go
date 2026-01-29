package reader_test

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"microservice-products-catalog/cmd/http/handlers/reader"
	"microservice-products-catalog/cmd/http/handlers/reader/mocks"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetProductByID(t *testing.T) {
	mockProduct := &domain.Product{
		ID:          uuid.New().String(),
		Name:        "Gopher",
		Description: "Realistic replic for the Gopher animal",
		Stock:       50,
	}

	testCases := []struct {
		name                 string
		setupMock            func(mock *mocks.MockProductService)
		request              *http.Request
		setupRequest         func(req *http.Request)
		expectedStatus       int
		expectedBodyContains string
		expectedJSONResponse *domain.Product
	}{
		{
			name: "Success - 200 OK",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					GetProductByID(gomock.Any(), mockProduct.ID).
					Return(mockProduct, nil).
					Times(1)
			},
			request: httptest.NewRequest(
				http.MethodGet,
				"/api/products/"+mockProduct.ID,
				nil,
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "fake-token")
			},
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: mockProduct,
		},
		{
			name:      "Failure - 400 Invalid UUID",
			setupMock: func(mock *mocks.MockProductService) {},
			request: httptest.NewRequest(
				http.MethodGet,
				"/api/products/not-a-uuid",
				nil,
			),
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "invalid product id format",
		},
		{
			name: "Failure - 404 Product not found",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					GetProductByID(gomock.Any(), mockProduct.ID).
					Return(nil, domain.ErrProductNotFound).
					Times(1)
			},
			request: httptest.NewRequest(
				http.MethodGet,
				"/api/products/"+mockProduct.ID,
				nil,
			),
			expectedStatus:       http.StatusNotFound,
			expectedBodyContains: "error fetching product",
		},
		{
			name: "Failure - 500 Internal Server Error",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					GetProductByID(gomock.Any(), mockProduct.ID).
					Return(nil, errors.New("database is down")).
					Times(1)
			},
			request: httptest.NewRequest(
				http.MethodGet,
				"/api/products/"+mockProduct.ID,
				nil,
			),
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error fetching product",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenGenerator := mocks.NewMockTokenGenerator(ctrl)
			mockTokenGenerator.EXPECT().
				Generate(gomock.Any(), gomock.Any()).
				Return("fake-jwt-token", nil).
				AnyTimes()

			mockOrderService := mocks.NewMockOrderService(ctrl)
			mockProductService := mocks.NewMockProductService(ctrl)

			if tc.setupMock != nil {
				tc.setupMock(mockProductService)
			}

			handler := reader.NewReaderHandler(
				mockProductService,
				mockOrderService,
				mockTokenGenerator,
			)

			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			handler.HandleGetProductByID(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBodyContains != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedBodyContains)
			}

			if tc.expectedJSONResponse != nil {
				expectedJSON, err := json.Marshal(tc.expectedJSONResponse)
				require.NoError(t, err)

				assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
				assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			}
		})
	}
}
