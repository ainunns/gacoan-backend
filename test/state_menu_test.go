package test

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
	"fp-kpl/domain/menu/category"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/shared"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for menu availability tests
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetAllCategories(ctx context.Context, tx interface{}) ([]category.Category, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]category.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, tx interface{}, id string) (category.Category, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(category.Category), args.Error(1)
}

type MockMenuRepositoryForAvailability struct {
	mock.Mock
}

func (m *MockMenuRepositoryForAvailability) GetAllMenus(ctx context.Context, tx interface{}) ([]menu.Menu, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForAvailability) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu.Menu, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForAvailability) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu.Menu, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForAvailability) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu.Menu, error) {
	args := m.Called(ctx, tx, id, isAvailable)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func TestUpdateMenuAvailability_Success_AvailableToUnavailable(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Updated menu (unavailable)
	updatedMenu := menu.Menu{
		ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
		CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name:        "Test Burger",
		ImageURL:    shared.URL{Path: "https://example.com/burger.jpg"},
		Price:       shared.Price{Price: decimal.NewFromInt(25000)},
		IsAvailable: false, // Changed to false
		CookingTime: 30 * time.Minute,
		Description: "A delicious test burger",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(), // Should be updated
		},
	}

	categoryDetail := category.Category{
		ID:   identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name: "Burgers",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	// Set up expectations
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, false).Return(updatedMenu, nil)
	mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(categoryDetail, nil)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, menuID, result.ID)
	assert.Equal(t, "Test Burger", result.Name)
	assert.Equal(t, false, result.IsAvailable) // Should be unavailable
	assert.Equal(t, "Burgers", result.Category.Name)
	assert.Equal(t, decimal.NewFromInt(25000), result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateMenuAvailability_Success_UnavailableToAvailable(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Updated menu (available)
	updatedMenu := menu.Menu{
		ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
		CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name:        "Test Pizza",
		ImageURL:    shared.URL{Path: "https://example.com/pizza.jpg"},
		Price:       shared.Price{Price: decimal.NewFromInt(35000)},
		IsAvailable: true, // Changed to true
		CookingTime: 45 * time.Minute,
		Description: "A delicious test pizza",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(), // Should be updated
		},
	}

	categoryDetail := category.Category{
		ID:   identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name: "Pizza",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	// Set up expectations
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, true).Return(updatedMenu, nil)
	mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(categoryDetail, nil)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, menuID, result.ID)
	assert.Equal(t, "Test Pizza", result.Name)
	assert.Equal(t, true, result.IsAvailable) // Should be available
	assert.Equal(t, "Pizza", result.Category.Name)
	assert.Equal(t, decimal.NewFromInt(35000), result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateMenuAvailability_Success_SameState(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Menu stays in the same state (available)
	updatedMenu := menu.Menu{
		ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
		CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name:        "Test Salad",
		ImageURL:    shared.URL{Path: "https://example.com/salad.jpg"},
		Price:       shared.Price{Price: decimal.NewFromInt(18000)},
		IsAvailable: true, // Stays available
		CookingTime: 15 * time.Minute,
		Description: "A healthy test salad",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	categoryDetail := category.Category{
		ID:   identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name: "Salads",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	// Set up expectations
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, true).Return(updatedMenu, nil)
	mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(categoryDetail, nil)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, menuID, result.ID)
	assert.Equal(t, "Test Salad", result.Name)
	assert.Equal(t, true, result.IsAvailable) // Should remain available
	assert.Equal(t, "Salads", result.Category.Name)
	assert.Equal(t, decimal.NewFromInt(18000), result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateMenuAvailability_MenuNotFound(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()

	// Set up expectations - menu not found
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, false).Return(menu.Menu{}, menu.ErrorMenuNotFound)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, false)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, menu.ErrorUpdateMenuAvailability, err)
	assert.Equal(t, response.Menu{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "GetCategoryByID")
}

func TestUpdateMenuAvailability_CategoryNotFound(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Updated menu
	updatedMenu := menu.Menu{
		ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
		CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name:        "Test Item",
		ImageURL:    shared.URL{Path: "https://example.com/item.jpg"},
		Price:       shared.Price{Price: decimal.NewFromInt(20000)},
		IsAvailable: false,
		CookingTime: 25 * time.Minute,
		Description: "A test item",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	// Set up expectations
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, false).Return(updatedMenu, nil)
	mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(category.Category{}, category.ErrorCategoryNotFound)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, false)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, category.ErrorGetCategoryByID, err)
	assert.Equal(t, response.Menu{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateMenuAvailability_RepositoryError(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()

	// Set up expectations - repository error
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, true).Return(menu.Menu{}, menu.ErrorUpdateMenuAvailability)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, true)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, menu.ErrorUpdateMenuAvailability, err)
	assert.Equal(t, response.Menu{}, result)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "GetCategoryByID")
}

