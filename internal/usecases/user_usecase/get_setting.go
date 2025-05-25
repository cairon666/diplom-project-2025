package user_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

type GetSettingResponse struct {
	Email       string
	HasPassword bool
	HasTelegram bool
	FirstName   string
	LastName    string
}

func (userUsecase *UserUsecase) GetSetting(ctx context.Context) (GetSettingResponse, error) {
	authClaims, err := checkReadPermissions(ctx)
	if err != nil {
		return GetSettingResponse{}, err
	}

	user, err := userUsecase.userService.GetUserById(ctx, authClaims.UserID)
	if err != nil {
		userUsecase.logger.Error("failed to get user by id", logger.Error(err))
		return GetSettingResponse{}, err
	}

	resDto := GetSettingResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	_, err = userUsecase.authService.GetPasswordByUserId(ctx, authClaims.UserID)
	if errors.Is(err, apperrors.ErrNotFound) {
		resDto.HasPassword = false
	} else if err != nil {
		userUsecase.logger.Error("failed to get user password by id", logger.Error(err))
		return GetSettingResponse{}, err
	} else {
		resDto.HasPassword = true
	}

	authProviders, err := userUsecase.authService.GetAuthProvidersByUserId(ctx, authClaims.UserID)
	if errors.Is(err, apperrors.ErrNotFound) {
		resDto.HasTelegram = false
	} else if err != nil {
		userUsecase.logger.Error("failed to get user auth providers by id", logger.Error(err))
		return GetSettingResponse{}, err
	}

	for _, authProvider := range authProviders {
		if authProvider.ProviderName == models.TelegramProviderName {
			resDto.HasTelegram = true
			break
		}
	}

	return resDto, nil
}
