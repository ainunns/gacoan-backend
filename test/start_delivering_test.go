package test

import (
	"context"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
	menu_item "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/platform/pagination"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for StartDelivering tests
type MockTransactionRepositoryForStartDelivering struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForStartDelivering) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	args := m.Called(ctx, tx, queueCode)
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepositoryForStartDelivering) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

// Implement other methods as no-op for interface compliance
func (m *MockTransactionRepositoryForStartDelivering) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForStartDelivering) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForStartDelivering) GetNextOrder(ctx context.Context, tx interface{}, userID string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForStartDelivering) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForStartDelivering) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}

// Mocks for other repositories (minimal, not used in these tests)
type MockUserRepositoryForStartDelivering struct{ mock.Mock }

func (m *MockUserRepositoryForStartDelivering) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForStartDelivering) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForStartDelivering) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForStartDelivering) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepositoryForStartDelivering struct{ mock.Mock }

func (m *MockTableRepositoryForStartDelivering) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepositoryForStartDelivering) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepositoryForStartDelivering struct{ mock.Mock }

func (m *MockOrderRepositoryForStartDelivering) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepositoryForStartDelivering) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockMenuRepositoryForStartDelivering struct{ mock.Mock }

func (m *MockMenuRepositoryForStartDelivering) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForStartDelivering) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepositoryForStartDelivering) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForStartDelivering) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockPaymentGatewayPortForStartDelivering struct{ mock.Mock }

func (m *MockPaymentGatewayPortForStartDelivering) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	return port.ProcessPaymentResponse{}, nil
}
func (m *MockPaymentGatewayPortForStartDelivering) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	return nil
}

func TestStartDelivering_Success(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()
	menuID := uuid.New()
	orderID := uuid.New()

	// Mock transaction schema with ready_to_serve status
	menuSchema := &schema.Menu{
		ID:    menuID,
		Name:  "Nasi Goreng",
		Price: decimal.NewFromInt(25000),
	}
	orderSchema := schema.Order{
		ID:       orderID,
		Menu:     menuSchema,
		MenuID:   menuID,
		Quantity: 2,
	}
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusReadyToServe,
		QueueCode:   &queueCode,
		Orders:      []schema.Order{orderSchema},
	}

	// Mock transaction entity for repository calls
	transactionEntity := transaction.Transaction{
		ID: identity.NewIDFromSchema(transactionID),
	}

	// Set up expectations
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionDeliveringStatusStart", ctx, nil, transactionID.String()).Return(transactionEntity, nil)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	assert.Len(t, result.Orders, 1)
	assert.Equal(t, "Nasi Goreng", result.Orders[0].Menu.Name)
	assert.Equal(t, "25000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_TransactionNotFound(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Set up expectations - transaction not found
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, nil)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.StartDelivering{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_InvalidTransactionType(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Set up expectations - return invalid type
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return("invalid_type", nil)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrorInvalidTransaction, err)
	assert.Equal(t, response.StartDelivering{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_InvalidOrderStatus(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()

	// Mock transaction schema with wrong status (should be ready_to_serve)
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusPending, // Wrong status
		QueueCode:   &queueCode,
	}

	// Set up expectations
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrorInvalidOrderStatus, err)
	assert.Equal(t, response.StartDelivering{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_GetTransactionError(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	// Set up expectations - repository error
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, assert.AnError)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, response.StartDelivering{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_UpdateStatusError(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()
	menuID := uuid.New()
	orderID := uuid.New()

	// Mock transaction schema with ready_to_serve status
	menuSchema := &schema.Menu{
		ID:    menuID,
		Name:  "Nasi Goreng",
		Price: decimal.NewFromInt(25000),
	}
	orderSchema := schema.Order{
		ID:       orderID,
		Menu:     menuSchema,
		MenuID:   menuID,
		Quantity: 2,
	}
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusReadyToServe,
		QueueCode:   &queueCode,
		Orders:      []schema.Order{orderSchema},
	}

	// Set up expectations - update status fails
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionDeliveringStatusStart", ctx, nil, transactionID.String()).Return(transaction.Transaction{}, assert.AnError)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, response.StartDelivering{}, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestStartDelivering_WithMultipleOrders(t *testing.T) {
	// Arrange
	mockTransactionRepo := new(MockTransactionRepositoryForStartDelivering)
	mockUserRepo := new(MockUserRepositoryForStartDelivering)
	mockTableRepo := new(MockTableRepositoryForStartDelivering)
	mockOrderRepo := new(MockOrderRepositoryForStartDelivering)
	mockMenuRepo := new(MockMenuRepositoryForStartDelivering)
	mockPaymentGateway := new(MockPaymentGatewayPortForStartDelivering)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()
	menuID1 := uuid.New()
	menuID2 := uuid.New()
	orderID1 := uuid.New()
	orderID2 := uuid.New()

	// Mock transaction schema with multiple orders
	menuSchema1 := &schema.Menu{
		ID:    menuID1,
		Name:  "Nasi Goreng",
		Price: decimal.NewFromInt(25000),
	}
	menuSchema2 := &schema.Menu{
		ID:    menuID2,
		Name:  "Es Teh",
		Price: decimal.NewFromInt(5000),
	}
	orderSchema1 := schema.Order{
		ID:       orderID1,
		Menu:     menuSchema1,
		MenuID:   menuID1,
		Quantity: 2,
	}
	orderSchema2 := schema.Order{
		ID:       orderID2,
		Menu:     menuSchema2,
		MenuID:   menuID2,
		Quantity: 1,
	}
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusReadyToServe,
		QueueCode:   &queueCode,
		Orders:      []schema.Order{orderSchema1, orderSchema2},
	}

	// Mock transaction entity for repository calls
	transactionEntity := transaction.Transaction{
		ID: identity.NewIDFromSchema(transactionID),
	}

	// Set up expectations
	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionDeliveringStatusStart", ctx, nil, transactionID.String()).Return(transactionEntity, nil)

	// Act
	req := request.StartDelivering{QueueCode: queueCode}
	result, err := transactionService.StartDelivering(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	assert.Len(t, result.Orders, 2)
	assert.Equal(t, "Nasi Goreng", result.Orders[0].Menu.Name)
	assert.Equal(t, "25000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)
	assert.Equal(t, "Es Teh", result.Orders[1].Menu.Name)
	assert.Equal(t, "5000", result.Orders[1].Menu.Price)
	assert.Equal(t, 1, result.Orders[1].Quantity)

	mockTransactionRepo.AssertExpectations(t)
}
