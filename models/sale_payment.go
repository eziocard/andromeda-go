package models

import (
	"time"

	"gorm.io/gorm"
)

type SalePayment struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Method string `json:"method" gorm:"size:10;not null"`
	Amount uint   `json:"amount" gorm:"not null"`

	// Foreign Key Sale
	SaleID uint `json:"saleId"`
	Sale   Sale `json:"-"`
}
