package health_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/health_usecase"
	"github.com/gin-gonic/gin"
)

type HealthUsecase interface {
	CreateSteps(ctx context.Context, dto health_usecase.CreateStepsRequest) error
	GetSteps(ctx context.Context, dto health_usecase.GetStepsRequest) (health_usecase.GetStepsResponse, error)
}

type HealthRouter struct {
	healthUsecase   HealthUsecase
	identityService *indentity.IdentityService
}

func NewHealthRouter(healthUsecase HealthUsecase, identityService *indentity.IdentityService) *HealthRouter {
	return &HealthRouter{
		healthUsecase:   healthUsecase,
		identityService: identityService,
	}
}

func (hr *HealthRouter) Register(router gin.IRouter) {
	group := router.Group("/v1")
	group.Use(hr.identityService.AuthMiddleware())

	group.GET("/health/steps", hr.GetSteps)
	group.POST("/health/steps", hr.CreateSteps)

	group.GET("/health/heart-rates", nil)
	group.POST("/health/heart-rates", nil)

	group.GET("/health/temperatures", nil)
	group.POST("/health/temperatures", nil)

	group.GET("/health/weights", nil)
	group.POST("/health/weights", nil)

	group.GET("/health/sleeps", nil)
	group.POST("/health/sleeps", nil)
}
