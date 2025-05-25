package user_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

func (userUsecase *UserUsecase) UpdateSettingTelegramUnlink(ctx context.Context) error {
	authClaims, err := checkWritePermissions(ctx)
	if err != nil {
		return err
	}

	currentAuthProvider, err := userUsecase.authService.GetAuthProviderByUserIdAndProviderName(ctx, authClaims.UserID, models.TelegramProviderName)
	if errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.ErrTelegramIsNotLinked
	} else if err != nil {
		userUsecase.logger.Error("failed to get auth provider", logger.Error(err))
		return apperrors.ErrInternalError
	}

	if err := userUsecase.authService.DeleteAuthProviderById(ctx, currentAuthProvider.ID); err != nil {
		userUsecase.logger.Error("failed to delete auth provider", logger.Error(err))
		return apperrors.ErrInternalError
	}

	return nil
}
