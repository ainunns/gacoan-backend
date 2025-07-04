package test

import (
	"context"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/shared"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
	"fp-kpl/platform/pagination"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for transaction tests
type MockTransactionRepositoryForGetByID struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForGetByID) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, userID, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	args := m.Called(ctx, tx, userID, id)
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(response.NextOrder), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForGetByID) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (transaction.Query, error) {
	return transaction.Query{}, nil
}

func (m *MockTransactionRepositoryForGetByID) GetDetailedTransactionByID(ctx context.Context, tx interface{}, id string) (transaction.Query, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(transaction.Query), args.Error(1)
}

// Mock other repositories
type MockUserRepositoryForTransaction struct {
	mock.Mock
}

func (m *MockUserRepositoryForTransaction) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	args := m.Called(ctx, tx, userEntity)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForTransaction) CreateUser(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	args := m.Called(ctx, tx, userEntity)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForTransaction) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForTransaction) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	args := m.Called(ctx, tx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForTransaction) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	args := m.Called(ctx, tx, email)
	return args.Get(0).(user.User), args.Get(1).(bool), args.Error(2)
}

type MockTableRepositoryForTransaction struct {
	mock.Mock
}

func (m *MockTableRepositoryForTransaction) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]table.Table), args.Error(1)
}

func (m *MockTableRepositoryForTransaction) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(table.Table), args.Error(1)
}

type MockOrderRepositoryForTransaction struct {
	mock.Mock
}

func (m *MockOrderRepositoryForTransaction) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	args := m.Called(ctx, tx, orderEntity)
	return args.Get(0).(order.Order), args.Error(1)
}

func (m *MockOrderRepositoryForTransaction) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).([]order.Order), args.Error(1)
}

type MockMenuRepositoryForTransaction struct {
	mock.Mock
}

func (m *MockMenuRepositoryForTransaction) GetAllMenus(ctx context.Context, tx interface{}) ([]menu.Menu, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForTransaction) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu.Menu, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForTransaction) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu.Menu, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).([]menu.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForTransaction) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu.Menu, error) {
	args := m.Called(ctx, tx, id, isAvailable)
	return args.Get(0).(menu.Menu), args.Error(1)
}

type MockPaymentGatewayPortForGetByID struct {
	mock.Mock
}

func (m *MockPaymentGatewayPortForGetByID) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(port.ProcessPaymentResponse), args.Error(1)
}

func (m *MockPaymentGatewayPortForGetByID) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	args := m.Called(ctx, tx, transactionID, datas)
	return args.Error(0)
}

type MockOrderServiceForGetByID struct{ mock.Mock }

func (m *MockOrderServiceForGetByID) CalculateTotalPrice(ctx context.Context, orders []request.Order) (shared.Price, error) {
	args := m.Called(ctx, orders)
	return args.Get(0).(shared.Price), args.Error(1)
}

type MockTransactionDomainServiceForGetByID struct{ mock.Mock }

func (m *MockTransactionDomainServiceForGetByID) GenerateQueueCode(ctx context.Context, transactionID string) (string, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTransactionDomainServiceForGetByID) CalculateMaxCookingTime(orders []transaction.OrderQuery) time.Duration {
	args := m.Called(orders)
	return args.Get(0).(time.Duration)
}

func (m *MockTransactionDomainServiceForGetByID) GetOrderDelayStatus(maxCookingTime time.Duration, cookedAt *time.Time, servedAt *time.Time) bool {
	args := m.Called(maxCookingTime, cookedAt, servedAt)
	return args.Get(0).(bool)
}

