package auth_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

type LoginRequest struct {
	Email    string
	Password string
}

func NewLoginRequest(email, password string) LoginRequest {
	return LoginRequest{
		Email:    email,
		Password: password,
	}
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	Id           string
}

func NewLoginResponse(accessToken, refreshToken, id string) LoginResponse {
	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Id:           id,
	}
}

func (authUsecase *AuthUsecase) Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error) {
	user, err := authUsecase.userService.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return LoginResponse{}, apperrors.ErrLoginNotRegistered
		}

		authUsecase.logger.Error("failed to get user by email", logger.Error(err))
		return LoginResponse{}, apperrors.ErrInternalError
	}

	userPassword, err := authUsecase.authService.GetPasswordByUserId(ctx, user.ID)

	if errors.Is(err, apperrors.ErrNotFound) {
		return LoginResponse{}, apperrors.ErrWrongPassword
	} else if err != nil {
		authUsecase.logger.Error("failed to get user password", logger.Error(err))
		return LoginResponse{}, apperrors.ErrInternalError
	}

	if !authUsecase.passwordHasher.Compare(loginReq.Password, userPassword.Salt, userPassword.PasswordHash) {
		return LoginResponse{}, apperrors.ErrWrongPassword
	}

	accessToken, refreshToken, err := authUsecase.authService.GenerateJWT(ctx, user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to generate jwt: (%v): %w", err, apperrors.ErrInternalError)
	}

	return NewLoginResponse(accessToken, refreshToken, user.ID.String()), nil
}
