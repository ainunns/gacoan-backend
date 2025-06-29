package middleware

import (
	"fp-kpl/application/service"
	"fp-kpl/domain/user"
	"fp-kpl/presentation"
	"fp-kpl/presentation/message"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorize(userService service.UserService, allowedRoles []user.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("user_id").(string)
		thisUser, err := userService.GetUserByID(ctx.Request.Context(), userID)
		if err != nil {
			response := presentation.BuildResponseFailed(message.FailedProcessRequest, message.FailedDeniedAccess, nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		authorized := false
		for _, role := range allowedRoles {
			if thisUser.Role == role.Name {
				authorized = true
				break
			}
		}

		if !authorized {
			response := presentation.BuildResponseFailed(message.FailedProcessRequest, message.FailedDeniedAccess, nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		ctx.Next()
	}
}
