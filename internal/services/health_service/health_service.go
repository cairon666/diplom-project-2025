package health_service

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type HealthRepo interface {
	GetHeartRates(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.HeartRate, error)
	CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error
	CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error
	GetSteps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Step, error)
	CreateStep(ctx context.Context, step models.Step) error
	CreateSteps(ctx context.Context, steps []models.Step) error
	GetSleeps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Sleep, error)
	CreateSleep(ctx context.Context, sleep models.Sleep) error
	CreateSleeps(ctx context.Context, sleeps []models.Sleep) error
	GetTemperatures(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Temperature, error)
	CreateTemperature(ctx context.Context, temperature models.Temperature) error
	CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error
	GetWeights(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Weight, error)
	CreateWeight(ctx context.Context, weight models.Weight) error
	CreateWeights(ctx context.Context, weights []models.Weight) error
}

type HealthService struct {
	healthRepo HealthRepo
}

func NewHealthService(healthRepo HealthRepo) *HealthService {
	return &HealthService{
		healthRepo: healthRepo,
	}
}

func (hs *HealthService) GetHeartRates(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.HeartRate, error) {
	return hs.healthRepo.GetHeartRates(ctx, userID, from, to)
}

func (hs *HealthService) CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error {
	return hs.healthRepo.CreateHeartRate(ctx, heartRate)
}

func (hs *HealthService) CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error {
	return hs.healthRepo.CreateHeartRates(ctx, heartRates)
}

func (hs *HealthService) GetSteps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Step, error) {
	return hs.healthRepo.GetSteps(ctx, userID, from, to)
}

func (hs *HealthService) CreateStep(ctx context.Context, step models.Step) error {
	return hs.healthRepo.CreateStep(ctx, step)
}

func (hs *HealthService) CreateSteps(ctx context.Context, steps []models.Step) error {
	return hs.healthRepo.CreateSteps(ctx, steps)
}

func (hs *HealthService) GetSleeps(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Sleep, error) {
	return hs.healthRepo.GetSleeps(ctx, userID, from, to)
}

func (hs *HealthService) CreateSleep(ctx context.Context, sleep models.Sleep) error {
	return hs.healthRepo.CreateSleep(ctx, sleep)
}

func (hs *HealthService) CreateSleeps(ctx context.Context, sleeps []models.Sleep) error {
	return hs.healthRepo.CreateSleeps(ctx, sleeps)
}

func (hs *HealthService) GetTemperatures(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Temperature, error) {
	return hs.healthRepo.GetTemperatures(ctx, userID, from, to)
}

func (hs *HealthService) CreateTemperature(ctx context.Context, temperature models.Temperature) error {
	return hs.healthRepo.CreateTemperature(ctx, temperature)
}

func (hs *HealthService) CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error {
	return hs.healthRepo.CreateTemperatures(ctx, temperatures)
}

func (hs *HealthService) GetWeights(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]models.Weight, error) {
	return hs.healthRepo.GetWeights(ctx, userID, from, to)
}

func (hs *HealthService) CreateWeight(ctx context.Context, weight models.Weight) error {
	return hs.healthRepo.CreateWeight(ctx, weight)
}

func (hs *HealthService) CreateWeights(ctx context.Context, weights []models.Weight) error {
	return hs.healthRepo.CreateWeights(ctx, weights)
}
