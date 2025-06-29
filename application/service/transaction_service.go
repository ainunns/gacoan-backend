package service

import (
	"context"
	"fmt"
	"fp-kpl/application"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/port"
	"fp-kpl/domain/table"
	"fp-kpl/domain/transaction"
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/validation"
	"fp-kpl/platform/pagination"
	"time"

	"github.com/google/uuid"
)

type (
	TransactionService interface {
		CreateTransaction(ctx context.Context, userID string, req request.TransactionCreate) (response.TransactionCreate, error)
		HookTransaction(ctx context.Context, datas map[string]interface{}) error
		GetAllTransactionsWithPagination(ctx context.Context, userID string, req pagination.Request) (pagination.ResponseWithData, error)
		CalculateMaxCookingTime(ctx context.Context, orders []order.Order) (time.Duration, error)
	}

	transactionService struct {
		transactionRepository transaction.Repository
		userRepository        user.Repository
		tableRepository       table.Repository
		orderRepository       order.Repository
		menuRepository        menu.Repository
		paymentGatewayPort    port.PaymentGatewayPort
		transaction           interface{}
		orderService          OrderService
	}
)

func NewTransactionService(
	transactionRepository transaction.Repository,
	userRepository user.Repository,
	tableRepository table.Repository,
	orderRepository order.Repository,
	menuRepository menu.Repository,
	paymentGatewayPort port.PaymentGatewayPort,
	transaction interface{},
	orderService OrderService,
) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		tableRepository:       tableRepository,
		orderRepository:       orderRepository,
		menuRepository:        menuRepository,
		paymentGatewayPort:    paymentGatewayPort,
		transaction:           transaction,
		orderService:          orderService,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, userID string, req request.TransactionCreate) (response.TransactionCreate, error) {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = application.RecoveredFromPanic(r)
		}
		validatedTransaction.CommitOrRollback(ctx, tx, err)
	}()

	retrievedUser, err := s.userRepository.GetUserByID(ctx, tx, userID)
	if err != nil {
		return response.TransactionCreate{}, err
	}
	retrievedTable, err := s.tableRepository.GetTableByID(ctx, tx, req.TableID)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	orderStatus, err := transaction.NewOrderStatus(transaction.OrderStatusPending)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	totalPrice, err := s.orderService.CalculateTotalPrice(ctx, req.Orders)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	paymentStatus, err := transaction.NewPayment("", transaction.PaymentStatusPending)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	transactionEntity := transaction.Transaction{
		UserID:      retrievedUser.ID,
		TableID:     retrievedTable.ID,
		OrderStatus: orderStatus,
		Payment:     paymentStatus,
		TotalPrice:  totalPrice,
	}

	createdTransaction, err := s.transactionRepository.CreateTransaction(ctx, tx, transactionEntity)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	var createdOrders []response.OrderForTransactionCreate
	for _, orderItem := range req.Orders {
		retrievedMenu, err := s.menuRepository.GetMenuByID(ctx, tx, orderItem.MenuID)
		if err != nil {
			return response.TransactionCreate{}, err
		}

		orderEntity := order.Order{
			TransactionID: createdTransaction.ID,
			MenuID:        retrievedMenu.ID,
			Quantity:      orderItem.Quantity,
		}

		createdOrder, err := s.orderRepository.CreateOrder(ctx, tx, orderEntity)
		if err != nil {
			return response.TransactionCreate{}, err
		}

		createdOrders = append(createdOrders, response.OrderForTransactionCreate{
			Menu: response.MenuForTransaction{
				ID:    retrievedMenu.ID.String(),
				Name:  retrievedMenu.Name,
				Price: retrievedMenu.Price.Price.String(),
			},
			Quantity: createdOrder.Quantity,
		})
	}

	paymentURL, err := s.paymentGatewayPort.ProcessPayment(ctx, tx, createdTransaction)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	return response.TransactionCreate{
		TransactionID: createdTransaction.ID.String(),
		TotalPrice:    totalPrice.Price.String(),
		PaymentLink:   paymentURL,
		Orders:        createdOrders,
	}, nil
}

