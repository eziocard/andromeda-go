package models

import (
	"time"

	"gorm.io/gorm"
)

// aqui solo almaceno el total de la venta
type Sale struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Total uint `json:"total" gorm:"not null"`

	Details  []SaleDetail  `json:"details,omitempty" gorm:"foreignKey:SaleID"`
	Payments []SalePayment `json:"payments,omitempty" gorm:"foreignKey:SaleID"`
}
