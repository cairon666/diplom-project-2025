package user_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/pkg/logger"
)

type UpdateSettingProfileRequest struct {
	FirstName string
	LastName  string
}

func NewUpdateSettingProfileRequest(firstName string, lastName string) UpdateSettingProfileRequest {
	return UpdateSettingProfileRequest{
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (userUsecase *UserUsecase) UpdateSettingProfile(ctx context.Context, dto UpdateSettingProfileRequest) error {
	authClaims, err := checkWritePermissions(ctx)
	if err != nil {
		return err
	}

	user, err := userUsecase.userService.GetUserById(ctx, authClaims.UserID)
	if err != nil {
		userUsecase.logger.Error("failed to get user by id", logger.Error(err))

		return err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName

	err = userUsecase.userService.UpdateUser(ctx, user)
	if err != nil {
		userUsecase.logger.Error("failed to update user", logger.Error(err))

		return err
	}

	return nil
}
