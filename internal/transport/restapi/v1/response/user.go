package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

// UserResp documents the REST user response shape.
type UserResp struct {
	ID        string    `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	Username  string    `example:"johndoe"                              json:"username"`
	Email     string    `example:"john@example.com"                     json:"email"`
	CreatedAt time.Time `example:"2026-01-01T00:00:00Z"                 json:"created_at"`
	UpdatedAt time.Time `example:"2026-01-01T00:00:00Z"                 json:"updated_at"`
} // @name v1.UserResp

func NewUserResp(user *domain.User) UserResp {
	return UserResp{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
