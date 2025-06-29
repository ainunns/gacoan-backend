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
	TransactionController interface {
		CreateTransaction(ctx *gin.Context)
	}

	transactionController struct {
		transactionService service.TransactionService
	}
)

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (t transactionController) CreateTransaction(ctx *gin.Context) {
	var req request.TransactionCreate
	if err := ctx.ShouldBind(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := ctx.MustGet("user_id").(string)
	result, err := t.transactionService.CreateTransaction(ctx.Request.Context(), userID, req)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedCreateTransaction, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessCreateTransaction, result)
	ctx.JSON(http.StatusCreated, res)
}
