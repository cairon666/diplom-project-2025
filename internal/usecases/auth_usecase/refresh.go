package auth_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

type RefreshTokenRequest struct {
	RefreshToken string
}

func NewRefreshTokenRequest(refreshToken string) RefreshTokenRequest {
	return RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

func NewRefreshTokenResponse(accessToken, refreshToken string) RefreshTokenResponse {
	return RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (authUsecase *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenReq RefreshTokenRequest) (RefreshTokenResponse, error) {
	claims, err := authUsecase.jwtService.ParseToken(refreshTokenReq.RefreshToken)
	if err != nil {
		authUsecase.logger.Error("failed to parse token", logger.Error(err))
		return RefreshTokenResponse{}, apperrors.ErrInvalidToken
	}

	user, err := authUsecase.userService.GetUserById(ctx, claims.UserID)
	if err != nil {
		authUsecase.logger.Error("failed to get user by id", logger.Error(err))
		return RefreshTokenResponse{}, apperrors.ErrInternalError
	}

	accessToken, refreshToken, err := authUsecase.authService.GenerateJWT(ctx, user)
	if err != nil {
		authUsecase.logger.Error("failed to generate jwt", logger.Error(err))
		return RefreshTokenResponse{}, apperrors.ErrInternalError
	}

	return NewRefreshTokenResponse(accessToken, refreshToken), nil
}
