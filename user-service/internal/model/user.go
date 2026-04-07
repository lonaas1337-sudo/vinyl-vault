package model

import "time"

type User struct {
	id           int64     `json:"id"`
	email        string    `json:"email"`
	passwordHash string    `json:"-"`
	createdAt    time.Time `json:"created_at"`
	updatedAt    time.Time `json:"updated_at"`
}
