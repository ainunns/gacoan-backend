package controller

import (
	"errors"
	"fp-kpl/application/response"
	"fp-kpl/application/service"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	MenuController interface {
		GetAllMenus(ctx *gin.Context)
		GetMenuByID(ctx *gin.Context)
	}

	menuController struct {
		menuService service.MenuService
	}
)

func NewMenuController(menuService service.MenuService) MenuController {
	return &menuController{menuService: menuService}
}

func (c *menuController) GetAllMenus(ctx *gin.Context) {
	categoryID := ctx.Query("category_id")

	var err error
	var allMenus []response.Menu

	if categoryID != "" {
		allMenus, err = c.menuService.GetMenusByCategoryID(ctx.Request.Context(), categoryID)
	} else {
		allMenus, err = c.menuService.GetAllMenus(ctx.Request.Context())
	}

	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetAllMenus, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	menus := make([]interface{}, len(allMenus))
	for i, menu := range allMenus {
		menus[i] = menu
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetAllMenus, menus)
	ctx.JSON(http.StatusOK, res)
}

func (c *menuController) GetMenuByID(ctx *gin.Context) {
	id := ctx.Param("id")
	responseMenu, err := c.menuService.GetMenuByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, menu.ErrorMenuNotFound) {
			res := presentation.BuildResponseFailed(message.FailedGetMenu, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}

		res := presentation.BuildResponseFailed(message.FailedGetMenu, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetMenu, responseMenu)
	ctx.JSON(http.StatusOK, res)
}
