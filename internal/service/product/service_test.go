package product

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"microservice-products-catalog/internal/service/product/mocks"
	"testing"
)

// TestNewService verifies that the service constructor correctly initializes
// the service with its dependencies.
func TestNewService(t *testing.T) {
	// Arrange: Set up the gomock controller and create mock dependencies.
	ctrl := gomock.NewController(t)

	// Create mock instances using auto-generated constructors.
	mockTransaction := mocks.NewMockTransactionManager(ctrl)
	mockStorage := mocks.NewMockStorageRepository(ctrl)

	// Act: Call the constructor function that we are testing.
	service := NewService(mockStorage, mockTransaction)

	// Assert: Verify the outcome.
	// 1. Ensure the service object was actually created.
	assert.NotNil(t, service)

	// 2. Ensure the dependencies were assigned to correct fields.
	// This confirms the service holds the dependencies it needs to operate.
	assert.Equal(t, mockStorage, service.Storage, "Storage should be the provided mock instance")
}
