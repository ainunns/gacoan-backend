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
		transactionGroup.GET("/", middleware.Authenticate(jwtService), transactionController.GetAllTransactionsWithPagination)
		transactionGroup.GET("/:id", middleware.Authenticate(jwtService), transactionController.GetTransactionByID)
		transactionGroup.POST("/hook", transactionController.HookTransaction)

		// Kitchen
		transactionGroup.GET("/next-order", middleware.Authenticate(jwtService), transactionController.GetNextOrder)
		transactionGroup.POST("/start-cooking", middleware.Authenticate(jwtService), transactionController.StartCooking)
		transactionGroup.POST("/finish-cooking", middleware.Authenticate(jwtService), transactionController.FinishCooking)

		// Waiter
		transactionGroup.GET("/ready-to-serve", middleware.Authenticate(jwtService), transactionController.GetAllReadyToServeTransactionList)
		transactionGroup.POST("/start-delivering", middleware.Authenticate(jwtService), transactionController.StartDelivering)
		transactionGroup.POST("/finish-delivering", middleware.Authenticate(jwtService), transactionController.FinishDelivering)
	}
}
