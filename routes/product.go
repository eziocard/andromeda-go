package routes

import (
	"github.com/eziocard/andromeda-go/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	productGroup := r.Group("/products")
	{
		productGroup.POST("", controllers.ProductCreate)
		productGroup.GET("", controllers.ProductIndex)
		productGroup.GET("/:id", controllers.ProductShow)
		productGroup.GET("/barcode/:barcode", controllers.ProductShowByBarcode)
		productGroup.PATCH("/:barcode", controllers.ProductsUpdate)
		productGroup.DELETE("/:id", controllers.ProductsDelete)
	}
}
