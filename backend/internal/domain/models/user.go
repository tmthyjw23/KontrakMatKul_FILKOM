package models

import (
	"time"
)

type User struct {
	ID            uint64    `json:"id"`
	StudentNumber string    `json:"student_number"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Hide password in JSON responses
	Role          string    `json:"role"` // 'STUDENT' or 'ADMIN'
	MaxSKS        int       `json:"max_sks"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserLoginRequest struct {
	StudentNumber string `json:"student_number" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID            uint64 `json:"id"`
	StudentNumber string `json:"student_number"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
}
