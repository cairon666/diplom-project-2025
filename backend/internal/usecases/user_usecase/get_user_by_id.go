package user_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/google/uuid"
)

type GetUserByIdResponse struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

func NewGetUserByIdResponse(id uuid.UUID, email, firstName, lastName string) GetUserByIdResponse {
	return GetUserByIdResponse{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (userUsecase *UserUsecase) GetUserById(ctx context.Context) (GetUserByIdResponse, error) {
	authClaims, err := checkReadPermissions(ctx)
	if err != nil {
		return GetUserByIdResponse{}, err
	}

	user, err := userUsecase.userService.GetUserById(ctx, authClaims.UserID)
	if err != nil {
		return GetUserByIdResponse{}, apperrors.UserNotFoundf("failed to get user by id: %v", err)
	}

	resDTO := NewGetUserByIdResponse(user.ID, user.Email, user.FirstName, user.LastName)

	return resDTO, nil
}
