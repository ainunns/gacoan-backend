package route

import (
	"fp-kpl/application/service"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"

	"github.com/gin-gonic/gin"
)

func CategoryRoute(route *gin.Engine, categoryController controller.CategoryController, jwtService service.JWTService) {
	categoryGroup := route.Group("/api/category")
	{
		categoryGroup.GET("/", middleware.Authenticate(jwtService), categoryController.GetAllCategories)
		categoryGroup.GET("/:id", middleware.Authenticate(jwtService), categoryController.GetCategoryByID)
	}
}
