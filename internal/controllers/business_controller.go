package controllers

import (
	"net/http"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/models"
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

func BusinessIndex(c *gin.Context) {
	var business []models.Business
	initializers.DB.Find(&business)

	c.JSON(200, gin.H{
		"business": business,
	})
}

func BusinessUpdate(c *gin.Context) {
	id := c.Param("id")

	var business models.Business
	if err := initializers.DB.First(&business, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "negocio no encontrado",
		})
		return
	}

	var body struct {
		Name     *string `json:"name"`
		Address  *string `json:"address"`
		City     *string `json:"city"`
		Contact  *string `json:"contact"`
		Email    *string `json:"email" binding:"omitempty,email"`
		IsActive *bool   `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updates := map[string]interface{}{}

	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Address != nil {
		updates["address"] = *body.Address
	}
	if body.City != nil {
		updates["city"] = *body.City
	}
	if body.Contact != nil {
		updates["contact"] = *body.Contact
	}
	if body.Email != nil {
		updates["email"] = *body.Email
	}
	if body.IsActive != nil {
		updates["is_active"] = *body.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se enviaron campos para actualizar",
		})
		return
	}

	if err := initializers.DB.Model(&business).Updates(updates).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se pudo actualizar el negocio: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"business": business,
	})
}
