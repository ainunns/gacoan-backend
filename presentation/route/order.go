package route

import (
	"fp-kpl/application/service"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoute(route *gin.Engine, orderController controller.OrderController, jwtService service.JWTService) {
	orderGroup := route.Group("/api/order")
	{
		orderGroup.POST("/calculate-total-price", middleware.Authenticate(jwtService), orderController.CalculateTotalPrice)
	}
}
