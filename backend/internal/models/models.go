package models

import (
	"time"
)

type Item struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Email     string    `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username" validate:"required,username"`
	Password string `json:"password,omitempty" db:"password" validate:"required,strong_password"`
	Email    string `json:"email" db:"email" validate:"required,email"`
}

type CreateItemRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Timestamp map[string]interface{} `json:"timestamp"`
	Database  string                 `json:"database,omitempty"`
}