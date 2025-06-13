package rr_intervals_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

// GetRRStatisticsRequest представляет запрос на получение статистики R-R интервалов.
type GetRRStatisticsRequest struct {
	From             time.Time
	To               time.Time
	IncludeHistogram bool // Включать ли данные гистограммы
	IncludeHRV       bool // Включать ли HRV метрики
	BinsCount        int  // Количество bins для гистограммы (если включена)
}

func NewGetRRStatisticsRequest(from, to time.Time) GetRRStatisticsRequest {
	return GetRRStatisticsRequest{
		From:             from,
		To:               to,
		IncludeHistogram: false,
		IncludeHRV:       false,
		BinsCount:        0, // Автоматический выбор
	}
}

// GetRRStatisticsResponse представляет ответ со статистикой R-R интервалов.
type GetRRStatisticsResponse struct {
	Summary    *models.RRStatisticalSummary
	Histogram  *models.RRHistogramData
	HRVMetrics *models.HRVMetrics
	TimeRange  TimeRange
}

// GetRRStatistics получает статистические данные R-R интервалов.
func (uc *RRIntervalsUsecase) GetRRStatistics(ctx context.Context, req GetRRStatisticsRequest) (GetRRStatisticsResponse, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return GetRRStatisticsResponse{}, apperrors.Forbidden()
	}

	// Валидируем временной диапазон
	if req.From.After(req.To) {
		return GetRRStatisticsResponse{}, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", req.From, req.To)
	}

	// Ограничиваем максимальный диапазон (30 дней для статистики)
	maxDuration := 30 * 24 * time.Hour
	if req.To.Sub(req.From) > maxDuration {
		return GetRRStatisticsResponse{}, apperrors.TimeRangeTooLargef("time range too large for statistics (max %v), got %v", maxDuration, req.To.Sub(req.From))
	}

	// Минимальный период (5 минут)
	minDuration := 5 * time.Minute
	if req.To.Sub(req.From) < minDuration {
		return GetRRStatisticsResponse{}, apperrors.TimeRangeTooSmallf("time range too small for statistics (min %v), got %v", minDuration, req.To.Sub(req.From))
	}

	// Получаем базовую статистику
	summary, err := uc.rrIntervalsService.GetStatisticalSummary(ctx, authClaims.UserID, req.From, req.To)
	if err != nil {
		uc.logger.Error("failed to get statistical summary",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))

		return GetRRStatisticsResponse{}, apperrors.DataProcessingErrorf("failed to get statistics: %v", err)
	}

	if summary.Count == 0 {
		return GetRRStatisticsResponse{}, apperrors.NoValidDataf("no RR intervals found for the specified time range")
	}

	response := GetRRStatisticsResponse{
		Summary: summary,
		TimeRange: TimeRange{
			From: req.From,
			To:   req.To,
		},
	}

	// Опционально получаем гистограмму
	if req.IncludeHistogram {
		histogram, err := uc.rrIntervalsService.BuildHistogram(ctx, authClaims.UserID, req.From, req.To, req.BinsCount)
		if err != nil {
			uc.logger.Warn("failed to build histogram for statistics",
				logger.String("user_id", authClaims.UserID.String()),
				logger.Error(err))
			// Не прерываем выполнение, просто не включаем гистограмму
		} else {
			response.Histogram = histogram
		}
	}

	// Опционально получаем HRV метрики
	if req.IncludeHRV {
		hrv, err := uc.rrIntervalsService.CalculateHRVMetrics(ctx, authClaims.UserID, req.From, req.To)
		if err != nil {
			uc.logger.Warn("failed to calculate HRV metrics for statistics",
				logger.String("user_id", authClaims.UserID.String()),
				logger.Error(err))
			// Не прерываем выполнение, просто не включаем HRV
		} else {
			response.HRVMetrics = hrv
		}
	}

	return response, nil
}

// GetAggregatedRRData получает агрегированные данные R-R интервалов по временным окнам.
func (uc *RRIntervalsUsecase) GetAggregatedRRData(ctx context.Context, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnProfile) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	if intervalMinutes < 1 || intervalMinutes > 60 {
		return nil, apperrors.ParameterOutOfRangef("interval must be between 1 and 60 minutes, got %d", intervalMinutes)
	}

	// Ограничиваем максимальный диапазон для агрегации
	maxDuration := 7 * 24 * time.Hour // 7 дней
	if to.Sub(from) > maxDuration {
		return nil, apperrors.TimeRangeTooLargef("time range too large for aggregation (max %v), got %v", maxDuration, to.Sub(from))
	}

	// Получаем агрегированные данные
	data, err := uc.rrIntervalsService.GetAggregatedByInterval(ctx, authClaims.UserID, from, to, intervalMinutes)
	if err != nil {
		uc.logger.Error("failed to get aggregated data",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))

		return nil, apperrors.DataProcessingErrorf("failed to get aggregated data: %v", err)
	}

	return data, nil
}
