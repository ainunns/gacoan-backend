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
		GetTransactionByID(ctx context.Context, id string) (response.Transaction, error)
		GetAllReadyToServeTransactionList(ctx context.Context, req pagination.Request) (pagination.ResponseWithData, error)
		GetNextOrder(ctx context.Context) (response.NextOrder, error)
		StartCooking(ctx context.Context, req request.StartCooking) (response.StartCooking, error)
		FinishCooking(ctx context.Context, req request.FinishCooking) (response.FinishCooking, error)
		StartDelivering(ctx context.Context, req request.StartDelivering) (response.StartDelivering, error)
		FinishDelivering(ctx context.Context, req request.FinishDelivering) (response.FinishDelivering, error)
	}

	transactionService struct {
		transactionRepository    transaction.Repository
		userRepository           user.Repository
		tableRepository          table.Repository
		orderRepository          order.Repository
		menuRepository           menu.Repository
		transactionDomainService transaction.Service
		paymentGatewayPort       port.PaymentGatewayPort
		transaction              interface{}
		orderService             OrderService
	}
)

func NewTransactionService(
	transactionRepository transaction.Repository,
	userRepository user.Repository,
	tableRepository table.Repository,
	orderRepository order.Repository,
	menuRepository menu.Repository,
	transactionDomainService transaction.Service,
	paymentGatewayPort port.PaymentGatewayPort,
	transaction interface{},
	orderService OrderService,
) TransactionService {
	return &transactionService{
		transactionRepository:    transactionRepository,
		userRepository:           userRepository,
		tableRepository:          tableRepository,
		orderRepository:          orderRepository,
		menuRepository:           menuRepository,
		transactionDomainService: transactionDomainService,
		paymentGatewayPort:       paymentGatewayPort,
		transaction:              transaction,
		orderService:             orderService,
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

	payment, err := s.paymentGatewayPort.ProcessPayment(ctx, tx, createdTransaction)
	if err != nil {
		return response.TransactionCreate{}, err
	}

	return response.TransactionCreate{
		TransactionID: createdTransaction.ID.String(),
		TotalPrice:    totalPrice.Price.String(),
		Token:         payment.Token,
		PaymentLink:   payment.PaymentLink,
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
		transactionQuery, ok := retrievedTransaction.(transaction.Query)
		if !ok {
			return pagination.ResponseWithData{}, transaction.ErrorInvalidTransaction
		}

		var orderResponses []response.OrderForTransaction
		for _, orderQuery := range transactionQuery.Orders {
			orderResponses = append(orderResponses, response.OrderForTransaction{
				Menu: response.MenuForTransaction{
					ID:    orderQuery.Menu.ID.String(),
					Name:  orderQuery.Menu.Name,
					Price: orderQuery.Menu.Price.Price.String(),
				},
				Quantity: orderQuery.Order.Quantity,
			})
		}

		maxCookingTime := s.transactionDomainService.CalculateMaxCookingTime(transactionQuery.Orders)

		now := time.Now()
		isDelayed := false

		if transactionQuery.Transaction.CookedAt != nil {
			expectedFinishTime := transactionQuery.Transaction.CookedAt.Add(maxCookingTime)
			if transactionQuery.Transaction.ServedAt != nil {
				isDelayed = transactionQuery.Transaction.ServedAt.After(expectedFinishTime)
			} else {
				isDelayed = now.After(expectedFinishTime)
			}
		}

		data = append(data, response.Transaction{
			ID:           transactionQuery.Transaction.ID.String(),
			QueueCode:    transactionQuery.Transaction.QueueCode.Code,
			EstimateTime: maxCookingTime.String(),
			Orders:       orderResponses,
			TotalPrice:   transactionQuery.Transaction.TotalPrice.Price,
			Table: response.Table{
				ID:          transactionQuery.Table.ID.String(),
				TableNumber: transactionQuery.Table.TableNumber,
			},
			OrderStatus: transactionQuery.Transaction.OrderStatus.Status,
			IsDelayed:   isDelayed,
		})
	}

	return pagination.ResponseWithData{
		Data:     data,
		Response: retrievedData.Response,
	}, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string) (response.Transaction, error) {
	retrievedData, err := s.transactionRepository.GetDetailedTransactionByID(ctx, nil, id)
	if err != nil {
		return response.Transaction{}, err
	}

	var orderResponses []response.OrderForTransaction
	for _, orderQuery := range retrievedData.Orders {
		orderResponses = append(orderResponses, response.OrderForTransaction{
			Menu: response.MenuForTransaction{
				ID:    orderQuery.Menu.ID.String(),
				Name:  orderQuery.Menu.Name,
				Price: orderQuery.Menu.Price.Price.String(),
			},
			Quantity: orderQuery.Order.Quantity,
		})
	}

	maxCookingTime := s.transactionDomainService.CalculateMaxCookingTime(retrievedData.Orders)

	now := time.Now()
	isDelayed := false

	if retrievedData.Transaction.CookedAt != nil {
		expectedFinishTime := retrievedData.Transaction.CookedAt.Add(maxCookingTime)
		if retrievedData.Transaction.ServedAt != nil {
			isDelayed = retrievedData.Transaction.ServedAt.After(expectedFinishTime)
		} else {
			isDelayed = now.After(expectedFinishTime)
		}
	}

	return response.Transaction{
		ID:           retrievedData.Transaction.ID.String(),
		QueueCode:    retrievedData.Transaction.QueueCode.Code,
		EstimateTime: maxCookingTime.String(),
		Orders:       orderResponses,
		OrderStatus:  retrievedData.Transaction.OrderStatus.Status,
		TotalPrice:   retrievedData.Transaction.TotalPrice.Price,
		Table: response.Table{
			ID:          retrievedData.Table.ID.String(),
			TableNumber: retrievedData.Table.TableNumber,
		},
		IsDelayed: isDelayed,
	}, nil
}

func (s *transactionService) GetAllReadyToServeTransactionList(ctx context.Context, req pagination.Request) (pagination.ResponseWithData, error) {
	retrievedData, err := s.transactionRepository.GetAllReadyToServeTransactionList(ctx, nil, req)
	if err != nil {
		return pagination.ResponseWithData{}, err
	}

	data := make([]any, 0, len(retrievedData.Data))
	for _, retrievedTransaction := range retrievedData.Data {
		transactionQuery, ok := retrievedTransaction.(transaction.Query)
		if !ok {
			return pagination.ResponseWithData{}, transaction.ErrorInvalidTransaction
		}

		var orderResponses []response.OrderForWaiter
		for _, orderQuery := range transactionQuery.Orders {
			orderResponses = append(orderResponses, response.OrderForWaiter{
				Menu: response.MenuForWaiter{
					ID:   orderQuery.Menu.ID.String(),
					Name: orderQuery.Menu.Name,
				},
				Quantity: orderQuery.Order.Quantity,
			})
		}

		tableResponse := response.TableForWaiter{
			ID:          transactionQuery.Table.ID.String(),
			TableNumber: transactionQuery.Table.TableNumber,
		}

		data = append(data, response.TransactionForWaiter{
			QueueCode: transactionQuery.Transaction.QueueCode.Code,
			Orders:    orderResponses,
			Table:     tableResponse,
		})
	}

	return pagination.ResponseWithData{
		Data:     data,
		Response: retrievedData.Response,
	}, nil
}

func (s *transactionService) GetNextOrder(ctx context.Context) (response.NextOrder, error) {
	retrievedNextOrder, err := s.transactionRepository.GetNextOrder(ctx, nil)
	if err != nil {
		return response.NextOrder{}, err
	}

	if retrievedNextOrder.QueueCode == "" {
		return response.NextOrder{}, transaction.ErrorNextOrderNotFound
	}

	return retrievedNextOrder, nil
}

func (s *transactionService) StartCooking(ctx context.Context, req request.StartCooking) (response.StartCooking, error) {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return response.StartCooking{}, err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return response.StartCooking{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = application.RecoveredFromPanic(r)
		}
		validatedTransaction.CommitOrRollback(ctx, tx, err)
	}()

	retrievedData, err := s.transactionRepository.GetTransactionByQueueCode(ctx, tx, req.QueueCode)
	if err != nil {
		return response.StartCooking{}, err
	}

	if retrievedData.Transaction.OrderStatus.Status != transaction.OrderStatusPending {
		return response.StartCooking{}, transaction.ErrorInvalidOrderStatus
	}

	_, err = s.transactionRepository.UpdateTransactionCookingStatusStart(ctx, tx, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.StartCooking{}, err
	}

	_, err = s.transactionRepository.UpdateCookedAt(ctx, tx, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.StartCooking{}, err
	}

	var orderResponses []response.OrderForTransaction
	for _, orderQuery := range retrievedData.Orders {
		orderResponses = append(orderResponses, response.OrderForTransaction{
			Menu: response.MenuForTransaction{
				ID:    orderQuery.Menu.ID.String(),
				Name:  orderQuery.Menu.Name,
				Price: orderQuery.Menu.Price.Price.String(),
			},
			Quantity: orderQuery.Order.Quantity,
		})
	}

	return response.StartCooking{
		QueueCode: retrievedData.Transaction.QueueCode.Code,
		Orders:    orderResponses,
	}, nil
}

