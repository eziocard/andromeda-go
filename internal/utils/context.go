package utils

import (
	"errors"

	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) (models.User, error) {
	value, exists := c.Get("user")
	if !exists {
		return models.User{}, errors.New("usuario no encontrado en el contexto")
	}

	user, ok := value.(models.User)
	if !ok {
		return models.User{}, errors.New("tipo de usuario inválido en el contexto")
	}

	return user, nil
}
