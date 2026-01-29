package writer_test

import (
	"bytes"
	"context"
	"microservice-products-catalog/internal/domain"

	"encoding/json"
	"errors"
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
	name        = "Gopher"
	description = "Realistic replic for the Gopher animal"
	price       = 12.21
	stock       = 50
)

func TestHandleCreateProduct(t *testing.T) {

	payload := map[string]any{
		"name":        name,
		"description": description,
		"price":       price,
		"stock":       stock,
	}

	b, _ := json.Marshal(payload)

	type testCase struct {
		testName             string
		request              *http.Request
		setupRequest         func(req *http.Request)
		setupMock            func(mock *mocks.MockProductService)
		expectedStatus       int
		expectedBodyContains string
	}

	testCases := []testCase{
		{
			testName: "Success - 201 Created with valid input",
			request: func() *http.Request {
				b, _ := json.Marshal(payload)
				return httptest.NewRequest(
					http.MethodPost,
					"/api/products",
					bytes.NewReader(b),
				)
			}(),
			setupRequest: func(req *http.Request) {
				// add token header
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().
					CreateProduct(gomock.Any(), gomock.AssignableToTypeOf(domain.Product{})).
					DoAndReturn(func(ctx context.Context, p domain.Product) error {
						if p.Name != name || p.Description != description || p.Stock != stock {
							t.Errorf("unexpected product: %+v", p)
						}
						return nil
					}).Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			testName: "Failure - 400 Bad Request Reading Body",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/products",
				strings.NewReader("{invalid-json"),
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock:            func(mock *mocks.MockProductService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "error reading body",
		},

		{
			testName: "Failure - 500 Internal Server Error from error",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/products",
				strings.NewReader(string(b))),
			setupRequest: func(req *http.Request) {

			},
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().CreateProduct(gomock.Any(), gomock.Any()).Return(errors.New("database constraint violation"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error creating product",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockOrderService := mocks.NewMockOrderService(ctrl)
			mockProductService := mocks.NewMockProductService(ctrl)
			tc.setupMock(mockProductService)

			writerHandler := writer.NewWriteHandler(mockProductService, mockOrderService)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			writerHandler.HandleCreateProduct(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBodyContains != "" {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedBodyContains))
			}
		})
	}
}
