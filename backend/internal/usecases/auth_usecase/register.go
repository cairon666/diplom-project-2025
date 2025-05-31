package auth_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email      string
	Password   string
	FirstName  string
	SecondName string
}

func NewRegisterRequest(email, password, firstName, lastName string) RegisterRequest {
	return RegisterRequest{
		Email:      email,
		Password:   password,
		FirstName:  firstName,
		SecondName: lastName,
	}
}

func (authUsecase *AuthUsecase) Register(ctx context.Context, registerReq RegisterRequest) error {
	createdAt := time.Now()
	userId := uuid.New()
	passwordHash, salt, err := authUsecase.passwordHasher.Hash(registerReq.Password)
	if err != nil {
		authUsecase.logger.Error("failed to hash password", logger.Error(err))
		return apperrors.InternalError()
	}

	user := models.NewUser(userId, registerReq.Email, registerReq.FirstName, registerReq.SecondName, true, createdAt)
	if err := authUsecase.userService.CreateUser(ctx, user); err != nil {
		if errors.Is(err, apperrors.AlreadyExists()) {
			return apperrors.EmailAlreadyExists()
		}

		authUsecase.logger.Error("failed to create user", logger.Error(err))
		return apperrors.InternalError()
	}

	userPassword := models.NewUserPassword(userId, salt, passwordHash)
	if err := authUsecase.authService.CreateUserPassword(ctx, userPassword); err != nil {
		authUsecase.logger.Error("failed to create user password", logger.Error(err))
		return apperrors.InternalError()
	}

	if err := authUsecase.roleService.AssignRoleToUser(ctx, userId, models.RoleUser); err != nil {
		authUsecase.logger.Error("failed to assign role to user", logger.Error(err))
		return apperrors.InternalError()
	}

	return nil
}
