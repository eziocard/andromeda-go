package middlewares

import (
	"net/http"
	"strings"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/eziocard/andromeda-go/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "no autorizado",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "formato de token inválido",
		})
		return
	}

	userID, err := utils.ValidateAccessToken(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token inválido o expirado",
		})
		return
	}

	var user models.User

	err = initializers.DB.
		Preload("Role").
		Preload("Business").
		First(&user, userID).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "usuario no encontrado",
		})
		return
	}
	if !user.IsActive {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "cuenta inactiva",
		})
		return
	}

	c.Set("user", user)

	c.Next()
}
