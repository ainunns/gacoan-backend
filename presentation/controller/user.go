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
	UserController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		Me(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var req request.UserRegister
	if err := ctx.ShouldBind(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	result, err := c.userService.Register(ctx.Request.Context(), req)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedRegister, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessRegister, result)
	ctx.JSON(http.StatusCreated, res)
}

func (c *userController) Login(ctx *gin.Context) {
	var req request.UserLogin
	if err := ctx.ShouldBind(&req); err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetDataFromBody, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.Verify(ctx.Request.Context(), req)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedLogin, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessLogin, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Me(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	result, err := c.userService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		res := presentation.BuildResponseFailed(message.FailedGetUser, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := presentation.BuildResponseSuccess(message.SuccessGetUser, result)
	ctx.JSON(http.StatusOK, res)
}
