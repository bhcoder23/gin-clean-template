package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/usecase/user"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var errUseCaseInternal = errors.New("internal server error")

func newUserUseCase(t *testing.T) (*user.UseCase, *MockUserStore) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := NewMockUserStore(ctrl)
	jwtManager := jwt.New("test-secret", time.Hour)
	useCase := user.New(repo, jwtManager)

	return useCase, repo
}

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("register success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().Store(context.Background(), gomock.Any()).Return(nil)

		u, err := uc.Register(context.Background(), "testuser", "test@example.com", "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, u.ID)
		assert.Equal(t, "testuser", u.Username)
		assert.Equal(t, "test@example.com", u.Email)
	})

	t.Run("register duplicate", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().Store(context.Background(), gomock.Any()).Return(domain.ErrUserAlreadyExists)

		_, err := uc.Register(context.Background(), "testuser", "test@example.com", "password123")

		require.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("register trims username and email", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().Store(context.Background(), gomock.Any()).DoAndReturn(func(_ context.Context, u *domain.User) error {
			assert.Equal(t, "testuser", u.Username)
			assert.Equal(t, "test@example.com", u.Email)

			return nil
		})

		u, err := uc.Register(context.Background(), "  testuser  ", "  test@example.com  ", "password123")

		require.NoError(t, err)
		assert.Equal(t, "testuser", u.Username)
		assert.Equal(t, "test@example.com", u.Email)
	})
}

func TestRegisterValidatesCoreFieldsBeforeRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  error
	}{
		{name: "short username", username: "ab", email: "test@example.com", password: "password123", wantErr: domain.ErrInvalidUsername},
		{name: "invalid email", username: "testuser", email: "not-an-email", password: "password123", wantErr: domain.ErrInvalidEmail},
		{name: "display email", username: "testuser", email: "Test <test@example.com>", password: "password123", wantErr: domain.ErrInvalidEmail},
		{name: "short password", username: "testuser", email: "test@example.com", password: "12345", wantErr: domain.ErrPasswordTooShort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc, repo := newUserUseCase(t)
			repo.EXPECT().Store(gomock.Any(), gomock.Any()).Times(0)

			_, err := uc.Register(context.Background(), tt.username, tt.email, tt.password)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("login success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		storedUser := domain.User{
			ID: "user-id-123", Username: "testuser",
			Email: "test@example.com", PasswordHash: string(hash),
		}
		repo.EXPECT().GetByEmail(context.Background(), "test@example.com").Return(storedUser, nil)

		token, err := uc.Login(context.Background(), "test@example.com", "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("login wrong password", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		storedUser := domain.User{
			ID: "user-id-123", Username: "testuser",
			Email: "test@example.com", PasswordHash: string(hash),
		}
		repo.EXPECT().GetByEmail(context.Background(), "test@example.com").Return(storedUser, nil)

		token, err := uc.Login(context.Background(), "test@example.com", "wrongpassword")

		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("login user not found", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().GetByEmail(context.Background(), "notfound@example.com").Return(domain.User{}, domain.ErrUserNotFound)

		token, err := uc.Login(context.Background(), "notfound@example.com", "password123")

		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("login repo generic error", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().GetByEmail(context.Background(), "broken@example.com").Return(domain.User{}, errUseCaseInternal)

		token, err := uc.Login(context.Background(), "broken@example.com", "password123")

		require.ErrorIs(t, err, errUseCaseInternal)
		assert.NotErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	expectedUser := domain.User{
		ID:       "user-id-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	t.Run("get user success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().GetByID(context.Background(), "user-id-123").Return(expectedUser, nil)

		u, err := uc.GetUser(context.Background(), "user-id-123")

		require.NoError(t, err)
		assert.Equal(t, expectedUser, u)
	})

	t.Run("get user not found", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().GetByID(context.Background(), "missing-id").Return(domain.User{}, domain.ErrUserNotFound)

		_, err := uc.GetUser(context.Background(), "missing-id")

		require.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestGetUser_GenericError(t *testing.T) {
	t.Parallel()

	uc, repo := newUserUseCase(t)

	repo.EXPECT().GetByID(context.Background(), "user-id-123").Return(domain.User{}, errUseCaseInternal)

	_, err := uc.GetUser(context.Background(), "user-id-123")

	require.Error(t, err)
	require.ErrorIs(t, err, errUseCaseInternal)
}
