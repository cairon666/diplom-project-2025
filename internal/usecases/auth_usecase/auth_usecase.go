package auth_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserById(ctx context.Context, ID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
}

type AuthService interface {
	GenerateJWT(ctx context.Context, user *models.User) (string, string, error)
	GenerateTempIDToCompleteRegistration(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	GetUserIdByTempID(ctx context.Context, tempID uuid.UUID) (uuid.UUID, error)
	DeleteTempID(ctx context.Context, tempID uuid.UUID) error
	GetPasswordByUserId(ctx context.Context, userId uuid.UUID) (models.UserPassword, error)
	CreateUserPassword(ctx context.Context, userPassword models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userPassword models.UserPassword) error
	GetAuthProviderByProviderUserIdAndProviderName(ctx context.Context, providerUsrId int64, providerName string) (models.AuthProvider, error)
	CreateUserAuthProvider(ctx context.Context, authProvider models.AuthProvider) error
}

type PasswordHasher interface {
	Hash(password string) (string, string, error)
	Compare(password, salt, hash string) bool
}

type TelegramService interface {
	Verify(data models.TelegramAuthData) bool
}

type RolesService interface {
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error
}

type AuthUsecase struct {
	userService     UserService
	jwtService      *indentity.JWTService
	passwordHasher  PasswordHasher
	telegramService TelegramService
	authService     AuthService
	roleService     RolesService
	config          *config.Config
	logger          *logger.Logger
}

func NewAuthUsecase(
	userService UserService,
	jwtService *indentity.JWTService,
	logger *logger.Logger,
	config *config.Config,
	passwordHasher PasswordHasher,
	telegramService TelegramService,
	authService AuthService,
	roleService RolesService,
) *AuthUsecase {
	return &AuthUsecase{
		userService:     userService,
		logger:          logger,
		jwtService:      jwtService,
		config:          config,
		passwordHasher:  passwordHasher,
		telegramService: telegramService,
		authService:     authService,
		roleService:     roleService,
	}
}
