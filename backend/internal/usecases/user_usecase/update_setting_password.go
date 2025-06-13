package user_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

type UpdateSettingPasswordRequest struct {
	Password string
}

func NewUpdateSettingPasswordRequest(password string) UpdateSettingPasswordRequest {
	return UpdateSettingPasswordRequest{
		Password: password,
	}
}

func (userUsecase *UserUsecase) UpdateSettingPassword(ctx context.Context, dto UpdateSettingPasswordRequest) error {
	authClaims, err := checkWritePermissions(ctx)
	if err != nil {
		return err
	}

	user, err := userUsecase.userService.GetUserById(ctx, authClaims.UserID)
	if err != nil {
		userUsecase.logger.Error("failed to get user by id", logger.Error(err))

		return err
	}

	passwordHash, salt, err := userUsecase.passwordHasher.Hash(dto.Password)
	if err != nil {
		userUsecase.logger.Error("failed to hash password", logger.Error(err))

		return apperrors.InternalError()
	}

	newUserPassword := models.NewUserPassword(user.ID, salt, passwordHash)

	_, err = userUsecase.authService.GetPasswordByUserId(ctx, authClaims.UserID)
	if errors.Is(err, apperrors.NotFound()) {
		err = userUsecase.authService.CreateUserPassword(ctx, newUserPassword)
		if err != nil {
			userUsecase.logger.Error("failed to create user password", logger.Error(err))

			return err
		}

		return nil
	} else if err != nil {
		userUsecase.logger.Error("failed to get user password by id", logger.Error(err))

		return err
	}

	if err := userUsecase.authService.UpdateUserPassword(ctx, newUserPassword); err != nil {
		userUsecase.logger.Error("failed to update user password", logger.Error(err))

		return err
	}

	return nil
}
