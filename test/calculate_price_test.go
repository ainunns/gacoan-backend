package test

import (
	"context"
	"fp-kpl/application/request"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/shared"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockOrderRepositoryForCalculatePrice struct {
	mock.Mock
}

func (m *MockOrderRepositoryForCalculatePrice) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	args := m.Called(ctx, tx, orderEntity)
	return args.Get(0).(order.Order), args.Error(1)
}

func (m *MockOrderRepositoryForCalculatePrice) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).([]order.Order), args.Error(1)
}

type MockMenuRepositoryForCalculatePrice struct {
	mock.Mock
}

// Mock order domain service
type MockOrderDomainService struct {
	mock.Mock
}

func (m *MockOrderDomainService) CalculatePrice(ctx context.Context, price shared.Price, quantity int64) (shared.Price, error) {
	args := m.Called(ctx, price, quantity)
	return args.Get(0).(shared.Price), args.Error(1)
}

func (m *MockMenuRepositoryForCalculatePrice) GetAllMenus(ctx context.Context, tx interface{}) ([]menu.Menu, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForCalculatePrice) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu.Menu, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForCalculatePrice) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu.Menu, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForCalculatePrice) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu.Menu, error) {
	args := m.Called(ctx, tx, id, isAvailable)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func TestCalculateTotalPrice_Success(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	// Test data
	menuID1 := uuid.New().String()
	menuID2 := uuid.New().String()

	orders := []request.Order{
		{MenuID: menuID1, Quantity: 2},
		{MenuID: menuID2, Quantity: 1},
	}

	// Mock menu responses
	menu1 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID1)),
		Name:  "Burger",
		Price: shared.Price{Price: decimal.NewFromInt(25000)},
	}

	menu2 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID2)),
		Name:  "Fries",
		Price: shared.Price{Price: decimal.NewFromInt(15000)},
	}

	// Set up expectations
	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID1).Return(menu1, nil)
	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID2).Return(menu2, nil)

	// Mock order domain service expectations
	mockOrderDomainService.On("CalculatePrice", ctx, menu1.Price, int64(2)).Return(shared.Price{Price: decimal.NewFromInt(50000)}, nil)
	mockOrderDomainService.On("CalculatePrice", ctx, menu2.Price, int64(1)).Return(shared.Price{Price: decimal.NewFromInt(15000)}, nil)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Expected: (25000 * 2) + (15000 * 1) = 65000
	expectedPrice := decimal.NewFromInt(65000)
	assert.True(t, result.Price.Equal(expectedPrice), "Expected price %s, got %s", expectedPrice, result.Price)

	mockMenuRepo.AssertExpectations(t)
}

func TestCalculateTotalPrice_SingleItem(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID := uuid.New().String()
	orders := []request.Order{
		{MenuID: menuID, Quantity: 3},
	}

	menu := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID)),
		Name:  "Pizza",
		Price: shared.Price{Price: decimal.NewFromInt(30000)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID).Return(menu, nil)

	// Mock order domain service expectation
	mockOrderDomainService.On("CalculatePrice", ctx, menu.Price, int64(3)).Return(shared.Price{Price: decimal.NewFromInt(90000)}, nil)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedPrice := decimal.NewFromInt(90000) // 30000 * 3
	assert.True(t, result.Price.Equal(expectedPrice), "Expected price %s, got %s", expectedPrice, result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertExpectations(t)
}

