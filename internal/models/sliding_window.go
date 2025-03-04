package models

import (
	"time"

	"gorm.io/gorm"
)

type RateLimit struct {
	gorm.Model
	Identifier string `gorm:"index;not null`
	Count int `gorm:"not null"`
	LastRequest time.Time `gorm:"not null"`
	ResetAt time.Time `gorm:"not null"`
}