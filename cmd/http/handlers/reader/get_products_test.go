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

func TestHandleGetProducts(t *testing.T) {

	mockProducts := []domain.Product{
		{ID: uuid.New().String(), Name: "Gopher", Description: "Realistic replic for the Gopher animal", Stock: 50},
		{ID: uuid.New().String(), Name: "Rusty", Description: "Realistic replic for the Rusty animal", Stock: 10},
	}

	testCases := []struct {
		name                 string
		setupMock            func(mock *mocks.MockProductService)
		request              *http.Request
		setupRequest         func(req *http.Request)
		expectedStatus       int
		expectedBodyContains string
		expectedJSONResponse []domain.Product
	}{
		{
			name: "Success - 200 OK with default limit",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					GetProducts(gomock.Any(), 10).
					Return(mockProducts, nil).
					Times(1)

			},
			request: httptest.NewRequest(http.MethodGet, "/api/products", nil),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk2MDg4NzgsImlhdCI6MTc2OTYwNzk3OCwicmVxdWVzdF9pZCI6IjJlODYxOTg3LTQwNzQtNGI1My04YmZjLTc4OTQ3NjUzZDRjNCIsInNjb3BlIjoicHJvZHVjdHM6cmVhZCJ9.3sjhs9uNj_337jGKiVfVNYBsXsb3eKKqB7o9USzck-E")
			},
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: mockProducts,
		},
		{
			name:      "Failure - 400 Bad Request for missing limit",
			setupMock: func(mock *mocks.MockProductService) {}, // No calls to the mock are expected
			request:   httptest.NewRequest(http.MethodGet, "/api/products?limit=-1", nil),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk2MDg4NzgsImlhdCI6MTc2OTYwNzk3OCwicmVxdWVzdF9pZCI6IjJlODYxOTg3LTQwNzQtNGI1My04YmZjLTc4OTQ3NjUzZDRjNCIsInNjb3BlIjoicHJvZHVjdHM6cmVhZCJ9.3sjhs9uNj_337jGKiVfVNYBsXsb3eKKqB7o9USzck-E")
			},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "error parsing limit: limit must be a positive number",
		},
		{
			name: "Failure - 500 Internal Server Error from service",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					GetProducts(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database is down"))
			},
			request: httptest.NewRequest(http.MethodGet, "/api/products", nil),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk2MDg4NzgsImlhdCI6MTc2OTYwNzk3OCwicmVxdWVzdF9pZCI6IjJlODYxOTg3LTQwNzQtNGI1My04YmZjLTc4OTQ3NjUzZDRjNCIsInNjb3BlIjoicHJvZHVjdHM6cmVhZCJ9.3sjhs9uNj_337jGKiVfVNYBsXsb3eKKqB7o9USzck-E")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error getting products: database is down",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenGenerator := mocks.NewMockTokenGenerator(ctrl)
			mockTokenGenerator.
				EXPECT().
				Generate(gomock.Any(), gomock.Any()).
				Return("fake-jwt-token", nil).
				AnyTimes()
			mockOrderService := mocks.NewMockOrderService(ctrl)
			mockProductService := mocks.NewMockProductService(ctrl)
			tc.setupMock(mockProductService)

			readerHandler := reader.NewReaderHandler(mockProductService, mockOrderService, mockTokenGenerator)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			readerHandler.HandleGetProducts(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBodyContains != "" {
				assert.Equal(t, recorder.Body.String(), tc.expectedBodyContains)
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
