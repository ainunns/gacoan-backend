package test

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	menu_item "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/platform/pagination"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repository for GetAllTransactionsWithPagination
type MockTransactionRepositoryForPagination struct {
	mock.Mock
}

type MockDBTransaction struct{}

func (m *MockDBTransaction) Begin(ctx context.Context) interface{} {
	// Return a dummy non-nil object, or a *gorm.DB if needed
	return &gorm.DB{}
}

func (m *MockTransactionRepositoryForPagination) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	args := m.Called(ctx, tx, userID, req)
	return args.Get(0).(pagination.ResponseWithData), args.Error(1)
}

// Implement other methods as no-op for interface compliance
func (m *MockTransactionRepositoryForPagination) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForPagination) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForPagination) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	return response.NextOrder{}, nil
}
func (m *MockTransactionRepositoryForPagination) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForPagination) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForPagination) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	return nil, nil
}

// Minimal mocks for other repositories
type MockUserRepositoryForPagination struct{ mock.Mock }

func (m *MockUserRepositoryForPagination) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForPagination) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForPagination) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepositoryForPagination) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepositoryForPagination struct{ mock.Mock }

func (m *MockTableRepositoryForPagination) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepositoryForPagination) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepositoryForPagination struct{ mock.Mock }

func (m *MockOrderRepositoryForPagination) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepositoryForPagination) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockMenuRepositoryForPagination struct{ mock.Mock }

func (m *MockMenuRepositoryForPagination) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForPagination) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepositoryForPagination) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepositoryForPagination) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockPaymentGatewayPortForPagination struct{ mock.Mock }

func (m *MockPaymentGatewayPortForPagination) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	return port.ProcessPaymentResponse{}, nil
}
func (m *MockPaymentGatewayPortForPagination) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	return nil
}

func TestGetAllTransactionsWithPagination_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForPagination)
	mockUserRepo := new(MockUserRepositoryForPagination)
	mockTableRepo := new(MockTableRepositoryForPagination)
	mockOrderRepo := new(MockOrderRepositoryForPagination)
	mockMenuRepo := new(MockMenuRepositoryForPagination)
	mockPaymentGateway := new(MockPaymentGatewayPortForPagination)

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
	userID := uuid.New().String()
	tableID := uuid.New()
	menuID := uuid.New()
	orderID := uuid.New()
	transactionID := uuid.New()
	queueCode := "Q0011"

	menuSchema := &schema.Menu{
		ID:    menuID,
		Name:  "Ayam Bakar",
		Price: decimal.NewFromInt(35000),
	}
	orderSchema := schema.Order{
		ID:       orderID,
		Menu:     menuSchema,
		MenuID:   menuID,
		Quantity: 1,
	}
	transactionSchema := schema.Transaction{
		ID:        transactionID,
		QueueCode: &queueCode,
		Table: &schema.Table{
			ID:          tableID,
			TableNumber: "A2",
		},
		Orders:    []schema.Order{orderSchema},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	paginated := pagination.ResponseWithData{
		Data:     []any{transactionSchema},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 1},
	}

	req := pagination.Request{Page: 1, PerPage: 10}
	mockTransactionRepo.On("GetAllTransactionsWithPagination", ctx, nil, userID, req).Return(paginated, nil)

	result, err := transactionService.GetAllTransactionsWithPagination(ctx, userID, req)

	assert.NoError(t, err)
	assert.Len(t, result.Data, 1)
	transactionResp, ok := result.Data[0].(response.Transaction)
	assert.True(t, ok)
	assert.Equal(t, queueCode, transactionResp.QueueCode)
	assert.Equal(t, "A2", transactionResp.Table.TableNumber)
	assert.Len(t, transactionResp.Orders, 1)
	assert.Equal(t, "Ayam Bakar", transactionResp.Orders[0].Menu.Name)
	assert.Equal(t, 1, transactionResp.Orders[0].Quantity)

	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllTransactionsWithPagination_Empty(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForPagination)
	mockUserRepo := new(MockUserRepositoryForPagination)
	mockTableRepo := new(MockTableRepositoryForPagination)
	mockOrderRepo := new(MockOrderRepositoryForPagination)
	mockMenuRepo := new(MockMenuRepositoryForPagination)
	mockPaymentGateway := new(MockPaymentGatewayPortForPagination)

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
	userID := uuid.New().String()
	req := pagination.Request{Page: 1, PerPage: 10}
	paginated := pagination.ResponseWithData{
		Data:     []any{},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 0},
	}

	mockTransactionRepo.On("GetAllTransactionsWithPagination", ctx, nil, userID, req).Return(paginated, nil)

	result, err := transactionService.GetAllTransactionsWithPagination(ctx, userID, req)

	assert.NoError(t, err)
	assert.Len(t, result.Data, 0)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllTransactionsWithPagination_InvalidType(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForPagination)
	mockUserRepo := new(MockUserRepositoryForPagination)
	mockTableRepo := new(MockTableRepositoryForPagination)
	mockOrderRepo := new(MockOrderRepositoryForPagination)
	mockMenuRepo := new(MockMenuRepositoryForPagination)
	mockPaymentGateway := new(MockPaymentGatewayPortForPagination)

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
	userID := uuid.New().String()
	req := pagination.Request{Page: 1, PerPage: 10}
	paginated := pagination.ResponseWithData{
		Data:     []any{"invalid_type"},
		Response: pagination.Response{Page: 1, PerPage: 10, MaxPage: 1, Count: 1},
	}

	mockTransactionRepo.On("GetAllTransactionsWithPagination", ctx, nil, userID, req).Return(paginated, nil)

	result, err := transactionService.GetAllTransactionsWithPagination(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, pagination.ResponseWithData{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllTransactionsWithPagination_RepoError(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForPagination)
	mockUserRepo := new(MockUserRepositoryForPagination)
	mockTableRepo := new(MockTableRepositoryForPagination)
	mockOrderRepo := new(MockOrderRepositoryForPagination)
	mockMenuRepo := new(MockMenuRepositoryForPagination)
	mockPaymentGateway := new(MockPaymentGatewayPortForPagination)

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
	userID := uuid.New().String()
	req := pagination.Request{Page: 1, PerPage: 10}
	repoErr := assert.AnError

	mockTransactionRepo.On("GetAllTransactionsWithPagination", ctx, nil, userID, req).Return(pagination.ResponseWithData{}, repoErr)

	result, err := transactionService.GetAllTransactionsWithPagination(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, pagination.ResponseWithData{}, result)
	mockTransactionRepo.AssertExpectations(t)
}
