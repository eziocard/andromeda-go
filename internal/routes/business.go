package routes

import (
	"github.com/eziocard/andromeda-go/internal/controllers"
	"github.com/eziocard/andromeda-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func BusinessRoutes(r *gin.Engine) {
	businessGroup := r.Group("/business")
	{

	}
	businessGroup.Use(
		middlewares.RequireAuth,
		middlewares.RequireRole("superuser"),
	)

	businessGroup.POST("", controllers.BusinessCreate)
	businessGroup.GET("", controllers.BusinessIndex)
	businessGroup.PATCH("/:id", controllers.BusinessUpdate)

}
