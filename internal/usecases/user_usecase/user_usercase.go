package user_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserById(ctx context.Context, ID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type AuthService interface {
	GetPasswordByUserId(ctx context.Context, userId uuid.UUID) (models.UserPassword, error)
	GetAuthProvidersByUserId(ctx context.Context, userId uuid.UUID) ([]models.AuthProvider, error)
	CreateUserPassword(ctx context.Context, userPassword models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userPassword models.UserPassword) error
	GetAuthProviderByProviderUserIdAndProviderName(ctx context.Context, providerUsrId int64, providerName string) (models.AuthProvider, error)
	CreateUserAuthProvider(ctx context.Context, authProvider models.AuthProvider) error
	DeleteAuthProviderById(ctx context.Context, id uuid.UUID) error
	GetAuthProviderByUserIdAndProviderName(ctx context.Context, userId uuid.UUID, providerName string) (models.AuthProvider, error)
}

type PasswordHasher interface {
	Hash(password string) (string, string, error)
	Compare(password, salt, hash string) bool
}

type TelegramService interface {
	Verify(data models.TelegramAuthData) bool
}

type UserUsecase struct {
	logger          *logger.Logger
	passwordHasher  PasswordHasher
	userService     UserService
	telegramService TelegramService
	authService     AuthService
}

func NewUserUsecase(
	userService UserService,
	logger *logger.Logger,
	passwordHasher PasswordHasher,
	telegramService TelegramService,
	authService AuthService,
) *UserUsecase {
	return &UserUsecase{
		userService:     userService,
		logger:          logger,
		passwordHasher:  passwordHasher,
		telegramService: telegramService,
		authService:     authService,
	}
}

func checkReadPermissions(ctx context.Context) (*indentity.AuthClaims, error) {
	claims, ok := indentity.GetAuthClaims(ctx)
	if !ok {
		return nil, apperrors.ErrForbidden
	}

	if claims.HasPermission(permission.ReadAllUsers, permission.ReadOwnProfile) {
		return claims, nil
	}

	return nil, apperrors.ErrForbidden
}

func checkWritePermissions(ctx context.Context) (*indentity.AuthClaims, error) {
	claims, ok := indentity.GetAuthClaims(ctx)
	if !ok {
		return nil, apperrors.ErrForbidden
	}

	if claims.HasPermission(permission.UpdateOwnProfile) {
		return claims, nil
	}

	return nil, apperrors.ErrForbidden
}
