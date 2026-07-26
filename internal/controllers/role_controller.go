package controllers

import (
	"net/http"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/gin-gonic/gin"
)

func RoleCreate(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	role := models.Role{Name: body.Name}

	if err := initializers.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se pudo crear el rol: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"role": role,
	})
}

func RolesIndex(c *gin.Context) {
	var roles []models.Role
	initializers.DB.Find(&roles)
	c.JSON(200, gin.H{
		"roles": roles,
	})
}
