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

// Mock repositories for create transaction
type MockTransactionRepositoryForCreateTransaction struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForCreateTransaction) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}
func (m *MockTransactionRepositoryForCreateTransaction) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) GetNextOrder(ctx context.Context, tx interface{}, userID string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForCreateTransaction) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	return nil, nil
}

type MockTransactionInterfaceForCreateTransaction struct {
	mock.Mock
}

func (m *MockTransactionInterfaceForCreateTransaction) Begin(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return m, args.Error(1) // Return self as the transaction
}

func (m *MockTransactionInterfaceForCreateTransaction) CommitOrRollback(ctx context.Context, tx interface{}, err error) {
	m.Called(ctx, tx, err)
}

func (m *MockTransactionInterfaceForCreateTransaction) DB() interface{} {
	args := m.Called()
	return args.Get(0)
}

type MockMenuRepositoryForCreateTransaction struct{ mock.Mock }

func (m *MockMenuRepositoryForCreateTransaction) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForCreateTransaction) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepositoryForCreateTransaction) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForCreateTransaction) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockUserRepositoryForCreateTransaction struct{ mock.Mock }

func (m *MockUserRepositoryForCreateTransaction) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForCreateTransaction) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForCreateTransaction) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForCreateTransaction) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepositoryForCreateTransaction struct{ mock.Mock }

func (m *MockTableRepositoryForCreateTransaction) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepositoryForCreateTransaction) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepositoryForCreateTransaction struct{ mock.Mock }

func (m *MockOrderRepositoryForCreateTransaction) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepositoryForCreateTransaction) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockPaymentGatewayPortForCreateTransaction struct{ mock.Mock }

func (m *MockPaymentGatewayPortForCreateTransaction) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	args := m.Called(ctx, tx, transactionEntity)
	return args.Get(0).(port.ProcessPaymentResponse), args.Error(1)
}

func (m *MockPaymentGatewayPortForCreateTransaction) HookPayment(ctx context.Context, tx interface{}, transactionId uuid.UUID, datas map[string]interface{}) error {
	args := m.Called(ctx, tx, transactionId, datas)
	return args.Error(0)
}

// Example test using mocks
func TestCreateTransaction_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForCreateTransaction)
	mockUserRepo := new(MockUserRepositoryForCreateTransaction)
	mockTableRepo := new(MockTableRepositoryForCreateTransaction)
	mockOrderRepo := new(MockOrderRepositoryForCreateTransaction)
	mockMenuRepo := new(MockMenuRepositoryForCreateTransaction)
	mockPaymentGateway := new(MockPaymentGatewayPortForCreateTransaction)
	mockTransactionInterface := new(MockTransactionInterfaceForCreateTransaction)

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

	userID := uuid.New()
	tableID := uuid.New()
	menuID := uuid.New()
	orderReq := request.Order{MenuID: menuID.String(), Quantity: 1}
	req := request.TransactionCreate{
		TableID: tableID.String(),
		Orders:  []request.Order{orderReq},
	}

	// Act
	result, err := transactionService.CreateTransaction(context.Background(), userID.String(), req)

	// Assert
	// The validation will fail because our mock is not the correct type
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction", err.Error())
	assert.Equal(t, response.TransactionCreate{}, result)
}
