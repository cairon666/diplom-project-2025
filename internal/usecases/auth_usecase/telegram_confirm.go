package auth_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type CompleteRegisterRequest struct {
	TempId    uuid.UUID
	FirstName string
	LastName  string
	Email     string
}

func NewCompleteRegisterRequest(tempId uuid.UUID, email, firstName, lastName string) CompleteRegisterRequest {
	return CompleteRegisterRequest{
		TempId:    tempId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}

type CompleteRegisterResponse struct {
	AccessToken  string
	RefreshToken string
	UserId       uuid.UUID
}

func NewCompleteRegisterResponse(accessToken, refreshToken string, userId uuid.UUID) CompleteRegisterResponse {
	return CompleteRegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userId,
	}
}

func (authUsecase *AuthUsecase) TelegramConfirm(ctx context.Context, reqDTO CompleteRegisterRequest) (CompleteRegisterResponse, error) {
	userId, err := authUsecase.authService.GetUserIdByTempID(ctx, reqDTO.TempId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return CompleteRegisterResponse{}, apperrors.ErrTempIdNotFound
		}

		authUsecase.logger.Error("failed to get user id by temp id", logger.Error(err))
		return CompleteRegisterResponse{}, apperrors.ErrInternalError
	}

	user, err := authUsecase.userService.GetUserById(ctx, userId)
	if err != nil {
		authUsecase.logger.Error("failed to get user by id", logger.Error(err))
		return CompleteRegisterResponse{}, apperrors.ErrInternalError
	}

	user.Email = reqDTO.Email
	user.FirstName = reqDTO.FirstName
	user.LastName = reqDTO.LastName
	user.IsRegistrationComplete = true

	if err := authUsecase.userService.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return CompleteRegisterResponse{}, apperrors.ErrEmailAlreadyExists
		}

		authUsecase.logger.Error("failed to update user", logger.Error(err))
		return CompleteRegisterResponse{}, apperrors.ErrInternalError
	}

	if err := authUsecase.authService.DeleteTempID(ctx, reqDTO.TempId); err != nil {
		authUsecase.logger.Error("failed to delete temp id", logger.Error(err))
	}

	if err := authUsecase.roleService.AssignRoleToUser(ctx, userId, models.RoleUser); err != nil {
		authUsecase.logger.Error("failed to assign role to user", logger.Error(err))
		return CompleteRegisterResponse{}, apperrors.ErrInternalError
	}

	accessToken, refreshToken, err := authUsecase.authService.GenerateJWT(ctx, user)
	if err != nil {
		authUsecase.logger.Error("failed to generate jwt", logger.Error(err))
		return CompleteRegisterResponse{}, apperrors.ErrInternalError
	}

	return NewCompleteRegisterResponse(accessToken, refreshToken, user.ID), nil
}
