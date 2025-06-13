package health_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

// HealthService определяет интерфейс для сервиса здоровья.
type HealthService interface {
	// CRUD операции
	CreateStep(ctx context.Context, step models.Step) error
	CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error
	CreateWeight(ctx context.Context, weight models.Weight) error
	CreateTemperature(ctx context.Context, temperature models.Temperature) error
	CreateSleep(ctx context.Context, sleep models.Sleep) error

	// Batch операции
	CreateSteps(ctx context.Context, steps []models.Step) error
	CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error
	CreateWeights(ctx context.Context, weights []models.Weight) error
	CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error
	CreateSleeps(ctx context.Context, sleeps []models.Sleep) error

	// Получение сырых данных
	GetSteps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Step, error)
	GetHeartRates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.HeartRate, error)
	GetWeights(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Weight, error)
	GetTemperatures(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Temperature, error)
	GetSleeps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Sleep, error)
}

// AggregationService определяет интерфейс для сервиса агрегации.
type AggregationService interface {
	GetHourlySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error)
	GetDailySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error)
	GetHourlyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailyWeightAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetHourlyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailySleepDuration(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
}

// HealthUsecase реализует бизнес-логику для работы с данными о здоровье.
type HealthUsecase struct {
	healthService      HealthService
	aggregationService AggregationService
	logger             logger.ILogger
}

// NewHealthUsecase создает новый экземпляр usecase здоровья.
func NewHealthUsecase(healthService HealthService, aggregationService AggregationService, logger logger.ILogger) *HealthUsecase {
	return &HealthUsecase{
		healthService:      healthService,
		aggregationService: aggregationService,
		logger:             logger,
	}
}
