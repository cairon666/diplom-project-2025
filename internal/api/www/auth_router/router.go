package auth_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	Login(ctx context.Context, loginReq auth_usecase.LoginRequest) (auth_usecase.LoginResponse, error)
	Register(ctx context.Context, registerReq auth_usecase.RegisterRequest) error
	RefreshToken(ctx context.Context, refreshTokenReq auth_usecase.RefreshTokenRequest) (auth_usecase.RefreshTokenResponse, error)
	TelegramLogin(ctx context.Context, telegramLoginReq auth_usecase.TelegramLoginRequest) (auth_usecase.TelegramLoginResponse, error)
	TelegramConfirm(ctx context.Context, reqDTO auth_usecase.CompleteRegisterRequest) (auth_usecase.CompleteRegisterResponse, error)
}

type AuthRouter struct {
	authUsecase AuthUsecase
	jwtService  *indentity.JWTService
}

func NewAuthRouter(authUsecase AuthUsecase, jwtService *indentity.JWTService) *AuthRouter {
	return &AuthRouter{
		authUsecase: authUsecase,
		jwtService:  jwtService,
	}
}

func (r *AuthRouter) Register(router gin.IRouter) {
	router.POST("/v1/auth/register", r.RegisterRoute)
	router.POST("/v1/auth/login", r.LoginRoute)
	router.POST("/v1/auth/refresh", r.RefreshRoute)
	router.POST("/v1/auth/telegram/login", r.TelegramLoginRoute)
	router.POST("/v1/auth/telegram/confirm-register", r.TelegramConfirmRoute)
}
