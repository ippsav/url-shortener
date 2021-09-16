package models

import (
	"time"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewUser() *User {
	return &User{}
}
