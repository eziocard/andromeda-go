package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Barcode    string  `json:"barcode" gorm:"size:100;not null;uniqueIndex:idx_business_barcode"`
	Name       string  `json:"name" gorm:"size:255;not null"`
	Variety    *string `json:"variety" gorm:"size:255"`
	Image      *string `json:"image" gorm:"size:255"`
	Price      uint    `json:"price" gorm:"not null"`
	Stock      uint    `json:"stock" gorm:"default:0"`
	IsWeighted bool    `json:"isWeighted" gorm:"default:false"`

	BusinessID uint     `json:"businessId" gorm:"not null;uniqueIndex:idx_business_barcode"`
	Business   Business `json:"-" gorm:"foreignKey:BusinessID;references:ID"`
}
