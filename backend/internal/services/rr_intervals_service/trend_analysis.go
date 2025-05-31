package rr_intervals_service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// AnalyzeTrends анализирует тренды R-R интервалов
func (s *RRIntervalsService) AnalyzeTrends(ctx context.Context, userID uuid.UUID, from, to time.Time, windowSizeMinutes int) (*models.RRTrendAnalysis, error) {
	// Получаем временные ряды для анализа
	timeSeries, err := s.rrIntervalsRepo.GetTimeSeriesForTrends(ctx, userID, from, to, windowSizeMinutes)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get time series for trends: %v", err)
	}

	if len(timeSeries) == 0 {
		return &models.RRTrendAnalysis{
			Period:        formatPeriod(from, to),
			TrendPoints:   []models.TrendPoint{},
			OverallTrend:  "stable",
			Correlation:   0,
			Seasonality:   []float64{},
			TrendStrength: 0,
		}, nil
	}

	// Анализируем полученные тренды
	correlation := s.calculateTrendCorrelation(timeSeries)
	overallTrend := s.determineOverallTrend(correlation)
	trendStrength := math.Abs(correlation)

	return &models.RRTrendAnalysis{
		Period:        formatPeriod(from, to),
		TrendPoints:   timeSeries,
		OverallTrend:  overallTrend,
		Correlation:   correlation,
		Seasonality:   s.calculateSeasonalityFromTrendPoints(timeSeries),
		TrendStrength: trendStrength,
	}, nil
}

// calculateTrendCorrelation вычисляет корреляцию тренда
func (s *RRIntervalsService) calculateTrendCorrelation(trendPoints []models.TrendPoint) float64 {
	if len(trendPoints) < 2 {
		return 0
	}

	n := float64(len(trendPoints))
	var sumX, sumY, sumXY, sumX2, sumY2 float64

	for i, point := range trendPoints {
		x := float64(i) // время как индекс
		y := point.Value
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}

	// Формула корреляции Пирсона
	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
	
	if denominator == 0 {
		return 0
	}
	
	return numerator / denominator
}

// determineOverallTrend определяет общий тренд
func (s *RRIntervalsService) determineOverallTrend(correlation float64) models.OverallTrend {
	threshold := 0.3 // Порог для значимой корреляции
	
	if correlation > threshold {
		return models.OverallTrendIncreasing
	} else if correlation < -threshold {
		return models.OverallTrendDecreasing
	}
	return models.OverallTrendStable
}

// calculateSeasonalityFromTrendPoints вычисляет сезонные компоненты из точек тренда
func (s *RRIntervalsService) calculateSeasonalityFromTrendPoints(trendPoints []models.TrendPoint) []float64 {
	// Упрощенный анализ сезонности - средние значения по часам
	hourlyMeans := make(map[int][]float64)
	
	for _, point := range trendPoints {
		hour := point.Time.Hour()
		hourlyMeans[hour] = append(hourlyMeans[hour], point.Value)
	}
	
	seasonality := make([]float64, 24)
	for hour := 0; hour < 24; hour++ {
		if values, exists := hourlyMeans[hour]; exists && len(values) > 0 {
			var sum float64
			for _, v := range values {
				sum += v
			}
			seasonality[hour] = sum / float64(len(values))
		}
	}
	
	return seasonality
}

// formatPeriod форматирует период для отображения
func formatPeriod(from, to time.Time) string {
	return fmt.Sprintf("%s - %s", from.Format("2006-01-02 15:04"), to.Format("2006-01-02 15:04"))
} 