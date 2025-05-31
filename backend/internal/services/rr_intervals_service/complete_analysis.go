package rr_intervals_service

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// GetCompleteAnalysis выполняет полный комплексный анализ RR интервалов за один оптимизированный запрос
// Этот метод заменяет множественные вызовы отдельных методов анализа
func (s *RRIntervalsService) GetCompleteAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time, options *models.CompleteAnalysisOptions) (*models.CompleteAnalysisData, error) {
	// Устанавливаем дефолтные значения для опций
	if options == nil {
		options = s.getDefaultAnalysisOptions()
	} else {
		s.validateAndSetDefaults(options)
	}

	// Проверяем минимальную продолжительность для анализа
	minDuration := 2 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for complete analysis (min %v), got %v", minDuration, to.Sub(from))
	}

	// Получаем базовые данные и предвычисленную статистику через оптимизированный запрос
	analysisData, err := s.rrIntervalsRepo.GetCompleteAnalysisData(ctx, userID, from, to, *options)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get complete analysis data: %v", err)
	}

	// Дополняем данные вычислениями на уровне сервиса
	if err := s.enrichAnalysisData(ctx, analysisData, options); err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to enrich analysis data: %v", err)
	}

	return analysisData, nil
}

// getDefaultAnalysisOptions возвращает дефолтные опции для анализа
func (s *RRIntervalsService) getDefaultAnalysisOptions() *models.CompleteAnalysisOptions {
	return &models.CompleteAnalysisOptions{
		AggregationIntervalMinutes:    5,
		TrendWindowSizeMinutes:       15,
		HistogramBinsCount:           25,
		DiffHistogramBinsCount:       20,
		EnableFrequencyDomainAnalysis: false, // Пока отключено для производительности
		IncludeRawData:               true,
		MaxDataPoints:                10000,
	}
}

// validateAndSetDefaults проверяет и устанавливает дефолтные значения для опций
func (s *RRIntervalsService) validateAndSetDefaults(options *models.CompleteAnalysisOptions) {
	if options.AggregationIntervalMinutes <= 0 {
		options.AggregationIntervalMinutes = 5
	}
	if options.TrendWindowSizeMinutes <= 0 {
		options.TrendWindowSizeMinutes = 15
	}
	if options.HistogramBinsCount <= 0 {
		options.HistogramBinsCount = 25
	}
	if options.DiffHistogramBinsCount <= 0 {
		options.DiffHistogramBinsCount = 20
	}
	if options.MaxDataPoints <= 0 {
		options.MaxDataPoints = 10000
	}
}

// enrichAnalysisData дополняет базовые данные расширенными вычислениями
func (s *RRIntervalsService) enrichAnalysisData(ctx context.Context, analysisData *models.CompleteAnalysisData, options *models.CompleteAnalysisOptions) error {
	// Если у нас есть сырые данные, вычисляем дополнительные метрики
	if len(analysisData.RawValues) > 0 {
		// HRV метрики
		if err := s.calculateHRVMetricsFromRawData(analysisData, options); err != nil {
			return err
		}

		// Анализ трендов
		if err := s.calculateTrendAnalysisFromAggregated(analysisData, options); err != nil {
			return err
		}

		// Гистограммы
		if err := s.calculateHistogramsFromRawData(analysisData, options); err != nil {
			return err
		}

		// Скаттерограмма
		if err := s.calculateScatterplotFromRawData(analysisData, options); err != nil {
			return err
		}
	}

	return nil
}

// calculateHRVMetricsFromRawData вычисляет HRV метрики из сырых данных
func (s *RRIntervalsService) calculateHRVMetricsFromRawData(analysisData *models.CompleteAnalysisData, options *models.CompleteAnalysisOptions) error {
	if len(analysisData.RawValues) < 10 {
		// Недостаточно данных для HRV анализа
		analysisData.HRVMetrics = &models.HRVMetrics{}
		return nil
	}

	// Используем существующие методы для вычисления HRV
	rmssd := s.calculateRMSSD(analysisData.RawValues)
	sdnn := s.calculateSDNN(analysisData.RawValues)
	pnn50 := s.calculatePNN50(analysisData.RawValues)
	triangularIndex := s.calculateTriangularIndex(analysisData.RawValues)
	tinn := s.calculateTINN(analysisData.RawValues)

	analysisData.HRVMetrics = &models.HRVMetrics{
		RMSSD:           rmssd,
		SDNN:            sdnn,
		PNN50:           pnn50,
		TriangularIndex: triangularIndex,
		TINN:            tinn,
		// Частотный анализ пока не реализован для производительности
		VLFPower:   0,
		LFPower:    0,
		HFPower:    0,
		LFHFRatio:  0,
		TotalPower: 0,
	}

	return nil
}

