package main

import (
	"fp-kpl/application/service"
	"fp-kpl/command"
	"fp-kpl/infrastructure/database/config"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/repository"
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
	transactionRepository := db_transaction.NewRepository(db)

	userRepository := repository.NewUserRepository(transactionRepository)
	tableRepository := repository.NewTableRepository(transactionRepository)
	categoryRepository := repository.NewCategoryRepository(transactionRepository)
	menuRepository := repository.NewMenuRepository(transactionRepository)

	userService := service.NewUserService(userRepository, jwtService, transactionRepository)
	tableService := service.NewTableService(tableRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	menuService := service.NewMenuService(menuRepository, categoryRepository)

	userController := controller.NewUserController(userService)
	tableController := controller.NewTableController(tableService)
	categoryController := controller.NewCategoryController(categoryService)
	menuController := controller.NewMenuController(menuService)

	defer config.CloseDatabaseConnection(db)

	if !args(db) {
		return
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	route.UserRoute(server, userController, jwtService)
	route.TableRoute(server, tableController, jwtService)
	route.CategoryRoute(server, categoryController, jwtService)
	route.MenuRoute(server, menuController, jwtService)

	run(server)
}
