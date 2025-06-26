package main

import (
	"fp-kpl/application/service"
	"fp-kpl/command"
	"fp-kpl/infrastructure/database/config"
	infrastructure_table "fp-kpl/infrastructure/database/table"
	"fp-kpl/infrastructure/database/transaction"
	infrastructure_user "fp-kpl/infrastructure/database/user"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"
	"fp-kpl/presentation/route"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func args(db *gorm.DB) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(db)
		return flag
	}

	return true
}

func run(server *gin.Engine) {
	server.Static("/assets", "./assets")

	if os.Getenv("IS_LOGGER") == "true" {
		route.LoggerRoute(server)
	}

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func main() {
	db := config.SetUpDatabaseConnection()

	jwtService := service.NewJWTService()

	transactionRepository := transaction.NewRepository(db)
	userRepository := infrastructure_user.NewRepository(transactionRepository)
	tableRepository := infrastructure_table.NewRepository(transactionRepository)

	userService := service.NewUserService(userRepository, jwtService, transactionRepository)
	tableService := service.NewTableService(tableRepository)

	userController := controller.NewUserController(userService)
	tableController := controller.NewTableController(tableService)

	defer config.CloseDatabaseConnection(db)

	if !args(db) {
		return
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	route.UserRoute(server, userController, jwtService)
	route.TableRoute(server, tableController, jwtService)

	run(server)
}
