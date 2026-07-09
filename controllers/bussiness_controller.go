package controllers

import (
	"net/http"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/models"
	"github.com/gin-gonic/gin"
)

func BusinessCreate(c *gin.Context) {
	var body struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
		City    string `json:"city" binding:"required"`
		Contact string `json:"contact"`
		Email   string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	business := models.Business{
		Name:     body.Name,
		Address:  body.Address,
		City:     body.City,
		Contact:  body.Contact,
		Email:    body.Email,
		IsActive: true,
	}

	if err := initializers.DB.Create(&business).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se pudo crear el negocio: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"business": business,
	})
}
