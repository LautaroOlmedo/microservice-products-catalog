package writer

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/cmd/http/handlers/writer/mocks"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleDeleteProduct(t *testing.T) {
	productID := "b82a87a1-a15e-4767-8b26-99cbbcb8ae97"
	validRequest := httptest.NewRequest(
		http.MethodDelete, // es DELETE, no GET
		fmt.Sprintf("/products/%s", productID),
		nil,
	)

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
			request:  validRequest,
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().DeleteProduct(gomock.Any(), productID).Return(nil).Times(1)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName: "Failure - 400 Bad Request Invalid ID",
			request: httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/products/%s", "-1"),
				nil,
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock:            func(mock *mocks.MockProductService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "invalid product id format, must be UUID",
		},
		{
			testName: "Failure - 400 Bad Request Empty ID",
			request: httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/products/%s", ""),
				nil,
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock:            func(mock *mocks.MockProductService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "invalid product id",
		},
		{
			testName: "Failure - 404 Bad Request Product Not Found",
			request: httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/products/%s", "a17a87a1-a15e-4767-8b26-99cbbcb8ae65"),
				nil,
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().DeleteProduct(gomock.Any(), "a17a87a1-a15e-4767-8b26-99cbbcb8ae65").Return(domain.ErrProductNotFound).Times(1)
			},
			expectedStatus:       http.StatusNotFound,
			expectedBodyContains: "product not found",
		},
		{
			testName: "Failure - 500 Bad Internal Server Error",
			request: httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/products/%s", "a17a87a1-a15e-4767-8b26-99cbbcb8ae65"),
				nil,
			),
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockProductService) {
				mock.EXPECT().DeleteProduct(gomock.Any(), "a17a87a1-a15e-4767-8b26-99cbbcb8ae65").Return(errors.New("cannot connect with database")).Times(1)
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error deleting product",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockOrderService := mocks.NewMockOrderService(mockCtrl)
			mockProductService := mocks.NewMockProductService(mockCtrl)
			tc.setupMock(mockProductService)

			writerHandler := NewWriteHandler(mockProductService, mockOrderService)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			writerHandler.HandleDeleteProduct(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBodyContains != "" {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedBodyContains))
			}
		})
	}
}