func TestUpdateMenuAvailability_WithDecimalPrice(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Menu with decimal price
	updatedMenu := menu.Menu{
		ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
		CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name:        "Premium Coffee",
		ImageURL:    shared.URL{Path: "https://example.com/coffee.jpg"},
		Price:       shared.Price{Price: decimal.NewFromFloat(12500.50)},
		IsAvailable: true,
		CookingTime: 5 * time.Minute,
		Description: "A premium coffee drink",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	categoryDetail := category.Category{
		ID:   identity.NewIDFromSchema(uuid.MustParse(categoryID)),
		Name: "Beverages",
		Timestamp: shared.Timestamp{
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	// Set up expectations
	mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, true).Return(updatedMenu, nil)
	mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(categoryDetail, nil)

	// Act
	result, err := menuService.UpdateMenuAvailability(ctx, menuID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, menuID, result.ID)
	assert.Equal(t, "Premium Coffee", result.Name)
	assert.Equal(t, true, result.IsAvailable)
	assert.Equal(t, "Beverages", result.Category.Name)
	assert.Equal(t, decimal.NewFromFloat(12500.50), result.Price)

	mockMenuRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateMenuAvailability_StateTransitionValidation(t *testing.T) {
	// Arrange
	mockMenuRepo := new(MockMenuRepositoryForAvailability)
	mockCategoryRepo := new(MockCategoryRepository)

	menuService := service.NewMenuService(mockMenuRepo, mockCategoryRepo)

	ctx := context.Background()
	menuID := uuid.New().String()
	categoryID := uuid.New().String()

	// Test multiple state transitions
	testCases := []struct {
		name           string
		initialState   bool
		targetState    bool
		expectedResult bool
	}{
		{"Available to Unavailable", true, false, false},
		{"Unavailable to Available", false, true, true},
		{"Available to Available", true, true, true},
		{"Unavailable to Unavailable", false, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks for each test case
			mockMenuRepo.ExpectedCalls = nil
			mockCategoryRepo.ExpectedCalls = nil

			updatedMenu := menu.Menu{
				ID:          identity.NewIDFromSchema(uuid.MustParse(menuID)),
				CategoryID:  identity.NewIDFromSchema(uuid.MustParse(categoryID)),
				Name:        "Test Item",
				ImageURL:    shared.URL{Path: "https://example.com/item.jpg"},
				Price:       shared.Price{Price: decimal.NewFromInt(20000)},
				IsAvailable: tc.expectedResult,
				CookingTime: 20 * time.Minute,
				Description: "A test item",
				Timestamp: shared.Timestamp{
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now(),
				},
			}

			categoryDetail := category.Category{
				ID:   identity.NewIDFromSchema(uuid.MustParse(categoryID)),
				Name: "Test Category",
				Timestamp: shared.Timestamp{
					CreatedAt: time.Now().Add(-48 * time.Hour),
					UpdatedAt: time.Now().Add(-24 * time.Hour),
				},
			}

			// Set up expectations
			mockMenuRepo.On("UpdateMenuAvailability", ctx, nil, menuID, tc.targetState).Return(updatedMenu, nil)
			mockCategoryRepo.On("GetCategoryByID", ctx, nil, categoryID).Return(categoryDetail, nil)

			// Act
			result, err := menuService.UpdateMenuAvailability(ctx, menuID, tc.targetState)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedResult, result.IsAvailable, "Expected availability state %v, got %v", tc.expectedResult, result.IsAvailable)

			mockMenuRepo.AssertExpectations(t)
			mockCategoryRepo.AssertExpectations(t)
		})
	}
}
