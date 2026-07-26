package routes

import (
	"github.com/eziocard/andromeda-go/internal/controllers"
	"github.com/eziocard/andromeda-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	productGroup := r.Group("/products")
	{

	}

	protected := productGroup.Group("")
	protected.Use(middlewares.RequireAuth)

	{
		protected.POST("", controllers.ProductCreate)
		protected.GET("", controllers.ProductIndex)
		protected.GET("/:id", controllers.ProductShow)
		protected.GET("/barcode/:barcode", controllers.ProductShowByBarcode)
		protected.PATCH("/:barcode", controllers.ProductsUpdate)
		protected.DELETE("/:id", controllers.ProductsDelete)
	}
}
