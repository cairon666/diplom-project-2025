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

// GetRRHistogram получает гистограмму R-R интервалов
func (uc *RRIntervalsUsecase) GetRRHistogram(ctx context.Context, from, to time.Time, binsCount int) (*models.RRHistogramData, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	if binsCount < 0 || binsCount > 50 {
		return nil, apperrors.ParameterOutOfRangef("bins count must be between 0 and 50, got %d", binsCount)
	}

	// Проверяем минимальную продолжительность для анализа
	minDuration := 1 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for histogram analysis (min %v), got %v", minDuration, to.Sub(from))
	}

	// Строим гистограмму
	histogram, err := uc.rrIntervalsService.BuildHistogram(ctx, authClaims.UserID, from, to, binsCount)
	if err != nil {
		uc.logger.Error("failed to build histogram",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))
		return nil, apperrors.AnalysisNotPossiblef("failed to build histogram: %v", err)
	}

	return histogram, nil
}

// GetRRTrends получает анализ трендов R-R интервалов
func (uc *RRIntervalsUsecase) GetRRTrends(ctx context.Context, from, to time.Time, windowSize int) (*models.RRTrendAnalysis, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	if windowSize < 1 || windowSize > 60 {
		return nil, apperrors.ParameterOutOfRangef("window size must be between 1 and 60 minutes, got %d", windowSize)
	}

	// Проверяем минимальную продолжительность для анализа трендов
	minDuration := time.Duration(windowSize*2) * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for trend analysis (min %v for window size %d), got %v", minDuration, windowSize, to.Sub(from))
	}

	// Анализируем тренды
	trends, err := uc.rrIntervalsService.AnalyzeTrends(ctx, authClaims.UserID, from, to, windowSize)
	if err != nil {
		uc.logger.Error("failed to analyze trends",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))
		return nil, apperrors.AnalysisNotPossiblef("failed to analyze trends: %v", err)
	}

	return trends, nil
}

// GetHRVMetrics получает метрики вариабельности сердечного ритма
func (uc *RRIntervalsUsecase) GetHRVMetrics(ctx context.Context, from, to time.Time) (*models.HRVMetrics, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Для HRV анализа нужен минимальный период
	minDuration := 5 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for HRV analysis (min %v), got %v", minDuration, to.Sub(from))
	}

	// Вычисляем HRV метрики
	hrv, err := uc.rrIntervalsService.CalculateHRVMetrics(ctx, authClaims.UserID, from, to)
	if err != nil {
		uc.logger.Error("failed to calculate HRV metrics",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))
		return nil, apperrors.AnalysisNotPossiblef("failed to calculate HRV metrics: %v", err)
	}

	return hrv, nil
}

// GetRRDifferentialHistogram получает дифференциальную гистограмму R-R интервалов
func (uc *RRIntervalsUsecase) GetRRDifferentialHistogram(ctx context.Context, from, to time.Time, binsCount int) (*models.DifferentialHistogramData, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	if binsCount < 0 || binsCount > 50 {
		return nil, apperrors.ParameterOutOfRangef("bins count must be between 0 and 50, got %d", binsCount)
	}

	// Для дифференциальной гистограммы нужно минимум 2 точки
	minDuration := 10 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for differential histogram (min %v), got %v", minDuration, to.Sub(from))
	}

	// Строим дифференциальную гистограмму
	histogram, err := uc.rrIntervalsService.BuildDifferentialHistogram(ctx, authClaims.UserID, from, to, binsCount)
	if err != nil {
		uc.logger.Error("failed to build differential histogram",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))
		return nil, apperrors.AnalysisNotPossiblef("failed to build differential histogram: %v", err)
	}

	return histogram, nil
}

// GetRRScatterplot получает скаттерограмму (диаграмму Пуанкаре) R-R интервалов
func (uc *RRIntervalsUsecase) GetRRScatterplot(ctx context.Context, from, to time.Time) (*models.ScatterplotData, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return nil, apperrors.Forbidden()
	}

	// Валидируем параметры
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Для скаттерограммы нужно минимум 2 точки
	minDuration := 10 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for scatterplot (min %v), got %v", minDuration, to.Sub(from))
	}

	// Строим скаттерограмму
	scatterplot, err := uc.rrIntervalsService.BuildScatterplot(ctx, authClaims.UserID, from, to)
	if err != nil {
		uc.logger.Error("failed to build scatterplot",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Error(err))
		return nil, apperrors.AnalysisNotPossiblef("failed to build scatterplot: %v", err)
	}

	return scatterplot, nil
} 