package controller

import (
	"fp-kpl/application/request"
	"fp-kpl/application/service"
	"fp-kpl/platform/pagination"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	TransactionController interface {
		CreateTransaction(ctx *gin.Context)
		HookTransaction(ctx *gin.Context)
		GetAllTransactionsWithPagination(ctx *gin.Context)
		GetTransactionByID(ctx *gin.Context)
		GetNextOrder(ctx *gin.Context)
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

func (t transactionController) GetAllTransactionsWithPagination(ctx *gin.Context) {
	var req pagination.Request

	if err := ctx.ShouldBindQuery(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromQuery, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := ctx.MustGet("user_id").(string)
	result, err := t.transactionService.GetAllTransactionsWithPagination(ctx.Request.Context(), userID, req)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetAllTransactions, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetAllTransactions, result.Data, result.Response)
	ctx.JSON(http.StatusOK, res)
}

func (t transactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.MustGet("user_id").(string)

	result, err := t.transactionService.GetTransactionByID(ctx.Request.Context(), userID, id)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetTransaction, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetTransaction, result)
	ctx.JSON(http.StatusOK, res)
}

func (t transactionController) GetNextOrder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	result, err := t.transactionService.GetNextOrder(ctx.Request.Context(), userID)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetNextOrder, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetNextOrder, result)
	ctx.JSON(http.StatusOK, res)
}
