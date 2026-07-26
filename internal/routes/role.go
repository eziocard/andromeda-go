package routes

import (
	"github.com/eziocard/andromeda-go/internal/controllers"
	"github.com/eziocard/andromeda-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine) {
	roles := r.Group("/roles")
	roles.Use(
		middlewares.RequireAuth,
		middlewares.RequireRole("superuser"),
	)

	roles.GET("", controllers.RolesIndex)
	roles.POST("", controllers.RoleCreate)
}
