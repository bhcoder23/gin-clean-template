package entity

import "time"

// User -.
type User struct {
	ID           string    `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	Username     string    `example:"johndoe"                              json:"username"`
	Email        string    `example:"john@example.com"                     json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `example:"2026-01-01T00:00:00Z"                 json:"created_at"`
	UpdatedAt    time.Time `example:"2026-01-01T00:00:00Z"                 json:"updated_at"`
} // @name entity.User
