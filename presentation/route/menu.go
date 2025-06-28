package route

import (
	"fp-kpl/application/service"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"

	"github.com/gin-gonic/gin"
)

func MenuRoute(route *gin.Engine, menuController controller.MenuController, jwtService service.JWTService) {
	menuGroup := route.Group("/api/menu")
	{
		menuGroup.GET("/", middleware.Authenticate(jwtService), menuController.GetAllMenus)
		menuGroup.GET("/:id", middleware.Authenticate(jwtService), menuController.GetMenuByID)
	}
}
