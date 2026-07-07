package models

import (
	"time"

	"gorm.io/gorm"
)

type SaleDetail struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ProductName string `json:"productName" gorm:"size:255;not null"`
	UnitPrice   uint   `json:"unitPrice" gorm:"not null"`
	Quantity    uint   `json:"quantity" gorm:"not null"`

	// Foreign Key Sale
	SaleID uint `json:"saleId"`
	Sale   Sale `json:"-"`

	// Foreign Key Product
	ProductID uint    `json:"productId"`
	Product   Product `json:"product"`
}
