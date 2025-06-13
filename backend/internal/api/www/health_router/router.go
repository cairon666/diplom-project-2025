package health_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/health_usecase"
	"github.com/gin-gonic/gin"
)

// HealthUsecase определяет интерфейс для usecase здоровья.
type HealthUsecase interface {
	// Steps
	CreateStep(ctx context.Context, dto health_usecase.CreateStepRequest) (health_usecase.CreateStepResponse, error)
	CreateSteps(ctx context.Context, dto health_usecase.CreateStepsRequest) (health_usecase.CreateStepsResponse, error)
	GetSteps(ctx context.Context, dto health_usecase.GetStepsRequest) (health_usecase.GetStepsResponse, error)
	GetHourlySteps(ctx context.Context, dto health_usecase.GetHourlyStepsRequest) (health_usecase.GetHourlyStepsResponse, error)
	GetDailySteps(ctx context.Context, dto health_usecase.GetDailyStepsRequest) (health_usecase.GetDailyStepsResponse, error)

	// Weights
	CreateWeight(ctx context.Context, dto health_usecase.CreateWeightRequest) (health_usecase.CreateWeightResponse, error)
	CreateWeights(ctx context.Context, dto health_usecase.CreateWeightsRequest) (health_usecase.CreateWeightsResponse, error)
	GetWeights(ctx context.Context, dto health_usecase.GetWeightsRequest) (health_usecase.GetWeightsResponse, error)
	GetDailyWeightAvg(ctx context.Context, dto health_usecase.GetDailyWeightAvgRequest) (health_usecase.GetDailyWeightAvgResponse, error)

	// Temperatures
	CreateTemperature(ctx context.Context, dto health_usecase.CreateTemperatureRequest) (health_usecase.CreateTemperatureResponse, error)
	CreateTemperatures(ctx context.Context, dto health_usecase.CreateTemperaturesRequest) (health_usecase.CreateTemperaturesResponse, error)
	GetTemperatures(ctx context.Context, dto health_usecase.GetTemperaturesRequest) (health_usecase.GetTemperaturesResponse, error)
	GetHourlyTemperatureAvg(ctx context.Context, dto health_usecase.GetHourlyTemperatureAvgRequest) (health_usecase.GetHourlyTemperatureAvgResponse, error)
	GetDailyTemperatureAvg(ctx context.Context, dto health_usecase.GetDailyTemperatureAvgRequest) (health_usecase.GetDailyTemperatureAvgResponse, error)

	// Sleeps
	CreateSleep(ctx context.Context, dto health_usecase.CreateSleepRequest) (health_usecase.CreateSleepResponse, error)
	CreateSleeps(ctx context.Context, dto health_usecase.CreateSleepsRequest) (health_usecase.CreateSleepsResponse, error)
	GetSleeps(ctx context.Context, dto health_usecase.GetSleepsRequest) (health_usecase.GetSleepsResponse, error)
	GetDailySleepDuration(ctx context.Context, dto health_usecase.GetDailySleepDurationRequest) (health_usecase.GetDailySleepDurationResponse, error)
}

// HealthRouter реализует роутер для эндпоинтов здоровья.
type HealthRouter struct {
	healthUsecase   HealthUsecase
	identityService *indentity.IdentityService
}

// NewHealthRouter создает новый экземпляр роутера здоровья.
func NewHealthRouter(healthUsecase HealthUsecase, identityService *indentity.IdentityService) www.Router {
	return &HealthRouter{
		healthUsecase:   healthUsecase,
		identityService: identityService,
	}
}

// Register регистрирует маршруты здоровья.
func (r *HealthRouter) Register(router gin.IRouter) {
	healthGroup := router.Group("/v1/health")
	healthGroup.Use(r.identityService.AuthMiddleware())
	{
		// Steps endpoints
		stepsGroup := healthGroup.Group("/steps")
		{
			stepsGroup.POST("", r.CreateStep)
			stepsGroup.POST("/batch", r.CreateSteps)
			stepsGroup.GET("", r.GetSteps)
			stepsGroup.GET("/hourly", r.GetHourlySteps)
			stepsGroup.GET("/daily", r.GetDailySteps)
		}

		// Weights endpoints
		weightsGroup := healthGroup.Group("/weights")
		{
			weightsGroup.POST("", r.CreateWeight)
			weightsGroup.POST("/batch", r.CreateWeights)
			weightsGroup.GET("", r.GetWeights)
			weightsGroup.GET("/daily", r.GetDailyWeightAvg)
		}

		// Temperatures endpoints
		temperaturesGroup := healthGroup.Group("/temperatures")
		{
			temperaturesGroup.POST("", r.CreateTemperature)
			temperaturesGroup.POST("/batch", r.CreateTemperatures)
			temperaturesGroup.GET("", r.GetTemperatures)
			temperaturesGroup.GET("/hourly", r.GetHourlyTemperatureAvg)
			temperaturesGroup.GET("/daily", r.GetDailyTemperatureAvg)
		}

		// Sleeps endpoints
		sleepsGroup := healthGroup.Group("/sleeps")
		{
			sleepsGroup.POST("", r.CreateSleep)
			sleepsGroup.POST("/batch", r.CreateSleeps)
			sleepsGroup.GET("", r.GetSleeps)
			sleepsGroup.GET("/daily", r.GetDailySleepDuration)
		}
	}
}
