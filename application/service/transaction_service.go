package service

import (
	"context"
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
)

type (
	TransactionService interface {
		CreateTransaction(ctx context.Context, userID string, req request.TransactionCreate) (response.TransactionCreate, error)
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

	transactionEntity := transaction.Transaction{
		UserID:      retrievedUser.ID,
		TableID:     retrievedTable.ID,
		OrderStatus: orderStatus,
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
			Menu: response.MenuForTransactionCreate{
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
