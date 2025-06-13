package rr_intervals_service

import (
	"context"
	"math"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// CalculateHRVMetrics вычисляет метрики вариабельности сердечного ритма.
func (s *RRIntervalsService) CalculateHRVMetrics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.HRVMetrics, error) {
	// Проверяем минимальную продолжительность для HRV анализа
	minDuration := 5 * time.Minute
	if to.Sub(from) < minDuration {
		return nil, apperrors.TimeRangeTooSmallf("time range too small for HRV analysis (min %v), got %v", minDuration, to.Sub(from))
	}

	// Получаем сырые данные
	values, err := s.rrIntervalsRepo.GetRawValuesForAnalysis(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get RR intervals for HRV: %v", err)
	}

	if len(values) < 10 {
		return nil, apperrors.InsufficientDataf("need at least 10 valid RR intervals for HRV analysis, got %d", len(values))
	}

	// Вычисляем временные метрики HRV
	rmssd := s.calculateRMSSD(values)
	sdnn := s.calculateSDNN(values)
	pnn50 := s.calculatePNN50(values)
	triangularIndex := s.calculateTriangularIndex(values)

	return &models.HRVMetrics{
		RMSSD:           rmssd,
		SDNN:            sdnn,
		PNN50:           pnn50,
		TriangularIndex: triangularIndex,
		TINN:            s.calculateTINN(values),
		// Частотный анализ требует более сложной обработки
		VLFPower:   0, // TODO: Implement FFT analysis
		LFPower:    0,
		HFPower:    0,
		LFHFRatio:  0,
		TotalPower: 0,
	}, nil
}

// calculateRMSSD вычисляет RMSSD (Root Mean Square of Successive Differences).
func (s *RRIntervalsService) calculateRMSSD(values []int64) float64 {
	if len(values) < 2 {
		return 0
	}

	var sumSquares float64
	count := 0

	for i := 1; i < len(values); i++ {
		diff := float64(values[i] - values[i-1])
		sumSquares += diff * diff
		count++
	}

	if count == 0 {
		return 0
	}

	return math.Sqrt(sumSquares / float64(count))
}

// calculateSDNN вычисляет SDNN (Standard Deviation of NN intervals).
func (s *RRIntervalsService) calculateSDNN(values []int64) float64 {
	stats := s.calculateBasicStatistics(values)

	return stats.StdDev
}

// calculatePNN50 вычисляет pNN50 (процент соседних интервалов, различающихся более чем на 50 мс).
func (s *RRIntervalsService) calculatePNN50(values []int64) float64 {
	if len(values) < 2 {
		return 0
	}

	count50 := 0
	totalPairs := 0

	for i := 1; i < len(values); i++ {
		diff := values[i] - values[i-1]
		if diff < 0 {
			diff = -diff
		}

		if diff > 50 {
			count50++
		}
		totalPairs++
	}

	if totalPairs == 0 {
		return 0
	}

	return float64(count50) / float64(totalPairs) * 100
}

// calculateTriangularIndex вычисляет треугольный индекс.
func (s *RRIntervalsService) calculateTriangularIndex(values []int64) float64 {
	if len(values) < 10 {
		return 0
	}

	// Строим гистограмму с bin width = 7.8125 мс (стандарт)
	binWidth := int64(8) // Округляем до 8 мс
	stats := s.calculateBasicStatistics(values)

	binsCount := int((stats.Max-stats.Min)/binWidth) + 1
	bins := make([]int, binsCount)

	for _, value := range values {
		binIndex := int((value - stats.Min) / binWidth)
		if binIndex >= 0 && binIndex < binsCount {
			bins[binIndex]++
		}
	}

	// Находим максимальную высоту гистограммы
	maxHeight := 0
	for _, count := range bins {
		if count > maxHeight {
			maxHeight = count
		}
	}

	if maxHeight == 0 {
		return 0
	}

	return float64(len(values)) / float64(maxHeight)
}

// calculateTINN вычисляет TINN (Triangular Interpolation of NN histogram).
func (s *RRIntervalsService) calculateTINN(values []int64) float64 {
	// Упрощенная реализация TINN
	// В полной реализации требуется триангулярная интерполяция
	stats := s.calculateBasicStatistics(values)

	return float64(stats.Max - stats.Min)
}
