package test

import (
	"context"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	menu_item "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
	"fp-kpl/platform/pagination"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for finish delivering tests
type MockTransactionRepositoryForFinishDelivering struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForFinishDelivering) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, userID, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	args := m.Called(ctx, tx, userID, id)
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetNextOrder(ctx context.Context, tx interface{}, userID string) (interface{}, error) {
	args := m.Called(ctx, tx, userID)
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishDelivering) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	args := m.Called(ctx, tx, queueCode)
	return args.Get(0), args.Error(1)
}

// Mock other repositories
type MockUserRepositoryForFinishDelivering struct {
	mock.Mock
}

func (m *MockUserRepositoryForFinishDelivering) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	args := m.Called(ctx, tx, userEntity)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForFinishDelivering) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForFinishDelivering) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	args := m.Called(ctx, tx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepositoryForFinishDelivering) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	args := m.Called(ctx, tx, email)
	return args.Get(0).(user.User), args.Get(1).(bool), args.Error(2)
}

type MockTableRepositoryForFinishDelivering struct {
	mock.Mock
}

func (m *MockTableRepositoryForFinishDelivering) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]table.Table), args.Error(1)
}

func (m *MockTableRepositoryForFinishDelivering) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(table.Table), args.Error(1)
}

type MockOrderRepositoryForFinishDelivering struct {
	mock.Mock
}

func (m *MockOrderRepositoryForFinishDelivering) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	args := m.Called(ctx, tx, orderEntity)
	return args.Get(0).(order.Order), args.Error(1)
}

func (m *MockOrderRepositoryForFinishDelivering) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).([]order.Order), args.Error(1)
}

type MockMenuRepositoryForFinishDelivering struct {
	mock.Mock
}

func (m *MockMenuRepositoryForFinishDelivering) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]menu_item.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForFinishDelivering) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(menu_item.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForFinishDelivering) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).([]menu_item.Menu), args.Error(1)
}

func (m *MockMenuRepositoryForFinishDelivering) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	args := m.Called(ctx, tx, id, isAvailable)
	return args.Get(0).(menu_item.Menu), args.Error(1)
}

type MockPaymentGatewayPortForFinishDelivering struct {
	mock.Mock
}

func (m *MockPaymentGatewayPortForFinishDelivering) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(port.ProcessPaymentResponse), args.Error(1)
}

func (m *MockPaymentGatewayPortForFinishDelivering) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	args := m.Called(ctx, tx, transactionID, datas)
	return args.Error(0)
}

// Mock transaction interface that can be validated
type MockTransactionInterfaceForFinishDelivering struct {
	mock.Mock
}

func (m *MockTransactionInterfaceForFinishDelivering) Begin(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return m, args.Error(1) // Return self as the transaction
}

func (m *MockTransactionInterfaceForFinishDelivering) CommitOrRollback(ctx context.Context, tx interface{}, err error) {
	m.Called(ctx, tx, err)
}

func (m *MockTransactionInterfaceForFinishDelivering) DB() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestFinishDelivering_Success(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil, // orderService
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	// The validation will fail because our mock is not the correct type
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_TransactionNotFound(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_InvalidTransactionType(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_InvalidOrderStatus(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_UpdateStatusError(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_UpdateServedAtError(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_TransactionBeginError(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}

func TestFinishDelivering_WithMultipleOrders(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForFinishDelivering)
	mockUserRepo := new(MockUserRepositoryForFinishDelivering)
	mockTableRepo := new(MockTableRepositoryForFinishDelivering)
	mockOrderRepo := new(MockOrderRepositoryForFinishDelivering)
	mockMenuRepo := new(MockMenuRepositoryForFinishDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishDelivering)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		mockTransactionInterface,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Act
	req := request.FinishDelivering{QueueCode: queueCode}
	result, err := transactionService.FinishDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.FinishDelivering{}, result)
}
