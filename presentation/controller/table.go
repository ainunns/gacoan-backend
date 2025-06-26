package controller

import (
	"errors"
	"fp-kpl/application/service"
	"fp-kpl/domain/table"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	TableController interface {
		GetAllTables(ctx *gin.Context)
		GetTableByID(ctx *gin.Context)
	}

	tableController struct {
		tableService service.TableService
	}
)

func NewTableController(tableService service.TableService) TableController {
	return &tableController{tableService: tableService}
}

func (c *tableController) GetAllTables(ctx *gin.Context) {
	tables, err := c.tableService.GetAllTables(ctx.Request.Context())
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetAllTables, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetAllTables, tables)
	ctx.JSON(http.StatusOK, res)
}

func (c *tableController) GetTableByID(ctx *gin.Context) {
	id := ctx.Param("id")
	responseTable, err := c.tableService.GetTableByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, table.ErrorTableNotFound) {
			res := presentation.BuildResponseFailed(message.FailedGetTable, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}

		res := presentation.BuildResponseFailed(message.FailedGetTable, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetTable, responseTable)
	ctx.JSON(http.StatusOK, res)
}
