package routes

import (
	"github.com/eziocard/andromeda-go/internal/controllers"
	"github.com/eziocard/andromeda-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SaleRoutes(r *gin.Engine) {
	saleGroup := r.Group("/sales")
	saleGroup.Use(middlewares.RequireAuth)
	{
		saleGroup.POST("", controllers.SaleCreate)
		saleGroup.GET("", controllers.SaleIndex)
		saleGroup.POST("/:id/void", controllers.SaleVoid)
	}
}
