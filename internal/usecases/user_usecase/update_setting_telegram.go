package user_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type UpdateSettingTelegramRequest struct {
	Id        int64
	FirstName string
	LastName  string
	Username  string
	PhotoURL  string
	AuthDate  int
	Hash      string
}

func NewUpdateSettingTelegramRequest(id int64, firstName string, lastName string, username string, photoURL string, authDate int, hash string) UpdateSettingTelegramRequest {
	return UpdateSettingTelegramRequest{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		PhotoURL:  photoURL,
		AuthDate:  authDate,
		Hash:      hash,
	}
}

func (userUsecase *UserUsecase) UpdateSettingTelegram(ctx context.Context, dto UpdateSettingTelegramRequest) error {
	authClaims, err := checkWritePermissions(ctx)
	if err != nil {
		return err
	}

	if verifyOk := userUsecase.telegramService.Verify(models.TelegramAuthData{
		ID:        dto.Id,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Username:  dto.Username,
		PhotoURL:  dto.PhotoURL,
		AuthDate:  dto.AuthDate,
		Hash:      dto.Hash,
	}); !verifyOk {
		return apperrors.ErrInvalidTelegramHash
	}

	_, err = userUsecase.authService.GetAuthProviderByProviderUserIdAndProviderName(ctx, dto.Id, models.TelegramProviderName)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		userUsecase.logger.Error("failed to get auth provider", logger.Error(err))
		return err
	}

	if !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.ErrProviderAlreadyConnected
	}

	id := uuid.New()
	authProvider := models.NewAuthProvider(id, authClaims.UserID, models.TelegramProviderName, dto.Id, time.Now())
	if err := userUsecase.authService.CreateUserAuthProvider(ctx, authProvider); err != nil {
		userUsecase.logger.Error("failed to create auth provider", logger.Error(err))
		return err
	}

	return nil
}
