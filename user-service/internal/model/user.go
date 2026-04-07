package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           int64     `json:"id"`
	email        string    `json:"email"`
	passwordHash string    `json:"-"`
	createdAt    time.Time `json:"created_at"`
	updatedAt    time.Time `json:"updated_at"`
}

func NewUser(email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt password: %w", err)
	}

	now := time.Now()

	return &User{
		email:        email,
		passwordHash: string(hashed),
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (u *User) Email() string {
	return u.email
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}
