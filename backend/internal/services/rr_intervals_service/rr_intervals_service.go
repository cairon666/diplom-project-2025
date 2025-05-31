package rr_intervals_service

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type RRIntervalsRepo interface {
	CreateBatch(ctx context.Context, rrIntervals []models.RRInterval) error
	GetByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.RRInterval, error)
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID, from, to time.Time) ([]models.RRInterval, error)
	GetRawValuesForAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]int64, error)
	GetAggregatedByInterval(ctx context.Context, userID uuid.UUID, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error)
	GetStatisticalSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.RRStatisticalSummary, error)
	GetTimeSeriesForTrends(ctx context.Context, userID uuid.UUID, from, to time.Time, windowSizeMinutes int) ([]models.TrendPoint, error)
	GetCompleteAnalysisData(ctx context.Context, userID uuid.UUID, from, to time.Time, options models.CompleteAnalysisOptions) (*models.CompleteAnalysisData, error)
}

type RRIntervalsService struct {
	rrIntervalsRepo RRIntervalsRepo
}

func NewRRIntervalsService(rrIntervalsRepo RRIntervalsRepo) *RRIntervalsService {
	return &RRIntervalsService{
		rrIntervalsRepo: rrIntervalsRepo,
	}
}

// CRUD операции

func (s *RRIntervalsService) CreateBatch(ctx context.Context, rrIntervals []models.RRInterval) error {
	return s.rrIntervalsRepo.CreateBatch(ctx, rrIntervals)
}

func (s *RRIntervalsService) GetByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.RRInterval, error) {
	return s.rrIntervalsRepo.GetByUserID(ctx, userID, from, to)
}

func (s *RRIntervalsService) GetByDeviceID(ctx context.Context, deviceID uuid.UUID, from, to time.Time) ([]models.RRInterval, error) {
	return s.rrIntervalsRepo.GetByDeviceID(ctx, deviceID, from, to)
}

// Аналитические операции

func (s *RRIntervalsService) GetRawValuesForAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]int64, error) {
	return s.rrIntervalsRepo.GetRawValuesForAnalysis(ctx, userID, from, to)
}

func (s *RRIntervalsService) GetAggregatedByInterval(ctx context.Context, userID uuid.UUID, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error) {
	return s.rrIntervalsRepo.GetAggregatedByInterval(ctx, userID, from, to, intervalMinutes)
}

func (s *RRIntervalsService) GetStatisticalSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.RRStatisticalSummary, error) {
	return s.rrIntervalsRepo.GetStatisticalSummary(ctx, userID, from, to)
}

func (s *RRIntervalsService) GetTimeSeriesForTrends(ctx context.Context, userID uuid.UUID, from, to time.Time, windowSizeMinutes int) ([]models.TrendPoint, error) {
	return s.rrIntervalsRepo.GetTimeSeriesForTrends(ctx, userID, from, to, windowSizeMinutes)
} 