package routes

import (
	"github.com/eziocard/andromeda-go/controllers"
	"github.com/gin-gonic/gin"
)

func SaleRoutes(r *gin.Engine) {
	saleGroup := r.Group("/sales")
	{
		saleGroup.POST("", controllers.SaleCreate)
		saleGroup.GET("", controllers.SaleIndex)
	}
}