func (s *transactionService) FinishCooking(ctx context.Context, req request.FinishCooking) (response.FinishCooking, error) {
	retrievedData, err := s.transactionRepository.GetTransactionByQueueCode(ctx, nil, req.QueueCode)
	if err != nil {
		return response.FinishCooking{}, err
	}

	if retrievedData.Transaction.OrderStatus.Status != transaction.OrderStatusPreparing {
		return response.FinishCooking{}, transaction.ErrorInvalidOrderStatus
	}

	_, err = s.transactionRepository.UpdateTransactionCookingStatusFinish(ctx, nil, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.FinishCooking{}, err
	}

	var orderResponses []response.OrderForTransaction
	for _, orderQuery := range retrievedData.Orders {
		orderResponses = append(orderResponses, response.OrderForTransaction{
			Menu: response.MenuForTransaction{
				ID:    orderQuery.Menu.ID.String(),
				Name:  orderQuery.Menu.Name,
				Price: orderQuery.Menu.Price.Price.String(),
			},
			Quantity: orderQuery.Order.Quantity,
		})
	}

	return response.FinishCooking{
		QueueCode: retrievedData.Transaction.QueueCode.Code,
		Orders:    orderResponses,
	}, nil
}

func (s *transactionService) StartDelivering(ctx context.Context, req request.StartDelivering) (response.StartDelivering, error) {
	retrievedData, err := s.transactionRepository.GetTransactionByQueueCode(ctx, nil, req.QueueCode)
	if err != nil {
		return response.StartDelivering{}, err
	}

	if retrievedData.Transaction.OrderStatus.Status != transaction.OrderStatusReadyToServe {
		return response.StartDelivering{}, transaction.ErrorInvalidOrderStatus
	}

	_, err = s.transactionRepository.UpdateTransactionDeliveringStatusStart(ctx, nil, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.StartDelivering{}, err
	}

	var orderResponses []response.OrderForTransaction
	for _, orderQuery := range retrievedData.Orders {
		orderResponses = append(orderResponses, response.OrderForTransaction{
			Menu: response.MenuForTransaction{
				ID:    orderQuery.Menu.ID.String(),
				Name:  orderQuery.Menu.Name,
				Price: orderQuery.Menu.Price.Price.String(),
			},
			Quantity: orderQuery.Order.Quantity,
		})
	}

	return response.StartDelivering{
		QueueCode: retrievedData.Transaction.QueueCode.Code,
		Orders:    orderResponses,
	}, nil
}

func (s *transactionService) FinishDelivering(ctx context.Context, req request.FinishDelivering) (response.FinishDelivering, error) {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return response.FinishDelivering{}, err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return response.FinishDelivering{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = application.RecoveredFromPanic(r)
		}
		validatedTransaction.CommitOrRollback(ctx, tx, err)
	}()

	retrievedData, err := s.transactionRepository.GetTransactionByQueueCode(ctx, nil, req.QueueCode)
	if err != nil {
		return response.FinishDelivering{}, err
	}

	if retrievedData.Transaction.OrderStatus.Status != transaction.OrderStatusDelivering {
		return response.FinishDelivering{}, transaction.ErrorInvalidOrderStatus
	}

	_, err = s.transactionRepository.UpdateTransactionDeliveringStatusFinish(ctx, nil, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.FinishDelivering{}, err
	}

	_, err = s.transactionRepository.UpdateServedAt(ctx, tx, retrievedData.Transaction.ID.String())
	if err != nil {
		return response.FinishDelivering{}, err
	}

	return response.FinishDelivering{}, nil
}
