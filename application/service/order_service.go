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
		orderRepository    order.Repository
		menuRepository     menu.Repository
		orderDomainService order.Service
	}
)

func NewOrderService(
	orderRepository order.Repository,
	menuRepository menu.Repository,
	orderDomainService order.Service,
) OrderService {
	return &orderService{
		orderRepository:    orderRepository,
		menuRepository:     menuRepository,
		orderDomainService: orderDomainService,
	}
}

func (s *orderService) CalculateTotalPrice(ctx context.Context, orders []request.Order) (shared.Price, error) {
	totalPrice := decimal.NewFromInt(0)

	for _, orderItem := range orders {
		menuEntity, err := s.menuRepository.GetMenuByID(ctx, nil, orderItem.MenuID)
		if err != nil {
			return shared.Price{}, menu.ErrorMenuNotFound
		}

		menuPrice := menuEntity.Price

		orderPrice, err := s.orderDomainService.CalculatePrice(ctx, menuPrice, int64(orderItem.Quantity))
		if err != nil {
			return shared.Price{}, err
		}

		totalPrice = totalPrice.Add(orderPrice.Price)
	}

	return shared.NewPrice(totalPrice)
}
