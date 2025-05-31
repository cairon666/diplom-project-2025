package rr_intervals_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type RRIntervalsService interface {
	// CRUD операции
	CreateBatch(ctx context.Context, rrIntervals []models.RRInterval) error
	GetByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.RRInterval, error)
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID, from, to time.Time) ([]models.RRInterval, error)
	
	// Аналитические операции (низкоуровневые)
	GetRawValuesForAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]int64, error)
	GetAggregatedByInterval(ctx context.Context, userID uuid.UUID, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error)
	GetStatisticalSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.RRStatisticalSummary, error)
	
	// Высокоуровневые аналитические методы
	BuildHistogram(ctx context.Context, userID uuid.UUID, from, to time.Time, binsCount int) (*models.RRHistogramData, error)
	BuildDifferentialHistogram(ctx context.Context, userID uuid.UUID, from, to time.Time, binsCount int) (*models.DifferentialHistogramData, error)
	BuildScatterplot(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.ScatterplotData, error)
	AnalyzeTrends(ctx context.Context, userID uuid.UUID, from, to time.Time, windowSizeMinutes int) (*models.RRTrendAnalysis, error)
	CalculateHRVMetrics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.HRVMetrics, error)
	
	// Комплексный анализ (новый оптимизированный метод)
	GetCompleteAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time, options *models.CompleteAnalysisOptions) (*models.CompleteAnalysisData, error)
}

type DeviceService interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Device, error)
}

type RRIntervalsUsecase struct {
	rrIntervalsService RRIntervalsService
	deviceService      DeviceService
	logger             logger.ILogger
}

func NewRRIntervalsUsecase(
	rrIntervalsService RRIntervalsService,
	deviceService DeviceService,
	logger logger.ILogger,
) *RRIntervalsUsecase {
	return &RRIntervalsUsecase{
		rrIntervalsService: rrIntervalsService,
		deviceService:      deviceService,
		logger:             logger,
	}
}