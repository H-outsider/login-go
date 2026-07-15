package model

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	Phone     string
	Status    int8
	CreatedAt time.Time
	UpdatedAt time.Time
}
