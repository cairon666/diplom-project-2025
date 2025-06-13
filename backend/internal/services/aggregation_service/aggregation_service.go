package aggregation_service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// HealthAggregationRepo определяет интерфейс для репозитория агрегированных данных о здоровье.
type HealthAggregationRepo interface {
	// Агрегированные данные по шагам
	GetHourlySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error)
	GetDailySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error)

	// Агрегированные данные по пульсу
	GetHourlyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)

	// Агрегированные данные по весу
	GetDailyWeightAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)

	// Агрегированные данные по температуре
	GetHourlyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
	GetDailyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)

	// Агрегированные данные по сну
	GetDailySleepDuration(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error)
}

// AggregationService реализует бизнес-логику для работы с агрегированными данными о здоровье.
type AggregationService struct {
	healthAggregationRepo HealthAggregationRepo
}

// NewAggregationService создает новый экземпляр сервиса агрегации.
func NewAggregationService(healthAggregationRepo HealthAggregationRepo) *AggregationService {
	return &AggregationService{
		healthAggregationRepo: healthAggregationRepo,
	}
}

// GetHourlySteps получает агрегированные данные по шагам по часам.
func (s *AggregationService) GetHourlySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error) {
	return s.healthAggregationRepo.GetHourlySteps(ctx, userID, from, to)
}

// GetDailySteps получает агрегированные данные по шагам по дням.
func (s *AggregationService) GetDailySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error) {
	return s.healthAggregationRepo.GetDailySteps(ctx, userID, from, to)
}

// GetHourlyHeartRateAvg получает средний пульс по часам.
func (s *AggregationService) GetHourlyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetHourlyHeartRateAvg(ctx, userID, from, to)
}

// GetDailyHeartRateAvg получает средний пульс по дням.
func (s *AggregationService) GetDailyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetDailyHeartRateAvg(ctx, userID, from, to)
}

// GetDailyWeightAvg получает средний вес по дням.
func (s *AggregationService) GetDailyWeightAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetDailyWeightAvg(ctx, userID, from, to)
}

// GetHourlyTemperatureAvg получает среднюю температуру по часам.
func (s *AggregationService) GetHourlyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetHourlyTemperatureAvg(ctx, userID, from, to)
}

// GetDailyTemperatureAvg получает среднюю температуру по дням.
func (s *AggregationService) GetDailyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetDailyTemperatureAvg(ctx, userID, from, to)
}

// GetDailySleepDuration получает продолжительность сна по дням.
func (s *AggregationService) GetDailySleepDuration(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	return s.healthAggregationRepo.GetDailySleepDuration(ctx, userID, from, to)
}
