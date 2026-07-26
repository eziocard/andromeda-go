package controllers

import (
	"net/http"

	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/gin-gonic/gin"
)

// helper: extrae el usuario autenticado y confirma que tenga un negocio asignado

func getAuthUserBusiness(c *gin.Context) (models.User, uint, bool) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return models.User{}, 0, false
	}
	user, ok := userValue.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno de autenticación"})
		return models.User{}, 0, false
	}
	if user.BusinessID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "el usuario no tiene un negocio asignado"})
		return models.User{}, 0, false
	}
	return user, *user.BusinessID, true
}
