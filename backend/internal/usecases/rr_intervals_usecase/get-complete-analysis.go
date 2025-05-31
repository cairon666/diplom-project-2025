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

// GetCompleteAnalysisRequest представляет запрос для комплексного анализа RR интервалов
type GetCompleteAnalysisRequest struct {
	From    time.Time                       `json:"from"`
	To      time.Time                       `json:"to"`
	Options *models.CompleteAnalysisOptions `json:"options,omitempty"`
}

// GetCompleteAnalysisResponse представляет ответ комплексного анализа RR интервалов
type GetCompleteAnalysisResponse struct {
	Data *models.CompleteAnalysisData `json:"data"`
}

// GetCompleteAnalysis выполняет комплексный анализ RR интервалов за один оптимизированный запрос
// Этот метод заменяет множественные вызовы отдельных аналитических методов
func (uc *RRIntervalsUsecase) GetCompleteAnalysis(ctx context.Context, req GetCompleteAnalysisRequest) (GetCompleteAnalysisResponse, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return GetCompleteAnalysisResponse{}, apperrors.Forbidden()
	}

	// Валидируем параметры времени
	if req.From.After(req.To) {
		return GetCompleteAnalysisResponse{}, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", req.From, req.To)
	}

	// Проверяем минимальную продолжительность для комплексного анализа
	minDuration := 2 * time.Minute
	duration := req.To.Sub(req.From)
	if duration < minDuration {
		return GetCompleteAnalysisResponse{}, apperrors.TimeRangeTooSmallf("time range too small for complete analysis (min %v), got %v", minDuration, duration)
	}

	// Ограничиваем максимальную продолжительность во избежание перегрузки системы
	maxDuration := 24 * time.Hour
	if duration > maxDuration {
		return GetCompleteAnalysisResponse{}, apperrors.TimeRangeTooLargef("time range too large for complete analysis (max %v), got %v", maxDuration, duration)
	}

	// Валидируем и нормализуем опции анализа
	options := req.Options
	if options == nil {
		options = uc.getDefaultCompleteAnalysisOptions()
	} else {
		if err := uc.validateCompleteAnalysisOptions(options); err != nil {
			return GetCompleteAnalysisResponse{}, err
		}
	}

	// Выполняем комплексный анализ
	analysisData, err := uc.rrIntervalsService.GetCompleteAnalysis(ctx, authClaims.UserID, req.From, req.To, options)
	if err != nil {
		uc.logger.Error("failed to perform complete analysis",
			logger.String("user_id", authClaims.UserID.String()),
			logger.Time("from", req.From),
			logger.Time("to", req.To),
			logger.Error(err))
		return GetCompleteAnalysisResponse{}, apperrors.AnalysisNotPossiblef("failed to perform complete analysis: %v", err)
	}

	// Логируем успешное выполнение с метриками производительности
	uc.logger.Info("complete analysis performed successfully",
		logger.String("user_id", authClaims.UserID.String()),
		logger.Duration("duration", duration),
		logger.Duration("processing_time", analysisData.ProcessingTime),
		logger.Int64("raw_values_count", int64(len(analysisData.RawValues))),
		logger.Int64("aggregated_points_count", int64(len(analysisData.AggregatedData))),
		logger.Float64("data_quality", analysisData.DataQuality.QualityPercentage))

	return GetCompleteAnalysisResponse{
		Data: analysisData,
	}, nil
}

// getDefaultCompleteAnalysisOptions возвращает опции по умолчанию для комплексного анализа
func (uc *RRIntervalsUsecase) getDefaultCompleteAnalysisOptions() *models.CompleteAnalysisOptions {
	return &models.CompleteAnalysisOptions{
		AggregationIntervalMinutes:    5,     // 5-минутные интервалы агрегации
		TrendWindowSizeMinutes:       15,     // 15-минутное окно для анализа трендов
		HistogramBinsCount:           25,     // 25 bins для основной гистограммы
		DiffHistogramBinsCount:       20,     // 20 bins для дифференциальной гистограммы
		EnableFrequencyDomainAnalysis: false, // Частотный анализ отключен для производительности
		IncludeRawData:               true,   // Включаем сырые данные
		MaxDataPoints:                10000,  // Лимит на количество точек данных
	}
}

// validateCompleteAnalysisOptions валидирует опции комплексного анализа
func (uc *RRIntervalsUsecase) validateCompleteAnalysisOptions(options *models.CompleteAnalysisOptions) error {
	// Валидация интервала агрегации
	if options.AggregationIntervalMinutes < 1 || options.AggregationIntervalMinutes > 60 {
		return apperrors.ParameterOutOfRangef("aggregation interval must be between 1 and 60 minutes, got %d", options.AggregationIntervalMinutes)
	}

	// Валидация окна для анализа трендов
	if options.TrendWindowSizeMinutes < 5 || options.TrendWindowSizeMinutes > 120 {
		return apperrors.ParameterOutOfRangef("trend window size must be between 5 and 120 minutes, got %d", options.TrendWindowSizeMinutes)
	}

	// Валидация количества bins для гистограмм
	if options.HistogramBinsCount < 5 || options.HistogramBinsCount > 100 {
		return apperrors.ParameterOutOfRangef("histogram bins count must be between 5 and 100, got %d", options.HistogramBinsCount)
	}

	if options.DiffHistogramBinsCount < 5 || options.DiffHistogramBinsCount > 100 {
		return apperrors.ParameterOutOfRangef("differential histogram bins count must be between 5 and 100, got %d", options.DiffHistogramBinsCount)
	}

	// Валидация лимита точек данных
	if options.MaxDataPoints < 100 || options.MaxDataPoints > 100000 {
		return apperrors.ParameterOutOfRangef("max data points must be between 100 and 100000, got %d", options.MaxDataPoints)
	}

	return nil
} 