func TestGetTransactionByID_Success(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	userID := uuid.New().String()
	transactionID := uuid.New().String()
	tableID := uuid.New().String()
	menuID := uuid.New().String()

	// Create transaction.Query instead of schema.Transaction
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(uuid.MustParse(transactionID)),
			UserID:      identity.NewID(uuid.MustParse(userID)),
			TableID:     identity.NewID(uuid.MustParse(tableID)),
			QueueCode:   transaction.QueueCode{Code: "Q0001"},
			OrderStatus: transaction.OrderStatus{Status: "pending"},
			TotalPrice:  shared.Price{Price: decimal.NewFromInt(65000)},
			CookedAt:    timePtr(time.Now().Add(-30 * time.Minute)),
			ServedAt:    nil,
		},
		Orders: []transaction.OrderQuery{
			{
				Order: order.Order{
					ID:       identity.NewID(uuid.New()),
					Quantity: 2,
				},
				Menu: menu.Menu{
					ID:    identity.NewID(uuid.MustParse(menuID)),
					Name:  "Burger",
					Price: shared.Price{Price: decimal.NewFromInt(25000)},
				},
			},
			{
				Order: order.Order{
					ID:       identity.NewID(uuid.New()),
					Quantity: 1,
				},
				Menu: menu.Menu{
					ID:    identity.NewID(uuid.New()),
					Name:  "Fries",
					Price: shared.Price{Price: decimal.NewFromInt(15000)},
				},
			},
		},
		Table: table.Table{
			ID:          identity.NewID(uuid.MustParse(tableID)),
			TableNumber: "A1",
		},
	}

	// Set up expectations
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transactionQuery, nil)
	mockTransactionDomainService.On("CalculateMaxCookingTime", transactionQuery.Orders).Return(45 * time.Minute)
	mockTransactionDomainService.On("GetOrderDelayStatus", 45*time.Minute, transactionQuery.Transaction.CookedAt, transactionQuery.Transaction.ServedAt).Return(true)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, "Q0001", result.QueueCode)
	assert.Equal(t, "pending", result.OrderStatus)
	assert.Equal(t, decimal.NewFromInt(65000), result.TotalPrice)
	assert.Equal(t, tableID, result.Table.ID)
	assert.Equal(t, "A1", result.Table.TableNumber)
	assert.Len(t, result.Orders, 2)

	// Check first order
	assert.Equal(t, menuID, result.Orders[0].Menu.ID)
	assert.Equal(t, "Burger", result.Orders[0].Menu.Name)
	assert.Equal(t, "25000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)

	// Check second order
	assert.Equal(t, "Fries", result.Orders[1].Menu.Name)
	assert.Equal(t, "15000", result.Orders[1].Menu.Price)
	assert.Equal(t, 1, result.Orders[1].Quantity)

	// Check delay calculation (should be delayed since cooked 30 min ago but not served)
	assert.True(t, result.IsDelayed)

	mockTransactionRepo.AssertExpectations(t)
	mockTransactionDomainService.AssertExpectations(t)
}

func TestGetTransactionByID_Success_NotDelayed(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	userID := uuid.New().String()
	transactionID := uuid.New().String()
	tableID := uuid.New().String()

	// Mock transaction schema - recently cooked and served
	cookedAt := time.Now().Add(-10 * time.Minute)
	servedAt := time.Now().Add(-5 * time.Minute)

	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(uuid.MustParse(transactionID)),
			UserID:      identity.NewID(uuid.MustParse(userID)),
			TableID:     identity.NewID(uuid.MustParse(tableID)),
			QueueCode:   transaction.QueueCode{Code: "Q0002"},
			OrderStatus: transaction.OrderStatus{Status: "served"},
			TotalPrice:  shared.Price{Price: decimal.NewFromInt(45000)},
			CookedAt:    &cookedAt,
			ServedAt:    &servedAt,
		},
		Orders: []transaction.OrderQuery{
			{
				Order: order.Order{
					ID:       identity.NewID(uuid.New()),
					Quantity: 1,
				},
				Menu: menu.Menu{
					ID:          identity.NewID(uuid.New()),
					Name:        "Pizza",
					Price:       shared.Price{Price: decimal.NewFromInt(45000)},
					CookingTime: 45 * time.Minute,
				},
			},
		},
		Table: table.Table{
			ID:          identity.NewID(uuid.MustParse(tableID)),
			TableNumber: "B2",
		},
	}

	// Set up expectations
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transactionQuery, nil)
	mockTransactionDomainService.On("CalculateMaxCookingTime", transactionQuery.Orders).Return(45 * time.Minute)
	mockTransactionDomainService.On("GetOrderDelayStatus", 45*time.Minute, &cookedAt, &servedAt).Return(false)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, "Q0002", result.QueueCode)
	assert.Equal(t, "served", result.OrderStatus)
	assert.Equal(t, decimal.NewFromInt(45000), result.TotalPrice)
	assert.Equal(t, "B2", result.Table.TableNumber)
	assert.Len(t, result.Orders, 1)
	assert.Equal(t, "Pizza", result.Orders[0].Menu.Name)
	assert.Equal(t, "45000", result.Orders[0].Menu.Price)
	assert.Equal(t, 1, result.Orders[0].Quantity)

	// Should not be delayed since served within expected time
	assert.False(t, result.IsDelayed)

	mockTransactionRepo.AssertExpectations(t)
	mockTransactionDomainService.AssertExpectations(t)
}

