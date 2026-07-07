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
	Total     uint           `json:"total" gorm:"not null"`
	Details   []SaleDetail   `json:"details"`
	Payments  []SalePayment  `json:"payments"`
}
