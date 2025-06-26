package route

import (
	"fp-kpl/application/service"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"

	"github.com/gin-gonic/gin"
)

func TableRoute(route *gin.Engine, tableController controller.TableController, jwtService service.JWTService) {
	tableGroup := route.Group("/api/table")
	{
		tableGroup.GET("/", middleware.Authenticate(jwtService), tableController.GetAllTables)
		tableGroup.GET("/:id", middleware.Authenticate(jwtService), tableController.GetTableByID)
	}
}
