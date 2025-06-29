package service

import (
	"context"
	"fp-kpl/application/request"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/shared"
	"github.com/shopspring/decimal"
)

type (
	OrderService interface {
		CalculateTotalPrice(ctx context.Context, orders []request.Order) (shared.Price, error)
	}

	orderService struct {
		orderRepository order.Repository
		menuRepository  menu.Repository
	}
)

func NewOrderService(
	orderRepository order.Repository,
	menuRepository menu.Repository,
) OrderService {
	return &orderService{
		orderRepository: orderRepository,
		menuRepository:  menuRepository,
	}
}

func (s *orderService) CalculateTotalPrice(ctx context.Context, orders []request.Order) (shared.Price, error) {
	totalPrice := decimal.NewFromInt(0)

	for _, orderItem := range orders {
		menuEntity, err := s.menuRepository.GetMenuByID(ctx, nil, orderItem.MenuID)
		if err != nil {
			return shared.Price{}, menu.ErrorMenuNotFound
		}

		menuPrice := menuEntity.Price.Price

		if orderItem.Quantity <= 0 {
			return shared.Price{}, order.ErrorInvalidQuantity
		}

		orderPrice := menuPrice.Mul(decimal.NewFromInt(int64(orderItem.Quantity)))
		totalPrice = totalPrice.Add(orderPrice)
	}

	return shared.NewPrice(totalPrice)
}
