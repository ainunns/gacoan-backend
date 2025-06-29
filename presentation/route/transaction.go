package route

import (
	"fp-kpl/application/service"
	"fp-kpl/domain/user"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"

	"github.com/gin-gonic/gin"
)

func TransactionRoute(route *gin.Engine, transactionController controller.TransactionController, jwtService service.JWTService, userService service.UserService) {
	transactionGroup := route.Group("/api/transaction")
	{
		transactionGroup.POST("/",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleCustomer},
				{user.RoleWaiter},
				{user.RoleSuperAdmin},
			}),
			transactionController.CreateTransaction)
		transactionGroup.GET("/", middleware.Authenticate(jwtService), transactionController.GetAllTransactionsWithPagination)
		transactionGroup.GET("/:id", middleware.Authenticate(jwtService), transactionController.GetTransactionByID)
		transactionGroup.POST("/hook", transactionController.HookTransaction)

		// Kitchen
		transactionGroup.GET("/next-order",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleKitchen},
				{user.RoleSuperAdmin},
			}),
			transactionController.GetNextOrder)
		transactionGroup.POST("/start-cooking",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleKitchen},
				{user.RoleSuperAdmin},
			}),
			transactionController.StartCooking)
		transactionGroup.POST("/finish-cooking",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleKitchen},
				{user.RoleSuperAdmin},
			}),
			transactionController.FinishCooking)

		// Waiter
		transactionGroup.GET("/ready-to-serve",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleWaiter},
				{user.RoleSuperAdmin},
			}),
			transactionController.GetAllReadyToServeTransactionList)
		transactionGroup.POST("/start-delivering",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleWaiter},
				{user.RoleSuperAdmin},
			}),
			transactionController.StartDelivering)
		transactionGroup.POST("/finish-delivering",
			middleware.Authenticate(jwtService),
			middleware.Authorize(userService, []user.Role{
				{user.RoleWaiter},
				{user.RoleSuperAdmin},
			}),
			transactionController.FinishDelivering)
	}
}
