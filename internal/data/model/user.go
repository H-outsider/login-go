package model

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"size:32;not null;uniqueIndex"`
	Password  string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255"`
	Phone     string `gorm:"size:32"`
	Status    int8   `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
