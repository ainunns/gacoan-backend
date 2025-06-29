package controller

import (
	"fp-kpl/application/request"
	"fp-kpl/application/service"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	TransactionController interface {
		CreateTransaction(ctx *gin.Context)
		HookTransaction(ctx *gin.Context)
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

func (t transactionController) HookTransaction(ctx *gin.Context) {
	var datas map[string]interface{}
	if err := ctx.ShouldBindJSON(&datas); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := t.transactionService.HookTransaction(ctx.Request.Context(), datas)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedHookTransaction, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessHookTransaction, nil)
	ctx.JSON(http.StatusOK, res)
}
