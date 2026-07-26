package controllers

import (
	"strconv"
	"strings"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/eziocard/andromeda-go/internal/utils"
	"github.com/gin-gonic/gin"
)

func ProductCreate(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	barcode := strings.TrimSpace(c.PostForm("barcode"))
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	stockStr := c.PostForm("stock")

	price, _ := strconv.ParseUint(priceStr, 10, 32)
	stock, _ := strconv.ParseUint(stockStr, 10, 32)

	var variety *string
	if v := c.PostForm("variety"); v != "" {
		variety = &v
	}

	product := models.Product{
		Barcode:    barcode,
		Name:       name,
		Variety:    variety,
		Price:      uint(price),
		Stock:      uint(stock),
		BusinessID: businessID,
	}

	fileHeader, err := c.FormFile("image")
	if err == nil {
		imgPath, err := utils.SaveProductImage(fileHeader)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		product.Image = &imgPath
	}

	result := initializers.DB.Create(&product)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{"product": product})
}

func ProductIndex(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	var products []models.Product
	initializers.DB.Where("business_id = ?", businessID).Find(&products)
	c.JSON(200, gin.H{"products": products})
}

func ProductShow(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	id := c.Param("id")
	var product models.Product
	result := initializers.DB.Where("business_id = ?", businessID).First(&product, id)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(200, gin.H{"product": product})
}

func ProductShowByBarcode(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	barcode := strings.TrimSpace(c.Param("barcode"))
	var product models.Product
	result := initializers.DB.Where("barcode = ? AND business_id = ?", barcode, businessID).First(&product)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(200, gin.H{"product": product})
}

func ProductsUpdate(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	barcode := strings.TrimSpace(c.Param("barcode"))

	var product models.Product
	if err := initializers.DB.Where("barcode = ? AND business_id = ?", barcode, businessID).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Producto no encontrado"})
		return
	}

	if newBarcode := strings.TrimSpace(c.PostForm("barcode")); newBarcode != "" {
		product.Barcode = newBarcode
	}
	if name := c.PostForm("name"); name != "" {
		product.Name = name
	}
	if v := c.PostForm("variety"); v != "" {
		product.Variety = &v
	}
	if priceStr := c.PostForm("price"); priceStr != "" {
		price, _ := strconv.ParseUint(priceStr, 10, 32)
		product.Price = uint(price)
	}
	if stockStr := c.PostForm("stock"); stockStr != "" {
		stock, _ := strconv.ParseUint(stockStr, 10, 32)
		product.Stock = uint(stock)
	}

	fileHeader, err := c.FormFile("image")
	if err == nil {
		if product.Image != nil {
			_ = utils.DeleteProductImage(*product.Image)
		}
		imgPath, err := utils.SaveProductImage(fileHeader)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		product.Image = &imgPath
	}

	initializers.DB.Save(&product)
	c.JSON(200, gin.H{"product": product})
}

func ProductsDelete(c *gin.Context) {
	_, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	id := c.Param("id")

	result := initializers.DB.Where("business_id = ?", businessID).Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(200, gin.H{"message": "Producto eliminado correctamente"})
}
