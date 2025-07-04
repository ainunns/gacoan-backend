package test

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	"fp-kpl/domain/identity"
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

// Mock repository for GetAllReadyToServeTransactionList
type MockTransactionRepositoryForReadyToServe struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForReadyToServe) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

func (m *MockTransactionRepositoryForReadyToServe) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	return response.NextOrder{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (transaction.Query, error) {
	return transaction.Query{}, nil
}
func (m *MockTransactionRepositoryForReadyToServe) GetDetailedTransactionByID(ctx context.Context, tx interface{}, id string) (transaction.Query, error) {
	return transaction.Query{}, nil
}

// Minimal mocks for other repositories
type MockUserRepositoryForReadyToServe struct{ mock.Mock }

func (m *MockUserRepositoryForReadyToServe) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForReadyToServe) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForReadyToServe) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForReadyToServe) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepositoryForReadyToServe struct{ mock.Mock }

func (m *MockTableRepositoryForReadyToServe) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepositoryForReadyToServe) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepositoryForReadyToServe struct{ mock.Mock }

func (m *MockOrderRepositoryForReadyToServe) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepositoryForReadyToServe) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockMenuRepositoryForReadyToServe struct{ mock.Mock }

func (m *MockMenuRepositoryForReadyToServe) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForReadyToServe) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepositoryForReadyToServe) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForReadyToServe) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockPaymentGatewayPortForReadyToServe struct{ mock.Mock }

func (m *MockPaymentGatewayPortForReadyToServe) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	return port.ProcessPaymentResponse{}, nil
}
func (m *MockPaymentGatewayPortForReadyToServe) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	return nil
}

func TestGetAllReadyToServeTransactionList_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForReadyToServe)
	mockUserRepo := new(MockUserRepositoryForReadyToServe)
	mockTableRepo := new(MockTableRepositoryForReadyToServe)
	mockOrderRepo := new(MockOrderRepositoryForReadyToServe)
	mockMenuRepo := new(MockMenuRepositoryForReadyToServe)
	mockPaymentGateway := new(MockPaymentGatewayPortForReadyToServe)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	tableID := uuid.New()
	menuID := uuid.New()
	transactionID := uuid.New()
	queueCode := "Q0010"

	query := transaction.Query{
		Transaction: transaction.Transaction{
			ID:        identity.NewIDFromSchema(transactionID),
			QueueCode: transaction.QueueCode{Code: queueCode, Valid: true},
			TableID:   identity.NewIDFromSchema(tableID),
		},
		Table: table.Table{
			ID:          identity.NewIDFromSchema(tableID),
			TableNumber: "A1",
		},
		Orders: []transaction.OrderQuery{
			{
				Menu: menu_item.Menu{
					ID:    identity.NewIDFromSchema(menuID),
					Name:  "Sate Ayam",
					Price: shared.NewPriceFromSchema(decimal.NewFromInt(30000)),
				},
				Order: order.Order{
					MenuID:   identity.NewIDFromSchema(menuID),
					Quantity: 3,
				},
			},
		},
	}

	paginated := pagination.ResponseWithData{
		Data:     []any{query},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 1},
	}

	mockTransactionRepo.On("GetAllReadyToServeTransactionList", ctx, nil, mock.Anything).Return(paginated, nil)

	req := pagination.Request{Page: 1, PerPage: 10}
	result, err := transactionService.GetAllReadyToServeTransactionList(ctx, req)

	assert.NoError(t, err)
	if assert.Len(t, result.Data, 1) {
		transactionResp, ok := result.Data[0].(response.TransactionForWaiter)
		assert.True(t, ok)
		assert.Equal(t, queueCode, transactionResp.QueueCode)
		assert.Equal(t, "A1", transactionResp.Table.TableNumber)
		assert.Len(t, transactionResp.Orders, 1)
		assert.Equal(t, "Sate Ayam", transactionResp.Orders[0].Menu.Name)
		assert.Equal(t, 3, transactionResp.Orders[0].Quantity)
	}

	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllReadyToServeTransactionList_Empty(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForReadyToServe)
	mockUserRepo := new(MockUserRepositoryForReadyToServe)
	mockTableRepo := new(MockTableRepositoryForReadyToServe)
	mockOrderRepo := new(MockOrderRepositoryForReadyToServe)
	mockMenuRepo := new(MockMenuRepositoryForReadyToServe)
	mockPaymentGateway := new(MockPaymentGatewayPortForReadyToServe)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	paginated := pagination.ResponseWithData{
		Data:     []any{},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 1},
	}

	mockTransactionRepo.On("GetAllReadyToServeTransactionList", ctx, nil, mock.Anything).Return(paginated, nil)

	req := pagination.Request{Page: 1, PerPage: 10}
	result, err := transactionService.GetAllReadyToServeTransactionList(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, result.Data, 0)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllReadyToServeTransactionList_InvalidType(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForReadyToServe)
	mockUserRepo := new(MockUserRepositoryForReadyToServe)
	mockTableRepo := new(MockTableRepositoryForReadyToServe)
	mockOrderRepo := new(MockOrderRepositoryForReadyToServe)
	mockMenuRepo := new(MockMenuRepositoryForReadyToServe)
	mockPaymentGateway := new(MockPaymentGatewayPortForReadyToServe)

	transactionService := service.NewTransactionService(
		mockTransactionRepo,
		mockUserRepo,
		mockTableRepo,
		mockOrderRepo,
		mockMenuRepo,
		nil,
		mockPaymentGateway,
		nil,
		nil,
	)

	ctx := context.Background()
	paginated := pagination.ResponseWithData{
		Data:     []any{"invalid_type"},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 1},
	}

	mockTransactionRepo.On("GetAllReadyToServeTransactionList", ctx, nil, mock.Anything).Return(paginated, nil)

	req := pagination.Request{Page: 1, PerPage: 10}
	result, err := transactionService.GetAllReadyToServeTransactionList(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, pagination.ResponseWithData{}, result)
	mockTransactionRepo.AssertExpectations(t)
}
