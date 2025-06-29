package route

import (
	"fp-kpl/application/service"
	"fp-kpl/domain/user"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"

	"github.com/gin-gonic/gin"
)

func MenuRoute(route *gin.Engine, menuController controller.MenuController, jwtService service.JWTService, userService service.UserService) {
	menuGroup := route.Group("/api/menu")
	{
		menuGroup.GET("/", middleware.Authenticate(jwtService), menuController.GetAllMenus)
		menuGroup.GET("/:id", middleware.Authenticate(jwtService), menuController.GetMenuByID)
		menuGroup.PATCH("/:id/availability",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{Name: user.RoleKitchen},
				{Name: user.RoleSuperAdmin},
			}),
			menuController.UpdateMenuAvailability)
	}
}
