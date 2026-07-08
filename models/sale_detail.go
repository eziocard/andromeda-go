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

	Quantity    uint   `json:"quantity" gorm:"not null"`
	UnitPrice   uint   `json:"unitPrice" gorm:"not null"`
	ProductName string `json:"productName" gorm:"size:255;not null"`

	ProductID uint    `json:"productId" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`

	SaleID uint `json:"saleId" gorm:"not null"`
	Sale   Sale `json:"-" gorm:"foreignKey:SaleID;references:ID"`
}