// calculateTrendAnalysisFromAggregated вычисляет анализ трендов из агрегированных данных
func (s *RRIntervalsService) calculateTrendAnalysisFromAggregated(analysisData *models.CompleteAnalysisData, options *models.CompleteAnalysisOptions) error {
	if len(analysisData.AggregatedData) < 3 {
		// Недостаточно точек для анализа трендов
		analysisData.TrendAnalysis = &models.RRTrendAnalysis{
			Period:        "insufficient_data",
			TrendPoints:   []models.TrendPoint{},
			OverallTrend:  models.OverallTrendStable,
			Correlation:   0,
			Seasonality:   []float64{},
			TrendStrength: 0,
		}
		return nil
	}

	// Создаем тренд-точки из агрегированных данных
	trendPoints := make([]models.TrendPoint, len(analysisData.AggregatedData))
	var prevValue float64

	for i, aggregated := range analysisData.AggregatedData {
		direction := models.TrendDirectionStable
		if i > 0 {
			if aggregated.Mean > prevValue {
				direction = models.TrendDirectionUp
			} else if aggregated.Mean < prevValue {
				direction = models.TrendDirectionDown
			}
		}

		trendPoints[i] = models.TrendPoint{
			Time:      aggregated.Time,
			Value:     aggregated.Mean,
			Direction: direction,
		}
		prevValue = aggregated.Mean
	}

	// Вычисляем корреляцию для определения общего тренда
	correlation := s.calculateCorrelationFromAggregated(analysisData.AggregatedData)
	overallTrend := s.determineOverallTrend(correlation)

	analysisData.TrendAnalysis = &models.RRTrendAnalysis{
		Period:        s.formatPeriod(analysisData.TimeRange.From, analysisData.TimeRange.To),
		TrendPoints:   trendPoints,
		OverallTrend:  overallTrend,
		Correlation:   correlation,
		Seasonality:   []float64{}, // TODO: реализовать сезонный анализ
		TrendStrength: abs(correlation),
	}

	return nil
}

// calculateHistogramsFromRawData вычисляет гистограммы из сырых данных
func (s *RRIntervalsService) calculateHistogramsFromRawData(analysisData *models.CompleteAnalysisData, options *models.CompleteAnalysisOptions) error {
	if len(analysisData.RawValues) == 0 {
		analysisData.Histogram = &models.RRHistogramData{
			Bins:       []models.HistogramBin{},
			TotalCount: 0,
			BinWidth:   0,
			Statistics: analysisData.Statistics,
		}
		analysisData.DiffHistogram = &models.DifferentialHistogramData{
			Bins:       []models.DifferentialHistogramBin{},
			TotalCount: 0,
			BinWidth:   0,
			Statistics: &models.DifferentialStatistics{},
		}
		return nil
	}

	// Обычная гистограмма
	binsCount := options.HistogramBinsCount
	if binsCount <= 0 {
		binsCount = s.calculateOptimalBinsCount(analysisData.RawValues)
	}

	histogram := s.buildHistogramBins(analysisData.RawValues, binsCount, analysisData.Statistics)
	analysisData.Histogram = &models.RRHistogramData{
		Bins:       histogram,
		TotalCount: int64(len(analysisData.RawValues)),
		BinWidth:   s.calculateBinWidth(analysisData.RawValues, binsCount),
		Statistics: analysisData.Statistics,
	}

	// Дифференциальная гистограмма
	if len(analysisData.RawValues) >= 2 {
		differences := make([]int64, len(analysisData.RawValues)-1)
		for i := 1; i < len(analysisData.RawValues); i++ {
			differences[i-1] = analysisData.RawValues[i] - analysisData.RawValues[i-1]
		}

		diffBinsCount := options.DiffHistogramBinsCount
		if diffBinsCount <= 0 {
			diffBinsCount = s.calculateOptimalBinsCountForDifferences(differences)
		}

		diffStats := s.calculateDifferentialStatistics(differences)
		diffHistogram := s.buildDifferentialHistogramBins(differences, diffBinsCount, diffStats)

		analysisData.DiffHistogram = &models.DifferentialHistogramData{
			Bins:       diffHistogram,
			TotalCount: int64(len(differences)),
			BinWidth:   s.calculateDifferentialBinWidth(differences, diffBinsCount),
			Statistics: diffStats,
		}
	}

	return nil
}

// calculateScatterplotFromRawData вычисляет скаттерограмму из сырых данных
func (s *RRIntervalsService) calculateScatterplotFromRawData(analysisData *models.CompleteAnalysisData, options *models.CompleteAnalysisOptions) error {
	if len(analysisData.RawValues) < 2 {
		analysisData.Scatterplot = &models.ScatterplotData{
			Points:     []models.ScatterplotPoint{},
			TotalCount: 0,
			Statistics: &models.ScatterplotStatistics{},
			Ellipse:    &models.PoincarePlotEllipse{},
		}
		return nil
	}

	// Создаем точки скаттерограммы
	points := make([]models.ScatterplotPoint, len(analysisData.RawValues)-1)
	for i := 0; i < len(analysisData.RawValues)-1; i++ {
		points[i] = models.ScatterplotPoint{
			RRn:  analysisData.RawValues[i],
			RRn1: analysisData.RawValues[i+1],
		}
	}

	// Вычисляем статистику скаттерограммы
	scatterStats := s.calculateScatterplotStatistics(points)
	ellipse := s.calculatePoincarePlotEllipse(points, scatterStats)

	analysisData.Scatterplot = &models.ScatterplotData{
		Points:     points,
		TotalCount: int64(len(points)),
		Statistics: scatterStats,
		Ellipse:    ellipse,
	}

	return nil
}

// Вспомогательные функции

func (s *RRIntervalsService) calculateCorrelationFromAggregated(aggregatedData []models.AggregatedRRData) float64 {
	if len(aggregatedData) < 2 {
		return 0
	}

	// Простой алгоритм корреляции времени и значений
	n := float64(len(aggregatedData))
	var sumX, sumY, sumXY, sumX2, sumY2 float64

	for i, data := range aggregatedData {
		x := float64(i) // Временная позиция
		y := data.Mean

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}

	denominator := (n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY)
	if denominator == 0 {
		return 0
	}

	return (n*sumXY - sumX*sumY) / (denominator * 0.5)
}

func (s *RRIntervalsService) formatPeriod(from, to time.Time) string {
	duration := to.Sub(from)
	if duration < time.Hour {
		return "short_term"
	} else if duration < 24*time.Hour {
		return "daily"
	} else if duration < 7*24*time.Hour {
		return "weekly"
	}
	return "long_term"
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
} 