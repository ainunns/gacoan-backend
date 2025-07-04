package test

import (
	"context"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
	menu "fp-kpl/domain/menu/menu_item"
	menu_item "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/shared"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
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

func (m *MockTransactionRepositoryForFinishCooking) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (transaction.Query, error) {
	args := m.Called(ctx, tx, queueCode)
	if args.Get(0) == nil {
		return transaction.Query{}, args.Error(1)
	}
	return args.Get(0).(transaction.Query), args.Error(1)
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
	args := m.Called(ctx, tx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepositoryForFinishCooking) GetDetailedTransactionByID(ctx context.Context, tx interface{}, id string) (transaction.Query, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(transaction.Query), args.Error(1)
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

type MockOrderServiceForFinishCooking struct{ mock.Mock }

func (m *MockOrderServiceForFinishCooking) CalculateTotalPrice(ctx context.Context, orders []request.Order) (shared.Price, error) {
	args := m.Called(ctx, orders)
	return args.Get(0).(shared.Price), args.Error(1)
}

func TestFinishCooking_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForFinishCooking)
	mockUserRepo := new(MockUserRepositoryForFinishCooking)
	mockTableRepo := new(MockTableRepositoryForFinishCooking)
	mockOrderRepo := new(MockOrderRepositoryForFinishCooking)
	mockMenuRepo := new(MockMenuRepositoryForFinishCooking)
	mockPaymentGateway := new(MockPaymentGatewayPortForFinishCooking)
	mockTransactionInterface := new(MockTransactionInterfaceForFinishCooking)
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()

	// Create a proper transaction.Query
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(transactionID),
			OrderStatus: transaction.OrderStatus{Status: transaction.OrderStatusPreparing},
			QueueCode:   transaction.QueueCode{Code: queueCode},
		},
		Orders: []transaction.OrderQuery{}, // Empty orders for this test
	}

	transactionEntity := transaction.Transaction{
		ID: identity.NewID(transactionID),
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionQuery, nil)
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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, assert.AnError)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(nil, assert.AnError)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.Error(t, err)
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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()

	// Create transaction with wrong status
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(transactionID),
			OrderStatus: transaction.OrderStatus{Status: transaction.OrderStatusPending}, // Wrong status
			QueueCode:   transaction.QueueCode{Code: queueCode},
		},
		Orders: []transaction.OrderQuery{},
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionQuery, nil)

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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()

	// Create transaction with correct status
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(transactionID),
			OrderStatus: transaction.OrderStatus{Status: transaction.OrderStatusPreparing},
			QueueCode:   transaction.QueueCode{Code: queueCode},
		},
		Orders: []transaction.OrderQuery{},
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionQuery, nil)
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
	mockOrderService := new(MockOrderServiceForFinishCooking)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil, // transaction.Service - using nil for now
		mockPaymentGateway,
		mockTransactionInterface,
		mockOrderService,
	)

	ctx := context.Background()
	queueCode := "Q0001"
	transactionID := uuid.New()

	// Create transaction with multiple orders
	transactionQuery := transaction.Query{
		Transaction: transaction.Transaction{
			ID:          identity.NewID(transactionID),
			OrderStatus: transaction.OrderStatus{Status: transaction.OrderStatusPreparing},
			QueueCode:   transaction.QueueCode{Code: queueCode},
		},
		Orders: []transaction.OrderQuery{
			{
				Order: order.Order{
					ID:       identity.NewID(uuid.New()),
					Quantity: 2,
				},
				Menu: menu.Menu{
					ID:    identity.NewID(uuid.New()),
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
	}

	transactionEntity := transaction.Transaction{
		ID: identity.NewID(transactionID),
	}

	mockTransactionRepo.On("GetTransactionByQueueCode", ctx, nil, queueCode).Return(transactionQuery, nil)
	mockTransactionRepo.On("UpdateTransactionCookingStatusFinish", ctx, nil, transactionID.String()).Return(transactionEntity, nil)

	req := request.FinishCooking{QueueCode: queueCode}
	result, err := transactionService.FinishCooking(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	assert.Len(t, result.Orders, 2)
	assert.Equal(t, "Burger", result.Orders[0].Menu.Name)
	assert.Equal(t, "25000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)
	assert.Equal(t, "Fries", result.Orders[1].Menu.Name)
	assert.Equal(t, "15000", result.Orders[1].Menu.Price)
	assert.Equal(t, 1, result.Orders[1].Quantity)
	mockTransactionRepo.AssertExpectations(t)
}
