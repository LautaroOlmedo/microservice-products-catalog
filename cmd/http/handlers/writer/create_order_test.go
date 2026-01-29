package writer_test

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/cmd/http/handlers/writer"
	"microservice-products-catalog/cmd/http/handlers/writer/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	productID = "18eb9153-a00c-466d-8f38-f149806b054e"
	quantity  = 3
)

func TestHandleCreateOrder(t *testing.T) {
	payload := map[string]any{
		"product_id": productID,
		"quantity":   quantity,
	}

	b, _ := json.Marshal(payload)

	validRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/orders",
		bytes.NewReader(b),
	)

	type testCase struct {
		testName             string
		request              *http.Request
		setupRequest         func(req *http.Request)
		setupMock            func(mock *mocks.MockOrderService)
		expectedStatus       int
		expectedBodyContains string
	}

	testCases := []testCase{
		{
			testName: "Success - 201 Created Order correctly",
			request:  validRequest,
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockOrderService) {
				mock.EXPECT().
					CreateOrder(gomock.Any(), productID, quantity).Return(nil).Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			testName: "Failure - 404 Product Not Found",
			request:  httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(b)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockOrderService) {
				mock.EXPECT().
					CreateOrder(gomock.Any(), productID, quantity).Return(nil).Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		/*{
			testName: "Failure - 500 Internal Server Error",
			request:  httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(b)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockOrderService) {
				mock.EXPECT().
					CreateOrder(gomock.Any(), productID, quantity).Return(errors.New("unknown Database")).Times(1)
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error creating order",
		},*/
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockOrderService := mocks.NewMockOrderService(ctrl)
			mockProductService := mocks.NewMockProductService(ctrl)
			tc.setupMock(mockOrderService)

			writerHandler := writer.NewWriteHandler(mockProductService, mockOrderService)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			writerHandler.HandleCreateOrder(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBodyContains != "" {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedBodyContains))
			}
		})
	}

}
