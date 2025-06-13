package auth_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type TelegramLoginRequest struct {
	Id        int64
	FirstName string
	LastName  string
	Username  string
	PhotoURL  string
	AuthDate  int
	Hash      string
}

func NewTelegramLoginRequest(id int64, firstName, lastName, username, photoURL string, authDate int, hash string) TelegramLoginRequest {
	return TelegramLoginRequest{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		PhotoURL:  photoURL,
		AuthDate:  authDate,
		Hash:      hash,
	}
}

type TelegramLoginResponse struct {
	AccessToken  string
	RefreshToken string
	UserId       uuid.UUID
}

func NewTelegramLoginResponse(accessToken, refreshToken string, userId uuid.UUID) TelegramLoginResponse {
	return TelegramLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userId,
	}
}

func (authUsecase *AuthUsecase) TelegramLogin(ctx context.Context, telegramLoginReq TelegramLoginRequest) (TelegramLoginResponse, error) {
	authProvider, err := authUsecase.authService.GetAuthProviderByProviderUserIdAndProviderName(ctx, telegramLoginReq.Id, models.TelegramProviderName)
	if errors.Is(err, apperrors.NotFound()) {
		tempId, err := authUsecase.createPartialUserAndCreateTempId(ctx, TelegramPartialUser{
			FirstName: telegramLoginReq.FirstName,
			LastName:  telegramLoginReq.LastName,
			Id:        telegramLoginReq.Id,
		})
		if err != nil {
			authUsecase.logger.Error("failed to create partial user and create temp id", logger.Error(err))

			return TelegramLoginResponse{}, err
		}

		return TelegramLoginResponse{}, genNeedEndRegistrationError(tempId)
	} else if err != nil {
		authUsecase.logger.Error("failed to get auth provider", logger.Error(err))

		return TelegramLoginResponse{}, apperrors.InternalError()
	}

	user, err := authUsecase.userService.GetUserById(ctx, authProvider.UserId)
	if err != nil {
		authUsecase.logger.Error("failed to get user by id", logger.Error(err))

		return TelegramLoginResponse{}, apperrors.InternalError()
	}

	if !user.IsRegistrationComplete {
		return TelegramLoginResponse{}, authUsecase.generateTempIdAndReturnErrorNeedEndRegistration(ctx, user.ID)
	}

	if verifyOk := authUsecase.telegramService.Verify(models.TelegramAuthData{
		ID:        telegramLoginReq.Id,
		FirstName: telegramLoginReq.FirstName,
		LastName:  telegramLoginReq.LastName,
		Username:  telegramLoginReq.Username,
		PhotoURL:  telegramLoginReq.PhotoURL,
		AuthDate:  telegramLoginReq.AuthDate,
		Hash:      telegramLoginReq.Hash,
	}); !verifyOk {
		return TelegramLoginResponse{}, apperrors.InvalidTelegramHash()
	}

	accessToken, refreshToken, err := authUsecase.authService.GenerateJWT(ctx, user)
	if err != nil {
		authUsecase.logger.Error("failed to generate jwt", logger.Error(err))

		return TelegramLoginResponse{}, apperrors.InternalError()
	}

	return NewTelegramLoginResponse(accessToken, refreshToken, authProvider.UserId), nil
}
