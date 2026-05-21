package response

import (
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
)

type UserResp struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserResp(user *domain.User) UserResp {
	return UserResp{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
