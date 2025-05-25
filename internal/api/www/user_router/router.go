package user_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	GetUserById(ctx context.Context) (user_usecase.GetUserByIdResponse, error)
	UpdateSettingEmail(ctx context.Context, dto user_usecase.UpdateSettingEmailRequest) error
	GetSetting(ctx context.Context) (user_usecase.GetSettingResponse, error)
	UpdateSettingPassword(ctx context.Context, dto user_usecase.UpdateSettingPasswordRequest) error
	UpdateSettingProfile(ctx context.Context, dto user_usecase.UpdateSettingProfileRequest) error
	UpdateSettingTelegram(ctx context.Context, dto user_usecase.UpdateSettingTelegramRequest) error
	UpdateSettingTelegramUnlink(ctx context.Context) error
}

type UserRouter struct {
	userUsecase     UserUsecase
	identityService *indentity.IdentityService
}

func NewUserRouter(userUsecase UserUsecase, identityService *indentity.IdentityService) *UserRouter {
	return &UserRouter{
		userUsecase:     userUsecase,
		identityService: identityService,
	}
}

func (r *UserRouter) Register(router gin.IRouter) {
	group := router.Group("/v1")
	group.Use(r.identityService.AuthMiddleware())

	group.GET("/user", r.GetUserByIdRouter)

	group.GET("/user/setting", r.GetSetting)
	group.PATCH("/user/setting/email", r.UpdateSettingEmail)
	group.PATCH("/user/setting/password", r.UpdateSettingPassword)
	group.PATCH("/user/setting/telegram", r.UpdateSettingTelegram)
	group.DELETE("/user/setting/telegram", r.UpdateSettingTelegramUnlink)
	group.PATCH("/user/setting/profile", r.UpdateSettingProfile)
}
