package auth_service

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type TempIdRepo interface {
	Get(ctx context.Context, tempId string) (string, error)
	Set(ctx context.Context, tempId string, data string, exp time.Duration) error
	Delete(ctx context.Context, tempId string) error
}

type JWTService interface {
	CreateAccessToken(authClaims indentity.AuthClaims) (string, error)
	CreateRefreshToken(authClaims indentity.AuthClaims) (string, error)
}

type UserPasswordsRepo interface {
	GetUserPasswordByUserId(ctx context.Context, userId uuid.UUID) (models.UserPassword, error)
	CreateUserPassword(ctx context.Context, userPassword models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userPassword models.UserPassword) error
}

type AuthProviderRepo interface {
	GetAuthProviderByProviderUserIdAndProviderName(ctx context.Context, providerUserID int64, providerName string) (models.AuthProvider, error)
	CreateAuthProvider(ctx context.Context, authProvider models.AuthProvider) error
	GetAuthProvidersByUserId(ctx context.Context, userId uuid.UUID) ([]models.AuthProvider, error)
	DeleteAuthProviderById(ctx context.Context, id uuid.UUID) error
	GetAuthProviderByUserIdAndProviderName(ctx context.Context, userId uuid.UUID, providerName string) (models.AuthProvider, error)
}

type RoleRepos interface {
	GetRolesByUserID(ctx context.Context, userId uuid.UUID) ([]models.Role, error)
	GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Permission, error)
}

type AuthService struct {
	jwtService        JWTService
	tempIdRepo        TempIdRepo
	userPasswordsRepo UserPasswordsRepo
	authProviderRepo  AuthProviderRepo
	roleRepos         RoleRepos
}

func NewAuthService(
	jwtService JWTService,
	tempIdRepo TempIdRepo,
	userPasswordsRepo UserPasswordsRepo,
	authProviderRepo AuthProviderRepo,
	roleRepos RoleRepos,
) *AuthService {
	return &AuthService{
		jwtService:        jwtService,
		tempIdRepo:        tempIdRepo,
		userPasswordsRepo: userPasswordsRepo,
		authProviderRepo:  authProviderRepo,
		roleRepos:         roleRepos,
	}
}

func (userService *AuthService) GenerateJWT(ctx context.Context, user *models.User) (string, string, error) {
	roles, err := userService.roleRepos.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	permissions, err := userService.roleRepos.GetPermissionsByUserID(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	authClaims := indentity.AuthClaims{
		UserID:      user.ID,
		Roles:       roles,
		Permissions: permissions,
	}

	accessToken, err := userService.jwtService.CreateAccessToken(authClaims)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := userService.jwtService.CreateRefreshToken(authClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (userService *AuthService) GenerateTempIDToCompleteRegistration(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	tempID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	err = userService.tempIdRepo.Set(ctx, tempID.String(), userId.String(), time.Hour)
	if err != nil {
		return uuid.Nil, err
	}

	return tempID, nil
}

func (userService *AuthService) GetUserIdByTempID(ctx context.Context, tempID uuid.UUID) (uuid.UUID, error) {
	id, err := userService.tempIdRepo.Get(ctx, tempID.String())
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}

func (userService *AuthService) DeleteTempID(ctx context.Context, tempID uuid.UUID) error {
	return userService.tempIdRepo.Delete(ctx, tempID.String())
}

func (userService *AuthService) GetPasswordByUserId(ctx context.Context, userId uuid.UUID) (models.UserPassword, error) {
	return userService.userPasswordsRepo.GetUserPasswordByUserId(ctx, userId)
}

func (userService *AuthService) CreateUserPassword(ctx context.Context, userPassword models.UserPassword) error {
	return userService.userPasswordsRepo.CreateUserPassword(ctx, userPassword)
}

func (userService *AuthService) UpdateUserPassword(ctx context.Context, userPassword models.UserPassword) error {
	return userService.userPasswordsRepo.UpdateUserPassword(ctx, userPassword)
}

func (userService *AuthService) GetAuthProviderByProviderUserIdAndProviderName(ctx context.Context, providerUsrId int64, providerName string) (models.AuthProvider, error) {
	return userService.authProviderRepo.GetAuthProviderByProviderUserIdAndProviderName(ctx, providerUsrId, providerName)
}

func (userService *AuthService) CreateUserAuthProvider(ctx context.Context, authProvider models.AuthProvider) error {
	return userService.authProviderRepo.CreateAuthProvider(ctx, authProvider)
}

func (userService *AuthService) GetAuthProvidersByUserId(ctx context.Context, userId uuid.UUID) ([]models.AuthProvider, error) {
	return userService.authProviderRepo.GetAuthProvidersByUserId(ctx, userId)
}

func (userService *AuthService) GetAuthProviderByUserIdAndProviderName(ctx context.Context, userId uuid.UUID, providerName string) (models.AuthProvider, error) {
	return userService.authProviderRepo.GetAuthProviderByUserIdAndProviderName(ctx, userId, providerName)
}

func (userService *AuthService) DeleteAuthProviderById(ctx context.Context, id uuid.UUID) error {
	return userService.authProviderRepo.DeleteAuthProviderById(ctx, id)
}
