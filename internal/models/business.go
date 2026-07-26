package models

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name     string `json:"name" gorm:"size:255;not null"`
	Address  string `json:"address" gorm:"size:255;not null"`
	City     string `json:"city" gorm:"size:255;not null"`
	Contact  string `json:"contact" gorm:"size:255"`
	Email    string `json:"email" gorm:"size:255;not null;uniqueIndex"`
	IsActive bool   `json:"isActive" gorm:"default:true"`
}
