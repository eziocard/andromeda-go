package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/dto"
	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaleCreate(c *gin.Context) {
	user, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	var body dto.SaleCreateInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := initializers.DB
	var sale models.Sale

	err := db.Transaction(func(tx *gorm.DB) error {
		var total uint

		products := make(map[uint]models.Product)
		unitPrices := make(map[uint]uint)

		for _, item := range body.Details {
			var product models.Product
			if err := tx.Where("business_id = ?", businessID).First(&product, item.Product).Error; err != nil {
				return fmt.Errorf("producto no encontrado o no pertenece a tu negocio")
			}

			var unitPrice uint
			if product.IsWeighted {
				if item.UnitPrice == nil || *item.UnitPrice == 0 {
					return fmt.Errorf("debes ingresar un monto válido para %s", product.Name)
				}
				unitPrice = *item.UnitPrice
			} else {
				if product.Stock < item.Quantity {
					return fmt.Errorf("stock insuficiente para %s (disponible: %d, solicitado: %d)", product.Name, product.Stock, item.Quantity)
				}
				unitPrice = product.Price
			}

			total += unitPrice * item.Quantity
			products[product.ID] = product
			unitPrices[product.ID] = unitPrice
		}

		sale = models.Sale{
			Total:      total,
			BusinessID: businessID,
			SellerID:   user.ID,
			SellerName: user.Name + " " + user.LastName,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}

		for _, p := range body.Payments {
			payment := models.SalePayment{
				SaleID:     sale.ID,
				Method:     p.Method,
				Amount:     p.Amount,
				BusinessID: businessID,
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		for _, item := range body.Details {
			product := products[uint(item.Product)]
			unitPrice := unitPrices[product.ID]

			detail := models.SaleDetail{
				SaleID:      sale.ID,
				ProductID:   product.ID,
				ProductName: product.Name,
				UnitPrice:   unitPrice,
				Quantity:    item.Quantity,
				BusinessID:  businessID,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			if !product.IsWeighted {
				if err := tx.Model(&models.Product{}).
					Where("id = ?", product.ID).
					UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sale)
}
func SaleIndex(c *gin.Context) {
	user, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}
	_ = user

	year := c.Query("year")
	month := c.Query("month")
	day := c.Query("day")
	includeVoided := c.Query("includeVoided") == "true"

	db := initializers.DB.
		Preload("Details").
		Preload("Details.Product").
		Preload("Payments").
		Where("business_id = ?", businessID)

	if !includeVoided {
		db = db.Where("voided = ?", false)
	}

	if year != "" {
		db = db.Where("EXTRACT(YEAR FROM created_at) = ?", year)
	}
	if month != "" {
		db = db.Where("EXTRACT(MONTH FROM created_at) = ?", month)
	}
	if day != "" {
		db = db.Where("EXTRACT(DAY FROM created_at) = ?", day)
	}

	var sales []models.Sale
	if err := db.Order("created_at DESC").Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sales": sales})
}

func SaleVoid(c *gin.Context) {
	user, businessID, ok := getAuthUserBusiness(c)
	if !ok {
		return
	}

	saleIDParam := c.Param("id")
	saleID, err := strconv.ParseUint(saleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de venta inválido"})
		return
	}

	var body dto.SaleVoidInput
	_ = c.ShouldBindJSON(&body)

	db := initializers.DB
	var voidedSale models.VoidedSale

	err = db.Transaction(func(tx *gorm.DB) error {
		var sale models.Sale

		if err := tx.
			Preload("Details").
			Preload("Payments").
			Where("business_id = ?", businessID).
			First(&sale, saleID).Error; err != nil {
			return fmt.Errorf("venta no encontrada o no pertenece a tu negocio")
		}

		if sale.Voided {
			return fmt.Errorf("esta venta ya fue anulada")
		}

		now := time.Now()

		var reason *string
		if body.Reason != "" {
			reason = &body.Reason
		}

		voidedSale = models.VoidedSale{
			OriginalSaleID: sale.ID,
			Total:          sale.Total,
			SoldAt:         sale.CreatedAt,
			VoidedByUserID: user.ID,
			Reason:         reason,
			BusinessID:     businessID,
		}
		if err := tx.Create(&voidedSale).Error; err != nil {
			return err
		}

		for _, d := range sale.Details {
			voidedDetail := models.VoidedSaleDetail{
				Quantity:     d.Quantity,
				UnitPrice:    d.UnitPrice,
				ProductName:  d.ProductName,
				ProductID:    d.ProductID,
				VoidedSaleID: voidedSale.ID,
				BusinessID:   businessID,
			}
			if err := tx.Create(&voidedDetail).Error; err != nil {
				return err
			}

			var product models.Product
			if err := tx.Where("business_id = ?", businessID).First(&product, d.ProductID).Error; err == nil {
				if !product.IsWeighted {
					// Reincorporamos el stock del producto solo si aplica
					if err := tx.Model(&models.Product{}).
						Where("id = ? AND business_id = ?", d.ProductID, businessID).
						UpdateColumn("stock", gorm.Expr("stock + ?", d.Quantity)).Error; err != nil {
						return err
					}
				}
			}
			// Si el producto ya no existe (fue borrado), simplemente no reincorporamos stock.
		}
		for _, p := range sale.Payments {
			voidedPayment := models.VoidedSalePayment{
				Method:       p.Method,
				Amount:       p.Amount,
				VoidedSaleID: voidedSale.ID,
				BusinessID:   businessID,
			}
			if err := tx.Create(&voidedPayment).Error; err != nil {
				return err
			}
		}

		// 3. Marcamos la venta original como anulada (no la borramos)
		if err := tx.Model(&sale).Updates(map[string]interface{}{
			"voided":    true,
			"voided_at": now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "venta anulada correctamente",
		"voidedSale": voidedSale,
	})
}