func TestCalculateTotalPrice_ZeroQuantity(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID := uuid.New().String()
	orders := []request.Order{
		{MenuID: menuID, Quantity: 0},
	}

	// Mock menu response - GetMenuByID is called before quantity check
	menu := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID)),
		Name:  "Test Item",
		Price: shared.Price{Price: decimal.NewFromInt(10000)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID).Return(menu, nil)

	// Mock order domain service expectation for invalid quantity
	mockOrderDomainService.On("CalculatePrice", ctx, menu.Price, int64(0)).Return(shared.Price{}, order.ErrorInvalidQuantity)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, order.ErrorInvalidQuantity, err)
	assert.Equal(t, shared.Price{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertExpectations(t)
}

func TestCalculateTotalPrice_NegativeQuantity(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID := uuid.New().String()
	orders := []request.Order{
		{MenuID: menuID, Quantity: -1},
	}

	// Mock menu response - GetMenuByID is called before quantity check
	menu := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID)),
		Name:  "Test Item",
		Price: shared.Price{Price: decimal.NewFromInt(10000)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID).Return(menu, nil)

	// Mock order domain service expectation for invalid quantity
	mockOrderDomainService.On("CalculatePrice", ctx, menu.Price, int64(-1)).Return(shared.Price{}, order.ErrorInvalidQuantity)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, order.ErrorInvalidQuantity, err)
	assert.Equal(t, shared.Price{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertExpectations(t)
}

func TestCalculateTotalPrice_MenuNotFound(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID := uuid.New().String()
	orders := []request.Order{
		{MenuID: menuID, Quantity: 1},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID).Return(menu.Menu{}, menu.ErrorMenuNotFound)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, menu.ErrorMenuNotFound, err)
	assert.Equal(t, shared.Price{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertNotCalled(t, "CalculatePrice")
}

func TestCalculateTotalPrice_EmptyOrders(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	orders := []request.Order{}

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedPrice := decimal.NewFromInt(0)
	assert.True(t, result.Price.Equal(expectedPrice), "Expected price %s, got %s", expectedPrice, result.Price)

	mockMenuRepo.AssertNotCalled(t, "GetMenuByID")
	mockOrderDomainService.AssertNotCalled(t, "CalculatePrice")
}

func TestCalculateTotalPrice_DecimalPrices(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID1 := uuid.New().String()
	menuID2 := uuid.New().String()

	orders := []request.Order{
		{MenuID: menuID1, Quantity: 2},
		{MenuID: menuID2, Quantity: 1},
	}

	// Mock menu responses with decimal prices
	menu1 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID1)),
		Name:  "Coffee",
		Price: shared.Price{Price: decimal.NewFromFloat(12500.50)},
	}

	menu2 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID2)),
		Name:  "Cake",
		Price: shared.Price{Price: decimal.NewFromFloat(8750.25)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID1).Return(menu1, nil)
	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID2).Return(menu2, nil)

	// Mock order domain service expectations
	mockOrderDomainService.On("CalculatePrice", ctx, menu1.Price, int64(2)).Return(shared.Price{Price: decimal.NewFromFloat(25001.00)}, nil)
	mockOrderDomainService.On("CalculatePrice", ctx, menu2.Price, int64(1)).Return(shared.Price{Price: decimal.NewFromFloat(8750.25)}, nil)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Expected: (12500.50 * 2) + (8750.25 * 1) = 33751.25
	expectedPrice := decimal.NewFromFloat(33751.25)
	assert.True(t, result.Price.Equal(expectedPrice), "Expected price %s, got %s", expectedPrice, result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertExpectations(t)
}

func TestCalculateTotalPrice_LargeQuantities(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID := uuid.New().String()
	orders := []request.Order{
		{MenuID: menuID, Quantity: 100},
	}

	menu := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID)),
		Name:  "Bulk Item",
		Price: shared.Price{Price: decimal.NewFromInt(1000)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID).Return(menu, nil)

	// Mock order domain service expectation
	mockOrderDomainService.On("CalculatePrice", ctx, menu.Price, int64(100)).Return(shared.Price{Price: decimal.NewFromInt(100000)}, nil)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedPrice := decimal.NewFromInt(100000) // 1000 * 100
	assert.True(t, result.Price.Equal(expectedPrice), "Expected price %s, got %s", expectedPrice, result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockOrderDomainService.AssertExpectations(t)
}

func TestCalculateTotalPrice_MixedValidAndInvalid(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForCalculatePrice)
	mockOrderRepo := new(MockOrderRepositoryForCalculatePrice)
	mockOrderDomainService := new(MockOrderDomainService)

	orderService := service.NewOrderService(mockOrderRepo, mockMenuRepo, mockOrderDomainService)

	ctx := context.Background()

	menuID1 := uuid.New().String()
	menuID2 := uuid.New().String()

	orders := []request.Order{
		{MenuID: menuID1, Quantity: 2}, // Valid
		{MenuID: menuID2, Quantity: 0}, // Invalid - should fail after GetMenuByID
	}

	menu1 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID1)),
		Name:  "Valid Item",
		Price: shared.Price{Price: decimal.NewFromInt(10000)},
	}

	menu2 := menu.Menu{
		ID:    identity.NewIDFromSchema(uuid.MustParse(menuID2)),
		Name:  "Invalid Item",
		Price: shared.Price{Price: decimal.NewFromInt(5000)},
	}

	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID1).Return(menu1, nil)
	mockMenuRepo.On("GetMenuByID", ctx, nil, menuID2).Return(menu2, nil)

	// Mock order domain service expectations
	mockOrderDomainService.On("CalculatePrice", ctx, menu1.Price, int64(2)).Return(shared.Price{Price: decimal.NewFromInt(20000)}, nil)
	mockOrderDomainService.On("CalculatePrice", ctx, menu2.Price, int64(0)).Return(shared.Price{}, order.ErrorInvalidQuantity)

	// Act
	result, err := orderService.CalculateTotalPrice(ctx, orders)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, order.ErrorInvalidQuantity, err)
	assert.Equal(t, shared.Price{}, result)

	// Both GetMenuByID calls should be made since quantity check happens after
	mockMenuRepo.AssertNumberOfCalls(t, "GetMenuByID", 2)
	mockOrderDomainService.AssertNumberOfCalls(t, "CalculatePrice", 2)
}
