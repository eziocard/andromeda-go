package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name string `json:"name" gorm:"size:255;not null"`
}
