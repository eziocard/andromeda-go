package models

import "time"

type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	UserID    uint      `json:"userId" gorm:"index;not null"`
	TokenHash string    `json:"-" gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expiresAt"`
	Revoked   bool      `json:"-" gorm:"default:false"`
}
