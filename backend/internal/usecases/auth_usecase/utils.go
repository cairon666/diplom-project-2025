package auth_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

func genNeedEndRegistrationError(tempId uuid.UUID) error {
	errWithFields := apperrors.NeedEndRegistration().WithField("tempId", tempId.String())

	return errWithFields
}

type TelegramPartialUser struct {
	FirstName string
	LastName  string
	Id        int64
}

func (authUsecase *AuthUsecase) createPartialUserAndCreateTempId(ctx context.Context, req TelegramPartialUser) (uuid.UUID, error) {
	userId := uuid.New()
	createdAt := time.Now()
	user := models.NewUser(userId, "", req.FirstName, req.LastName, false, createdAt)

	if err := authUsecase.userService.CreateUser(ctx, user); err != nil {
		authUsecase.logger.Error("failed to create user", logger.Error(err))
		return uuid.Nil, apperrors.InternalError()
	}

	authProvider := models.NewAuthProvider(uuid.New(), userId, models.TelegramProviderName, req.Id, createdAt)
	if err := authUsecase.authService.CreateUserAuthProvider(ctx, authProvider); err != nil {
		authUsecase.logger.Error("failed to create auth provider", logger.Error(err))
		return uuid.Nil, apperrors.InternalError()
	}

	tempId, err := authUsecase.authService.GenerateTempIDToCompleteRegistration(ctx, userId)
	if err != nil {
		authUsecase.logger.Error("failed to generate temp id", logger.Error(err))
		return uuid.Nil, apperrors.InternalError()
	}

	return tempId, nil
}

func (authUsecase *AuthUsecase) generateTempIdAndReturnErrorNeedEndRegistration(ctx context.Context, userId uuid.UUID) error {
	tempId, err := authUsecase.authService.GenerateTempIDToCompleteRegistration(ctx, userId)
	if err != nil {
		authUsecase.logger.Error("failed to generate temp id", logger.Error(err))
		return apperrors.InternalError()
	}

	return genNeedEndRegistrationError(tempId)
}
