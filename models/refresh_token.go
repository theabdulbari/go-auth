package models

import (
	"time"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Token     string    `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
}