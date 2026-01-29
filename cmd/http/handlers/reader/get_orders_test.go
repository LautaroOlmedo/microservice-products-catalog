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
	"time"
)

func TestHandleGetOrders(t *testing.T) {

	mocksOrders := []domain.Order{
		{ID: uuid.New().String(), ProductID: uuid.New().String(), Quantity: 5, Total: 32.23, Date: time.Now()},
		{ID: uuid.New().String(), ProductID: uuid.New().String(), Quantity: 7, Total: 21.90, Date: time.Now()},
	}

	testCases := []struct {
		name                 string
		setupMock            func(mock *mocks.MockOrderService)
		request              *http.Request
		setupRequest         func(req *http.Request)
		expectedStatus       int
		expectedBodyContains string
		expectedJSONResponse []domain.Order
	}{
		{
			name: "Success - 200 Get Orders",
			setupMock: func(mock *mocks.MockOrderService) {
				mock.EXPECT().
					GetOrders(gomock.Any()).
					Return(mocksOrders, nil).
					Times(1)

			},
			request: httptest.NewRequest(http.MethodGet, "/api/orders", nil),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk2MDg4NzgsImlhdCI6MTc2OTYwNzk3OCwicmVxdWVzdF9pZCI6IjJlODYxOTg3LTQwNzQtNGI1My04YmZjLTc4OTQ3NjUzZDRjNCIsInNjb3BlIjoicHJvZHVjdHM6cmVhZCJ9.3sjhs9uNj_337jGKiVfVNYBsXsb3eKKqB7o9USzck-E")
			},
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: mocksOrders,
		},

		{
			name: "Failure - 500 Internal Server Error",
			setupMock: func(mock *mocks.MockOrderService) {
				mock.EXPECT().
					GetOrders(gomock.Any()).
					Return(nil, errors.New("database is down"))
			},
			request: httptest.NewRequest(http.MethodGet, "/api/orders", nil),
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
			tc.setupMock(mockOrderService)

			readerHandler := reader.NewReaderHandler(mockProductService, mockOrderService, mockTokenGenerator)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			readerHandler.HandleGetOrders(recorder, tc.request)

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
