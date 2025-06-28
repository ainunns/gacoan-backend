package controller

import (
	"errors"
	"fp-kpl/application/service"
	"fp-kpl/domain/menu/category"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	CategoryController interface {
		GetAllCategories(ctx *gin.Context)
		GetCategoryByID(ctx *gin.Context)
	}

	categoryController struct {
		categoryService service.CategoryService
	}
)

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &categoryController{categoryService: categoryService}
}

func (c *categoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.categoryService.GetAllCategories(ctx.Request.Context())
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetAllCategories, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetAllCategories, categories)
	ctx.JSON(http.StatusOK, res)
}

func (c *categoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	responseCategory, err := c.categoryService.GetCategoryByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, category.ErrorCategoryNotFound) {
			res := presentation.BuildResponseFailed(message.FailedGetCategory, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}

		res := presentation.BuildResponseFailed(message.FailedGetCategory, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetCategory, responseCategory)
	ctx.JSON(http.StatusOK, res)
}
