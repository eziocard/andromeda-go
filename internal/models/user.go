package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	LastLogin *time.Time     `json:"lastLogin"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name     string `json:"name" gorm:"size:255;not null"`
	LastName string `json:"lastName" gorm:"size:255;not null"`
	Contact  string `json:"contact" gorm:"size:255"`
	Email    string `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Password string `json:"-" gorm:"size:255;not null"`
	IsActive bool   `json:"isActive" gorm:"default:true"`

	FailedLoginAttempts int        `json:"-" gorm:"default:0"`
	LockedUntil         *time.Time `json:"-"`

	BusinessID *uint     `json:"businessId" gorm:"index"`
	Business   *Business `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
	//Business   *Business `json:"-" gorm:"foreignKey:businessId;references:ID"`

	RoleID uint `json:"roleId" gorm:"not null;index"`
	Role   Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
}