func (s *transactionService) HookTransaction(ctx context.Context, datas map[string]interface{}) error {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = application.RecoveredFromPanic(r)
		}
		validatedTransaction.CommitOrRollback(ctx, tx, err)
	}()

	transactionID, ok := datas["order_id"].(string)
	if !ok {
		return fmt.Errorf("order_id is required in datas")
	}

	err = s.paymentGatewayPort.HookPayment(ctx, tx, uuid.MustParse(transactionID), datas)
	if err != nil {
		return fmt.Errorf("failed to hook payment: %w", err)
	}

	return nil
}

func (s *transactionService) GetAllTransactionsWithPagination(ctx context.Context, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	retrievedData, err := s.transactionRepository.GetAllTransactionsWithPagination(ctx, nil, userID, req)
	if err != nil {
		return pagination.ResponseWithData{}, err
	}

	data := make([]any, 0, len(retrievedData.Data))
	for _, retrievedTransaction := range retrievedData.Data {
		transactionEntity, ok := retrievedTransaction.(transaction.Transaction)
		if !ok {
			return pagination.ResponseWithData{}, transaction.ErrorInvalidTransaction
		}

		var retrievedOrder []order.Order
		retrievedOrder, err = s.orderRepository.GetOrdersByTransactionID(ctx, nil, transactionEntity.ID.String())
		if err != nil {
			return pagination.ResponseWithData{}, order.ErrorGetOrdersByTransactionID
		}

		var orderResponses []response.OrderForTransaction
		for _, orderItem := range retrievedOrder {
			menuEntity, err := s.menuRepository.GetMenuByID(ctx, nil, orderItem.MenuID.String())
			if err != nil {
				return pagination.ResponseWithData{}, menu.ErrorGetMenuByID
			}

			orderResponses = append(orderResponses, response.OrderForTransaction{
				Menu: response.MenuForTransaction{
					ID:    menuEntity.ID.String(),
					Name:  menuEntity.Name,
					Price: menuEntity.Price.Price.String(),
				},
				Quantity: orderItem.Quantity,
			})
		}

		maxCookingTime, err := s.CalculateMaxCookingTime(ctx, retrievedOrder)
		if err != nil {
			return pagination.ResponseWithData{}, err
		}

		now := time.Now()
		isDelayed := false

		if transactionEntity.CookedAt != nil {
			expectedFinishTime := transactionEntity.CookedAt.Add(maxCookingTime)
			if transactionEntity.ServedAt != nil {
				isDelayed = transactionEntity.ServedAt.After(expectedFinishTime)
			} else {
				isDelayed = now.After(expectedFinishTime)
			}
		}

		var retrievedTable table.Table
		retrievedTable, err = s.tableRepository.GetTableByID(ctx, nil, transactionEntity.TableID.String())
		if err != nil {
			return pagination.ResponseWithData{}, table.ErrorGetTableByID
		}

		data = append(data, response.Transaction{
			ID:           transactionEntity.ID.String(),
			QueueCode:    transactionEntity.QueueCode.Code,
			EstimateTime: maxCookingTime.String(),
			Orders:       orderResponses,
			TotalPrice:   transactionEntity.TotalPrice.Price,
			Table: response.Table{
				ID:          retrievedTable.ID.String(),
				TableNumber: retrievedTable.TableNumber,
			},
			OrderStatus: transactionEntity.OrderStatus.Status,
			IsDelayed:   isDelayed,
		})
	}

	return pagination.ResponseWithData{
		Data:     data,
		Response: retrievedData.Response,
	}, nil
}

func (s *transactionService) CalculateMaxCookingTime(ctx context.Context, orders []order.Order) (time.Duration, error) {
	maxCookingTime := time.Duration(0)

	for _, orderItem := range orders {
		menuEntity, err := s.menuRepository.GetMenuByID(ctx, nil, orderItem.MenuID.String())
		if err != nil {
			return 0, menu.ErrorGetMenuByID
		}

		if menuEntity.CookingTime > maxCookingTime {
			maxCookingTime = menuEntity.CookingTime
		}
	}

	return maxCookingTime, nil
}
