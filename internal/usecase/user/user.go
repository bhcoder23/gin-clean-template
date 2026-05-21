package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	appports "github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 6

// Usecase coordinates user application workflows.
type Usecase struct {
	repo appports.UserRepo
	jwt  *jwt.Manager
}

// New -.
func New(r appports.UserRepo, j *jwt.Manager) *Usecase {
	return &Usecase{
		repo: r,
		jwt:  j,
	}
}

// Register -.
func (uc *Usecase) Register(ctx context.Context, username, email, password string) (domain.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if len([]rune(username)) < 3 || len([]rune(username)) > 255 {
		return domain.User{}, domain.ErrInvalidUsername
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return domain.User{}, domain.ErrInvalidEmail
	}

	if len(password) < minPasswordLength {
		return domain.User{}, domain.ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("user.Usecase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	now := time.Now().UTC()

	user := domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.repo.Store(ctx, &user)
	if err != nil {
		return domain.User{}, fmt.Errorf("user.Usecase - Register - uc.repo.Store: %w", err)
	}

	return user, nil
}

// Login -.
func (uc *Usecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}

		return "", fmt.Errorf("user.Usecase - Login - uc.repo.GetByEmail: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("user.Usecase - Login - uc.jwt.GenerateToken: %w", err)
	}

	return token, nil
}

// GetUser -.
func (uc *Usecase) GetUser(ctx context.Context, userID string) (domain.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("user.Usecase - GetUser - uc.repo.GetByID: %w", err)
	}

	return user, nil
}