func TestGetTransactionByID_TransactionNotFound(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	transactionID := uuid.New().String()

	// Set up expectations - transaction not found
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transaction.Query{}, assert.AnError)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, response.Transaction{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestGetTransactionByID_InvalidTransactionType(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	transactionID := uuid.New().String()

	// Set up expectations - return error for invalid transaction
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transaction.Query{}, assert.AnError)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, response.Transaction{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestGetTransactionByID_EmptyOrders(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	userID := uuid.New().String()
	transactionID := uuid.New().String()
	tableID := uuid.New().String()

	// Mock transaction schema with no orders
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(uuid.MustParse(transactionID)),
			UserID:      identity.NewID(uuid.MustParse(userID)),
			TableID:     identity.NewID(uuid.MustParse(tableID)),
			QueueCode:   transaction.QueueCode{Code: "Q0003"},
			OrderStatus: transaction.OrderStatus{Status: "pending"},
			TotalPrice:  shared.Price{Price: decimal.NewFromInt(0)},
			CookedAt:    nil,
			ServedAt:    nil,
		},
		Orders: []transaction.OrderQuery{}, // Empty orders
		Table: table.Table{
			ID:          identity.NewID(uuid.MustParse(tableID)),
			TableNumber: "C3",
		},
	}

	// Set up expectations
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transactionQuery, nil)
	mockTransactionDomainService.On("CalculateMaxCookingTime", transactionQuery.Orders).Return(0 * time.Minute)
	mockTransactionDomainService.On("GetOrderDelayStatus", 0*time.Minute, (*time.Time)(nil), (*time.Time)(nil)).Return(false)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, "Q0003", result.QueueCode)
	assert.Equal(t, "pending", result.OrderStatus)
	assert.Equal(t, decimal.NewFromInt(0), result.TotalPrice)
	assert.Equal(t, "C3", result.Table.TableNumber)
	assert.Len(t, result.Orders, 0)
	assert.False(t, result.IsDelayed) // No cooking time to calculate delay

	mockTransactionRepo.AssertExpectations(t)
	mockTransactionDomainService.AssertExpectations(t)
}

func TestGetTransactionByID_WithDecimalPrices(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForGetByID)
	mockUserRepo := new(MockUserRepositoryForTransaction)
	mockTableRepo := new(MockTableRepositoryForTransaction)
	mockOrderRepo := new(MockOrderRepositoryForTransaction)
	mockMenuRepo := new(MockMenuRepositoryForTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForGetByID)
	mockOrderService := new(MockOrderServiceForGetByID)
	mockTransactionDomainService := new(MockTransactionDomainServiceForGetByID)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockTransactionDomainService,
		mockPaymentGateway,
		nil, // interface{} - using nil for now
		mockOrderService,
	)

	ctx := context.Background()
	userID := uuid.New().String()
	transactionID := uuid.New().String()
	tableID := uuid.New().String()

	// Mock transaction schema with decimal prices
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(uuid.MustParse(transactionID)),
			UserID:      identity.NewID(uuid.MustParse(userID)),
			TableID:     identity.NewID(uuid.MustParse(tableID)),
			QueueCode:   transaction.QueueCode{Code: "Q0004"},
			OrderStatus: transaction.OrderStatus{Status: "pending"},
			TotalPrice:  shared.Price{Price: decimal.NewFromFloat(12500.50)},
			CookedAt:    nil,
			ServedAt:    nil,
		},
		Orders: []transaction.OrderQuery{
			{
				Order: order.Order{
					ID:       identity.NewID(uuid.New()),
					Quantity: 1,
				},
				Menu: menu.Menu{
					ID:          identity.NewID(uuid.New()),
					Name:        "Premium Coffee",
					Price:       shared.Price{Price: decimal.NewFromFloat(12500.50)},
					CookingTime: 5 * time.Minute,
				},
			},
		},
		Table: table.Table{
			ID:          identity.NewID(uuid.MustParse(tableID)),
			TableNumber: "D4",
		},
	}

	// Set up expectations
	mockTransactionRepo.On("GetDetailedTransactionByID", ctx, nil, transactionID).Return(transactionQuery, nil)
	mockTransactionDomainService.On("CalculateMaxCookingTime", transactionQuery.Orders).Return(5 * time.Minute)
	mockTransactionDomainService.On("GetOrderDelayStatus", 5*time.Minute, (*time.Time)(nil), (*time.Time)(nil)).Return(false)

	// Act
	result, err := transactionService.GetTransactionByID(ctx, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, "Q0004", result.QueueCode)
	assert.Equal(t, "pending", result.OrderStatus)
	assert.Equal(t, decimal.NewFromFloat(12500.50), result.TotalPrice)
	assert.Equal(t, "D4", result.Table.TableNumber)
	assert.Len(t, result.Orders, 1)
	assert.Equal(t, "Premium Coffee", result.Orders[0].Menu.Name)
	assert.Equal(t, "12500.5", result.Orders[0].Menu.Price)
	assert.Equal(t, 1, result.Orders[0].Quantity)
	assert.False(t, result.IsDelayed) // Not cooked yet

	mockTransactionRepo.AssertExpectations(t)
	mockTransactionDomainService.AssertExpectations(t)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
