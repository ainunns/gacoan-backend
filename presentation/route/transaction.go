package route

import (
	"fp-kpl/application/service"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"
	"github.com/gin-gonic/gin"
)

func TransactionRoute(route *gin.Engine, transactionController controller.TransactionController, jwtService service.JWTService) {
	transactionGroup := route.Group("/api/transaction")
	{
		transactionGroup.POST("/", middleware.Authenticate(jwtService), transactionController.CreateTransaction)
	}
}
