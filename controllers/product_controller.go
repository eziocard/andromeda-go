package controllers

import (
	"strconv"
	"strings"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/models"
	"github.com/eziocard/andromeda-go/utils"
	"github.com/gin-gonic/gin"
)

func ProductCreate(c *gin.Context) {
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
		Barcode: barcode,
		Name:    name,
		Variety: variety,
		Price:   uint(price),
		Stock:   uint(stock),
	}

	// Imagen es opcional, igual que blank=True en Django
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
	var products []models.Product
	initializers.DB.Find(&products)
	c.JSON(200, gin.H{"products": products})
}

func ProductShow(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	initializers.DB.First(&product, id)
	c.JSON(200, gin.H{"product": product})
}

// Equivalente al @action by_barcode de Django
func ProductShowByBarcode(c *gin.Context) {
	barcode := strings.TrimSpace(c.Param("barcode"))
	var product models.Product
	result := initializers.DB.Where("barcode = ?", barcode).First(&product)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(200, gin.H{"product": product})
}

func ProductsUpdate(c *gin.Context) {
	barcode := strings.TrimSpace(c.Param("barcode"))

	var product models.Product
	if err := initializers.DB.Where("barcode = ?", barcode).First(&product).Error; err != nil {
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
	barcode := strings.TrimSpace(c.Param("barcode"))

	var product models.Product
	if err := initializers.DB.Where("barcode = ?", barcode).First(&product).Error; err == nil && product.Image != nil {
		_ = utils.DeleteProductImage(*product.Image)
	}

	initializers.DB.Where("barcode = ?", barcode).Delete(&models.Product{})
	c.JSON(200, gin.H{"message": "Product removed Successfully"})
}
