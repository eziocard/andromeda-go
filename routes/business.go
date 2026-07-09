package routes

import (
	"github.com/eziocard/andromeda-go/controllers"
	"github.com/gin-gonic/gin"
)

func BusinessRoutes(r *gin.Engine) {
	businessGroup := r.Group("/business")
	{
		businessGroup.POST("", controllers.BusinessCreate)
	}
}
