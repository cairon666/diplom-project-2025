package health_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type HealthService interface {
	GetHeartRates(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.HeartRate, error)
	CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error
	GetSteps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Step, error)
	CreateSteps(ctx context.Context, steps []models.Step) error
	GetSleeps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Sleep, error)
	CreateSleeps(ctx context.Context, sleeps []models.Sleep) error
	GetTemperatures(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Temperature, error)
	CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error
	GetWeights(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Weight, error)
	CreateWeights(ctx context.Context, weights []models.Weight) error
}

type HealthUsecase struct {
	healthService HealthService
	logger        *logger.Logger
}

func NewHealthUsecase(healthService HealthService, logger *logger.Logger) *HealthUsecase {
	return &HealthUsecase{
		healthService: healthService,
		logger:        logger,
	}
}
