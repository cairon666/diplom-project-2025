package user_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/pkg/logger"
)

type UpdateSettingEmailRequest struct {
	Email string
}

func NewUpdateSettingEmailRequest(email string) UpdateSettingEmailRequest {
	return UpdateSettingEmailRequest{
		Email: email,
	}

}

func (userUsecase *UserUsecase) UpdateSettingEmail(ctx context.Context, dto UpdateSettingEmailRequest) error {
	authClaims, err := checkWritePermissions(ctx)
	if err != nil {
		return err
	}

	user, err := userUsecase.userService.GetUserById(ctx, authClaims.UserID)
	if err != nil {
		userUsecase.logger.Error("failed to get user by id", logger.Error(err))
		return err
	}

	user.Email = dto.Email
	err = userUsecase.userService.UpdateUser(ctx, user)
	if err != nil {
		userUsecase.logger.Error("failed to update user", logger.Error(err))
		return err
	}

	return nil
}
