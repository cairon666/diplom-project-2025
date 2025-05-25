package user_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GetUserByIdResponse struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

func NewGetUserByIdResponse(ID uuid.UUID, email string, firstName string, lastName string) GetUserByIdResponse {
	return GetUserByIdResponse{
		ID:        ID,
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
		return GetUserByIdResponse{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	resDTO := NewGetUserByIdResponse(user.ID, user.Email, user.FirstName, user.LastName)
	return resDTO, nil
}
