package controller

import (
	"fp-kpl/application/request"
	"fp-kpl/application/service"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	OrderController interface {
		CalculateTotalPrice(ctx *gin.Context)
	}

	orderController struct {
		orderService service.OrderService
	}
)

func NewOrderController(orderService service.OrderService) OrderController {
	return &orderController{
		orderService: orderService,
	}
}

func (c *orderController) CalculateTotalPrice(ctx *gin.Context) {
	var req request.CalculateTotalPrice
	if err := ctx.ShouldBind(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	totalPrice, err := c.orderService.CalculateTotalPrice(ctx.Request.Context(), req.Orders)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedCalculateTotalPrice, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessCalculateTotalPrice, totalPrice)
	ctx.JSON(http.StatusOK, res)
}
