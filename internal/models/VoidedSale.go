// internal/models/VoidedSale.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type VoidedSale struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	OriginalSaleID uint      `json:"originalSaleId" gorm:"not null;index"`
	Total          uint      `json:"total" gorm:"not null"`
	SoldAt         time.Time `json:"soldAt" gorm:"not null"`

	VoidedByUserID uint    `json:"voidedByUserId" gorm:"not null"`
	VoidedByUser   User    `json:"-" gorm:"foreignKey:VoidedByUserID;references:ID"`
	Reason         *string `json:"reason" gorm:"size:500"`

	BusinessID uint     `json:"businessId" gorm:"not null"`
	Business   Business `json:"-" gorm:"foreignKey:BusinessID;references:ID"`

	Details  []VoidedSaleDetail  `json:"details,omitempty" gorm:"foreignKey:VoidedSaleID"`
	Payments []VoidedSalePayment `json:"payments,omitempty" gorm:"foreignKey:VoidedSaleID"`
}

type VoidedSaleDetail struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	Quantity    uint   `json:"quantity" gorm:"not null"`
	UnitPrice   uint   `json:"unitPrice" gorm:"not null"`
	ProductName string `json:"productName" gorm:"size:255;not null"`
	ProductID   uint   `json:"productId" gorm:"not null"`

	VoidedSaleID uint `json:"voidedSaleId" gorm:"not null"`

	BusinessID uint `json:"businessId" gorm:"not null"`
}

type VoidedSalePayment struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	Method string `json:"method" gorm:"size:255;not null"`
	Amount uint   `json:"amount" gorm:"not null"`

	VoidedSaleID uint `json:"voidedSaleId" gorm:"not null"`

	BusinessID uint `json:"businessId" gorm:"not null"`
}
