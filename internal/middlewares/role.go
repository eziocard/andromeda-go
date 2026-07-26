package middlewares

import (
	"net/http"

	"github.com/eziocard/andromeda-go/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.CurrentUser(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "usuario no autenticado",
			})
			return
		}

		for _, role := range roles {
			if user.Role.Name == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no tienes permisos",
		})
	}
}
