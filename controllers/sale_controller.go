package controllers

import (
	"net/http"

	"github.com/eziocard/andromeda-go/dto"
	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/models"
	"github.com/gin-gonic/gin"
)

// Crear registro de venta

func SaleCreate(c *gin.Context) {
	var body dto.SaleCreateInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	db := initializers.DB

	// Calculo del total

	var total uint
	for _, item := range body.Details {
		var product models.Product
		if err := db.First(&product, item.Product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product no encontrado: no encontrado: " + err.Error(),
			})
			return
		}
		total += product.Price * item.Quantity
	}

	// Guardar total en tabla Sale
	sale := models.Sale{Total: total}
	if err := db.Create(&sale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Guardar metodos de pago y montos
	for _, p := range body.Payments {
		payment := models.SalePayment{
			SaleID: sale.ID,
			Method: p.Method,
			Amount: p.Amount,
		}
		if err := db.Create(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// Guardando detalle de ventas
	for _, item := range body.Details {
		var product models.Product
		if err := db.First(&product, item.Product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product no encontrado: " + err.Error(),
			})
			return
		}

		detail := models.SaleDetail{
			SaleID:      sale.ID,
			ProductID:   product.ID,
			ProductName: product.Name,
			UnitPrice:   product.Price,
			Quantity:    item.Quantity,
		}
		if err := db.Create(&detail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, sale)

}

// Traer ventas
func SaleIndex(c *gin.Context) {
	var sales []models.Sale

	initializers.DB.
		Preload("Details").
		Preload("Details.Product").
		Preload("Payments").
		Find(&sales)

	c.JSON(http.StatusOK, gin.H{
		"sales": sales,
	})
}
