package writer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/cmd/http/dto"
	"microservice-products-catalog/cmd/http/handlers/writer"
	"microservice-products-catalog/cmd/http/handlers/writer/mocks"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUpdateProduct(t *testing.T) {

	newName := "New Gopher"
	newPrice := 99.9
	newStock := 20

	updateBody := dto.UpdateProductRequest{
		Name:  &newName,
		Price: &newPrice,
		Stock: &newStock,
	}

	bodyBytes, _ := json.Marshal(updateBody)

	testCases := []struct {
		name                 string
		setupMock            func(mock *mocks.MockProductService)
		request              *http.Request
		setupRequest         func(req *http.Request)
		expectedStatus       int
		expectedBodyContains string
	}{
		{
			name: "Success - 204 No Content",
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					UpdateProduct(gomock.Any(), gomock.AssignableToTypeOf(&domain.Product{})).
					DoAndReturn(func(_ context.Context, p *domain.Product) error {
						assert.Equal(t, productID, p.ID)
						assert.Equal(t, newName, p.Name)
						assert.Equal(t, newPrice, p.Price)
						assert.Equal(t, newStock, p.Stock)
						return nil
					}).
					Times(1)
			},
			request: httptest.NewRequest(
				http.MethodPut,
				"/api/products/"+productID,
				bytes.NewReader(bodyBytes),
			),
			setupRequest: func(req *http.Request) {
				//req.Header.Set("Authorization", "fake-token")
				req.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Failure - 400 Invalid UUID",
			setupMock: func(mock *mocks.MockProductService) {},
			request: httptest.NewRequest(
				http.MethodPut,
				"/api/products/not-a-uuid",
				bytes.NewReader(bodyBytes),
			),
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "invalid product id format",
		},
		{
			name:      "Failure - 400 DTO validation error",
			setupMock: func(mock *mocks.MockProductService) {},
			request: httptest.NewRequest(
				http.MethodPut,
				"/api/products/"+productID,
				bytes.NewReader([]byte(`{"price": -1}`)),
			),
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "DTO validation error",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProductService := mocks.NewMockProductService(ctrl)
			mockOrderService := mocks.NewMockOrderService(ctrl)

			if tc.setupMock != nil {
				tc.setupMock(mockProductService)
			}

			handler := writer.NewWriteHandler(
				mockProductService,
				mockOrderService,
			)

			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			handler.HandleUpdateProduct(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBodyContains != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedBodyContains)
			}
		})
	}
}
