package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenBucket struct {
	gorm.Model
	Identifier string    `gorm:"index; not null"`
	Tokens     float64   `gorm:"not null"`
	LastRefill time.Time `gorm:"not null"`
}
