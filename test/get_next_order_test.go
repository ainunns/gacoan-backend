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

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for GetNextOrder tests
// (Reuse structure from finish_delivering_test.go)
type MockTransactionRepositoryForNextOrder struct {
	mock.Mock
}

func (m *MockTransactionRepositoryForNextOrder) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(response.NextOrder), args.Error(1)
}

// Implement other methods as no-op for interface compliance
func (m *MockTransactionRepositoryForNextOrder) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
	return pagination.ResponseWithData{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error) {
	return nil, nil
}
func (m *MockTransactionRepositoryForNextOrder) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	return "", nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
func (m *MockTransactionRepositoryForNextOrder) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error) {
	return nil, nil
}

// Mocks for other repositories (minimal, not used in these tests)
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepository) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	return user.User{}, nil
}
func (m *MockUserRepository) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	return user.User{}, false, nil
}

type MockTableRepository struct{ mock.Mock }

func (m *MockTableRepository) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	return nil, nil
}
func (m *MockTableRepository) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	return table.Table{}, nil
}

type MockOrderRepository struct{ mock.Mock }

func (m *MockOrderRepository) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	return order.Order{}, nil
}
func (m *MockOrderRepository) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	return nil, nil
}

type MockMenuRepository struct{ mock.Mock }

func (m *MockMenuRepository) GetAllMenus(ctx context.Context, tx interface{}) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepository) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}
func (m *MockMenuRepository) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu_item.Menu, error) {
	return nil, nil
}
func (m *MockMenuRepository) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu_item.Menu, error) {
	return menu_item.Menu{}, nil
}

type MockPaymentGatewayPort struct{ mock.Mock }

func (m *MockPaymentGatewayPort) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	return port.ProcessPaymentResponse{}, nil
}
func (m *MockPaymentGatewayPort) HookPayment(ctx context.Context, tx interface{}, transactionID uuid.UUID, datas map[string]interface{}) error {
	return nil
}

func TestGetNextOrder_Success(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForNextOrder)
	mockUserRepo := new(MockUserRepository)
	mockTableRepo := new(MockTableRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockMenuRepo := new(MockMenuRepository)
	mockPaymentGateway := new(MockPaymentGatewayPort)

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
	queueCode := "Q0005"
	menuID := uuid.New()
	orderID := uuid.New()
	transactionID := uuid.New()

	// Mock menu and order
	menuSchema := &schema.Menu{
		ID:    menuID,
		Name:  "Nasi Goreng",
		Price: decimal.NewFromInt(20000),
	}
	orderSchema := schema.Order{
		ID:       orderID,
		Menu:     menuSchema,
		MenuID:   menuID,
		Quantity: 2,
	}
	transactionSchema := schema.Transaction{
		ID:        transactionID,
		QueueCode: &queueCode,
		Orders:    []schema.Order{orderSchema},
	}

	mockTransactionRepo.On("GetNextOrder", ctx, nil, userID).Return(transactionSchema, nil)

	result, err := transactionService.GetNextOrder(ctx)

	assert.NoError(t, err)
	assert.Equal(t, queueCode, result.QueueCode)
	assert.Len(t, result.Orders, 1)
	assert.Equal(t, "Nasi Goreng", result.Orders[0].Menu.Name)
	assert.Equal(t, "20000", result.Orders[0].Menu.Price)
	assert.Equal(t, 2, result.Orders[0].Quantity)

	mockTransactionRepo.AssertExpectations(t)
}

func TestGetNextOrder_NoOrder(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForNextOrder)
	mockUserRepo := new(MockUserRepository)
	mockTableRepo := new(MockTableRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockMenuRepo := new(MockMenuRepository)
	mockPaymentGateway := new(MockPaymentGatewayPort)

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

	mockTransactionRepo.On("GetNextOrder", ctx, nil, userID).Return(nil, nil)

	result, err := transactionService.GetNextOrder(ctx)

	assert.NoError(t, err)
	assert.Equal(t, response.NextOrder{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetNextOrder_InvalidType(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForNextOrder)
	mockUserRepo := new(MockUserRepository)
	mockTableRepo := new(MockTableRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockMenuRepo := new(MockMenuRepository)
	mockPaymentGateway := new(MockPaymentGatewayPort)

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

	mockTransactionRepo.On("GetNextOrder", ctx, nil, userID).Return("invalid_type", nil)

	result, err := transactionService.GetNextOrder(ctx)

	assert.Error(t, err)
	assert.Equal(t, transaction.ErrorInvalidTransaction, err)
	assert.Equal(t, response.NextOrder{}, result)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetNextOrder_RepoError(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepositoryForNextOrder)
	mockUserRepo := new(MockUserRepository)
	mockTableRepo := new(MockTableRepository)
	mockOrderRepo := new(MockOrderRepository)
	mockMenuRepo := new(MockMenuRepository)
	mockPaymentGateway := new(MockPaymentGatewayPort)

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

	repoErr := assert.AnError
	mockTransactionRepo.On("GetNextOrder", ctx, nil, userID).Return(nil, repoErr)

	result, err := transactionService.GetNextOrder(ctx)

	assert.Error(t, err)
	assert.Equal(t, response.NextOrder{}, result)
	mockTransactionRepo.AssertExpectations(t)
}
