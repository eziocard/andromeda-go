package routes

import (
	"github.com/eziocard/andromeda-go/internal/controllers"
	"github.com/eziocard/andromeda-go/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", controllers.UserCreate)
		userGroup.POST("/login", controllers.UserLogin)
		userGroup.POST("/refresh", controllers.UserRefresh)
		userGroup.POST("/logout", controllers.UserLogout)
	}

	userGroup.Use(
		middlewares.RequireAuth,
		middlewares.RequireRole("superuser"),
	)

	userGroup.GET("/me", controllers.UserMe)
	userGroup.GET("", controllers.UsersIndex)

}
