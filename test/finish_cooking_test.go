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

// Mock repositories for FinishCooking tests
type MockTransactionRepositoryForFinishCooking struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForFinishCooking) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	args := m.Called(ctx, tx, queueCode)
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishCooking) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

// Implement other methods as no-op for interface compliance
func (m *MockTransactionRepositoryForFinishCooking) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForFinishCooking) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForFinishCooking) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	return response.NextOrder{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForFinishCooking) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}

// Mocks for other repositories (minimal, not used in these tests)
type MockUserRepositoryForFinishCooking struct{ mock.Mock }

func (m *MockUserRepositoryForFinishCooking) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForFinishCooking) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForFinishCooking) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForFinishCooking) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepositoryForFinishCooking struct{ mock.Mock }

func (m *MockTableRepositoryForFinishCooking) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepositoryForFinishCooking) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepositoryForFinishCooking struct{ mock.Mock }

func (m *MockOrderRepositoryForFinishCooking) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepositoryForFinishCooking) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockMenuRepositoryForFinishCooking struct{ mock.Mock }

func (m *MockMenuRepositoryForFinishCooking) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForFinishCooking) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepositoryForFinishCooking) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForFinishCooking) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockPaymentGatewayPortForFinishCooking struct{ mock.Mock }

func (m *MockPaymentGatewayPortForFinishCooking) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	return port.ProcessPaymentResponse{}, nil
}
func (m *MockPaymentGatewayPortForFinishCooking) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	return nil
}

type MockTransactionInterfaceForFinishCooking struct {
	mock.Mock
}

func (m *MockTransactionInterfaceForFinishCooking) Begin(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return m, args.Error(1)
}

func (m *MockTransactionInterfaceForFinishCooking) CommitOrRollback(ctx context.Context, tx interface{}, err error) {
	m.Called(ctx, tx, err)
}

func (m *MockTransactionInterfaceForFinishCooking) DB() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestFinishCooking_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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
	transactionID := uuid.New()

	// Create a proper schema.Transaction
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusPreparing,
		QueueCode:   &queueCode,
		Orders:      []schema.Order{}, // Empty orders for this test
	}

	transactionEntity := transaction.Transaction{
		ID: identity.NewID(transactionID),
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionCookingStatusFinish", ctx, nil, transactionID.String()).Return(transactionEntity, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_TransactionNotFound(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, response.FinishCooking{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_InvalidTransactionType(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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

	// Return invalid type
	invalidData := "invalid data"

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(invalidData, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, transaction.ErrorInvalidTransaction, err)
	assert.Equal(t, response.FinishCooking{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_InvalidOrderStatus(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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
	transactionID := uuid.New()

	// Create transaction with wrong status
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusPending, // Wrong status
		QueueCode:   &queueCode,
		Orders:      []schema.Order{},
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, transaction.ErrorInvalidOrderStatus, err)
	assert.Equal(t, response.FinishCooking{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_GetTransactionError(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, assert.AnError)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Equal(t, response.FinishCooking{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_UpdateStatusError(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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
	transactionID := uuid.New()

	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusPreparing,
		QueueCode:   &queueCode,
		Orders:      []schema.Order{},
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionCookingStatusFinish", ctx, nil, transactionID.String()).Return(transaction.Transaction{}, assert.AnError)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Equal(t, response.FinishCooking{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFinishCooking_WithMultipleOrders(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)

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
	transactionID := uuid.New()
	menuID1 := uuid.New()
	menuID2 := uuid.New()

	// Create transaction with multiple orders
	transactionSchema := schema.Transaction{
		ID:          transactionID,
		OrderStatus: transaction.OrderStatusPreparing,
		QueueCode:   &queueCode,
		Orders: []schema.Order{
			{
				ID:            uuid.New(),
				TransactionID: transactionID,
				MenuID:        menuID1,
				Quantity:      2,
				Menu: &schema.Menu{
					ID:    menuID1,
					Name:  "Nasi Goreng",
					Price: decimal.NewFromFloat(15000),
				},
			},
			{
				ID:            uuid.New(),
				TransactionID: transactionID,
				MenuID:        menuID2,
				Quantity:      1,
				Menu: &schema.Menu{
					ID:    menuID2,
					Name:  "Es Teh",
					Price: decimal.NewFromFloat(5000),
				},
			},
		},
	}

	transactionEntity := transaction.Transaction{
		ID: identity.NewID(transactionID),
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionSchema, nil)
	mockTransactionRepo.On("UpdateTransactionCookingStatusFinish", ctx, nil, transactionID.String()).Return(transactionEntity, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	assert.Len(t, result.Orders, 2)
	assert.Equal(t, "Nasi Goreng", result.Orders[0].Menu.Name)
	assert.Equal(t, "15000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)
	assert.Equal(t, "Es Teh", result.Orders[1].Menu.Name)
	assert.Equal(t, "5000", result.Orders[1].Menu.Price)
	assert.Equal(t, 1, result.Orders[1].Quantity)
	mockTransactionRepo.AssertExpectations(t)
}
