package models

import (
	"time"

	"gorm.io/gorm"
)

type Sale struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Total uint `json:"total" gorm:"not null"`

	Voided   bool       `json:"voided" gorm:"default:false;index"`
	VoidedAt *time.Time `json:"voidedAt"`

	BusinessID uint     `json:"businessId" gorm:"not null"`
	Business   Business `json:"-" gorm:"foreignKey:BusinessId;references:ID"`

	SellerID   uint   `json:"sellerId" gorm:"not null;default:0;index"`
	SellerName string `json:"sellerName" gorm:"size:255;not null;default:''"`
	Seller     User   `json:"seller,omitempty" gorm:"foreignKey:SellerID;references:ID"`

	Details  []SaleDetail  `json:"details,omitempty" gorm:"foreignKey:SaleID"`
	Payments []SalePayment `json:"payments,omitempty" gorm:"foreignKey:SaleID"`
}
