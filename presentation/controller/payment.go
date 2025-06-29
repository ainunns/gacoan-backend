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
	PaymentController interface {
		PaymentNotification(ctx *gin.Context)
	}

	paymentController struct {
		paymentService service.PaymentService
	}
)

func NewPaymentController(paymentService service.PaymentService) PaymentController {
	return &paymentController{
		paymentService: paymentService,
	}
}

func (c *paymentController) PaymentNotification(ctx *gin.Context) {
	var req request.MidtransNotification

	if err := ctx.ShouldBind(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.paymentService.PaymentNotification(ctx.Request.Context(), req)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedPaymentNotification, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
}
