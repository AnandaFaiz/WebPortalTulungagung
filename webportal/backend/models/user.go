package models

import "time"

type User struct {
	ID           int        `json:"id"`
	Nama         string     `json:"nama" binding:"required"`
	Email        string     `json:"email" binding:"required,email"`
	Password     string     `json:"password,omitempty"`
	Role         string     `json:"role" binding:"required"`
	OPD          string     `json:"opd"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nama  string `json:"nama" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}