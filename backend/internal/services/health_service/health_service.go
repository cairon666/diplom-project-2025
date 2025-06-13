package health_service

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// HealthDataRepo определяет интерфейс для репозитория данных о здоровье.
type HealthDataRepo interface {
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

// HealthService реализует бизнес-логику для работы с данными о здоровье.
type HealthService struct {
	healthDataRepo HealthDataRepo
}

// NewHealthService создает новый экземпляр сервиса.
func NewHealthService(healthDataRepo HealthDataRepo) *HealthService {
	return &HealthService{
		healthDataRepo: healthDataRepo,
	}
}

// CreateStep создает запись о шагах.
func (s *HealthService) CreateStep(ctx context.Context, step models.Step) error {
	return s.healthDataRepo.CreateStep(ctx, step)
}

// CreateHeartRate создает запись о пульсе.
func (s *HealthService) CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error {
	return s.healthDataRepo.CreateHeartRate(ctx, heartRate)
}

// CreateWeight создает запись о весе.
func (s *HealthService) CreateWeight(ctx context.Context, weight models.Weight) error {
	return s.healthDataRepo.CreateWeight(ctx, weight)
}

// CreateTemperature создает запись о температуре.
func (s *HealthService) CreateTemperature(ctx context.Context, temperature models.Temperature) error {
	return s.healthDataRepo.CreateTemperature(ctx, temperature)
}

// CreateSleep создает запись о сне.
func (s *HealthService) CreateSleep(ctx context.Context, sleep models.Sleep) error {
	return s.healthDataRepo.CreateSleep(ctx, sleep)
}

// CreateSteps создает множественные записи о шагах.
func (s *HealthService) CreateSteps(ctx context.Context, steps []models.Step) error {
	return s.healthDataRepo.CreateSteps(ctx, steps)
}

// CreateHeartRates создает множественные записи о пульсе.
func (s *HealthService) CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error {
	return s.healthDataRepo.CreateHeartRates(ctx, heartRates)
}

// CreateWeights создает множественные записи о весе.
func (s *HealthService) CreateWeights(ctx context.Context, weights []models.Weight) error {
	return s.healthDataRepo.CreateWeights(ctx, weights)
}

// CreateTemperatures создает множественные записи о температуре.
func (s *HealthService) CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error {
	return s.healthDataRepo.CreateTemperatures(ctx, temperatures)
}

// CreateSleeps создает множественные записи о сне.
func (s *HealthService) CreateSleeps(ctx context.Context, sleeps []models.Sleep) error {
	return s.healthDataRepo.CreateSleeps(ctx, sleeps)
}

// GetSteps получает сырые данные о шагах.
func (s *HealthService) GetSteps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Step, error) {
	return s.healthDataRepo.GetSteps(ctx, userID, from, to)
}

// GetHeartRates получает сырые данные о пульсе.
func (s *HealthService) GetHeartRates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.HeartRate, error) {
	return s.healthDataRepo.GetHeartRates(ctx, userID, from, to)
}

// GetWeights получает сырые данные о весе.
func (s *HealthService) GetWeights(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Weight, error) {
	return s.healthDataRepo.GetWeights(ctx, userID, from, to)
}

// GetTemperatures получает сырые данные о температуре.
func (s *HealthService) GetTemperatures(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Temperature, error) {
	return s.healthDataRepo.GetTemperatures(ctx, userID, from, to)
}

// GetSleeps получает сырые данные о сне.
func (s *HealthService) GetSleeps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Sleep, error) {
	return s.healthDataRepo.GetSleeps(ctx, userID, from, to)
}
