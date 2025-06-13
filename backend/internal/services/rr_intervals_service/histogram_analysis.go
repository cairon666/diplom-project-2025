package rr_intervals_service

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// BuildHistogram строит гистограмму R-R интервалов.
func (s *RRIntervalsService) BuildHistogram(ctx context.Context, userID uuid.UUID, from, to time.Time, binsCount int) (*models.RRHistogramData, error) {
	// Получаем сырые данные для анализа
	values, err := s.rrIntervalsRepo.GetRawValuesForAnalysis(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get RR intervals for histogram: %v", err)
	}

	if len(values) == 0 {
		return &models.RRHistogramData{
			Bins:       []models.HistogramBin{},
			TotalCount: 0,
			BinWidth:   0,
			Statistics: &models.RRStatisticalSummary{},
		}, nil
	}

	// Применяем адаптивный алгоритм выбора количества bins
	if binsCount <= 0 {
		binsCount = s.calculateOptimalBinsCount(values)
	}

	// Вычисляем базовую статистику
	statistics := s.calculateBasicStatistics(values)

	// Строим гистограмму
	histogram := s.buildHistogramBins(values, binsCount, statistics)

	return &models.RRHistogramData{
		Bins:       histogram,
		TotalCount: int64(len(values)),
		BinWidth:   s.calculateBinWidth(values, binsCount),
		Statistics: statistics,
	}, nil
}

// BuildDifferentialHistogram строит дифференциальную гистограмму разностей между соседними R-R интервалами.
func (s *RRIntervalsService) BuildDifferentialHistogram(ctx context.Context, userID uuid.UUID, from, to time.Time, binsCount int) (*models.DifferentialHistogramData, error) {
	// Получаем сырые данные для анализа
	values, err := s.rrIntervalsRepo.GetRawValuesForAnalysis(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get RR intervals for differential histogram: %v", err)
	}

	if len(values) < 2 {
		return &models.DifferentialHistogramData{
			Bins:       []models.DifferentialHistogramBin{},
			TotalCount: 0,
			BinWidth:   0,
			Statistics: &models.DifferentialStatistics{},
		}, nil
	}

	// Вычисляем разности между соседними интервалами
	differences := make([]int64, len(values)-1)
	for i := 1; i < len(values); i++ {
		differences[i-1] = values[i] - values[i-1]
	}

	// Применяем адаптивный алгоритм выбора количества bins
	if binsCount <= 0 {
		binsCount = s.calculateOptimalBinsCountForDifferences(differences)
	}

	// Вычисляем статистику разностей
	statistics := s.calculateDifferentialStatistics(differences)

	// Строим гистограмму
	histogram := s.buildDifferentialHistogramBins(differences, binsCount, statistics)

	return &models.DifferentialHistogramData{
		Bins:       histogram,
		TotalCount: int64(len(differences)),
		BinWidth:   s.calculateDifferentialBinWidth(differences, binsCount),
		Statistics: statistics,
	}, nil
}

// Вспомогательные методы для обычной гистограммы

// calculateOptimalBinsCount рассчитывает оптимальное количество bins по правилу Стерджеса.
func (s *RRIntervalsService) calculateOptimalBinsCount(values []int64) int {
	n := len(values)
	// Правило Стерджеса: bins = 1 + log2(n)
	binsCount := int(1 + math.Log2(float64(n)))

	// Ограничиваем диапазон для R-R интервалов
	if binsCount < 15 {
		binsCount = 15
	}
	if binsCount > 30 {
		binsCount = 30
	}

	return binsCount
}

// calculateBasicStatistics вычисляет базовую статистику.
func (s *RRIntervalsService) calculateBasicStatistics(values []int64) *models.RRStatisticalSummary {
	if len(values) == 0 {
		return &models.RRStatisticalSummary{}
	}

	// Сортируем для нахождения мин/макс
	sorted := make([]int64, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	// Вычисляем среднее
	var sum int64
	for _, v := range values {
		sum += v
	}
	mean := float64(sum) / float64(len(values))

	// Вычисляем стандартное отклонение
	var variance float64
	for _, v := range values {
		diff := float64(v) - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev := math.Sqrt(variance)

	return &models.RRStatisticalSummary{
		Mean:   mean,
		StdDev: stdDev,
		Min:    sorted[0],
		Max:    sorted[len(sorted)-1],
		Count:  int64(len(values)),
	}
}

// buildHistogramBins строит bins гистограммы.
func (s *RRIntervalsService) buildHistogramBins(values []int64, binsCount int, stats *models.RRStatisticalSummary) []models.HistogramBin {
	if len(values) == 0 || binsCount <= 0 {
		return []models.HistogramBin{}
	}

	binWidth := (stats.Max - stats.Min) / int64(binsCount)
	if binWidth == 0 {
		binWidth = 1
	}

	bins := make([]models.HistogramBin, binsCount)

	// Инициализируем bins
	for i := range binsCount {
		bins[i] = models.HistogramBin{
			RangeStart: stats.Min + int64(i)*binWidth,
			RangeEnd:   stats.Min + int64(i+1)*binWidth,
			Count:      0,
			Frequency:  0,
		}
	}

	// Заполняем bins данными
	for _, value := range values {
		binIndex := int((value - stats.Min) / binWidth)
		if binIndex >= binsCount {
			binIndex = binsCount - 1
		}
		if binIndex < 0 {
			binIndex = 0
		}
		bins[binIndex].Count++
	}

	// Вычисляем частоты
	totalCount := float64(len(values))
	for i := range bins {
		bins[i].Frequency = float64(bins[i].Count) / totalCount
	}

	return bins
}

// calculateBinWidth вычисляет ширину бина.
func (s *RRIntervalsService) calculateBinWidth(values []int64, binsCount int) int64 {
	if len(values) == 0 || binsCount <= 0 {
		return 0
	}

	stats := s.calculateBasicStatistics(values)

	return (stats.Max - stats.Min) / int64(binsCount)
}

// Вспомогательные методы для дифференциальной гистограммы

// calculateOptimalBinsCountForDifferences рассчитывает оптимальное количество bins для разностей.
func (s *RRIntervalsService) calculateOptimalBinsCountForDifferences(differences []int64) int {
	n := len(differences)
	// Правило Стерджеса адаптированное для разностей R-R интервалов
	binsCount := int(1 + math.Log2(float64(n)))

	// Ограничиваем диапазон для разностей (обычно меньший диапазон чем для основных интервалов)
	if binsCount < 10 {
		binsCount = 10
	}
	if binsCount > 25 {
		binsCount = 25
	}

	return binsCount
}

// calculateDifferentialStatistics вычисляет статистику разностей.
func (s *RRIntervalsService) calculateDifferentialStatistics(differences []int64) *models.DifferentialStatistics {
	if len(differences) == 0 {
		return &models.DifferentialStatistics{}
	}

	// Сортируем для нахождения мин/макс
	sorted := make([]int64, len(differences))
	copy(sorted, differences)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	// Вычисляем среднее
	var sum int64
	for _, v := range differences {
		sum += v
	}
	mean := float64(sum) / float64(len(differences))

	// Вычисляем стандартное отклонение
	var variance float64
	for _, v := range differences {
		diff := float64(v) - mean
		variance += diff * diff
	}
	variance /= float64(len(differences))
	stdDev := math.Sqrt(variance)

	// Вычисляем RMSSD (для разностей это квадратный корень из среднего квадратов разностей)
	var sumSquares float64
	for _, v := range differences {
		diff := float64(v)
		sumSquares += diff * diff
	}
	rmssd := math.Sqrt(sumSquares / float64(len(differences)))

	return &models.DifferentialStatistics{
		Mean:   mean,
		StdDev: stdDev,
		Min:    sorted[0],
		Max:    sorted[len(sorted)-1],
		Count:  int64(len(differences)),
		RMSSD:  rmssd,
	}
}

// buildDifferentialHistogramBins строит bins дифференциальной гистограммы.
func (s *RRIntervalsService) buildDifferentialHistogramBins(differences []int64, binsCount int, stats *models.DifferentialStatistics) []models.DifferentialHistogramBin {
	if len(differences) == 0 || binsCount <= 0 {
		return []models.DifferentialHistogramBin{}
	}

	binWidth := (stats.Max - stats.Min) / int64(binsCount)
	if binWidth == 0 {
		binWidth = 1
	}

	bins := make([]models.DifferentialHistogramBin, binsCount)

	// Инициализируем bins
	for i := range binsCount {
		bins[i] = models.DifferentialHistogramBin{
			RangeStart: stats.Min + int64(i)*binWidth,
			RangeEnd:   stats.Min + int64(i+1)*binWidth,
			Count:      0,
			Frequency:  0,
		}
	}

	// Заполняем bins данными
	for _, value := range differences {
		binIndex := int((value - stats.Min) / binWidth)
		if binIndex >= binsCount {
			binIndex = binsCount - 1
		}
		if binIndex < 0 {
			binIndex = 0
		}
		bins[binIndex].Count++
	}

	// Вычисляем частоты
	totalCount := float64(len(differences))
	for i := range bins {
		bins[i].Frequency = float64(bins[i].Count) / totalCount
	}

	return bins
}

// calculateDifferentialBinWidth вычисляет ширину бина для дифференциальной гистограммы.
func (s *RRIntervalsService) calculateDifferentialBinWidth(differences []int64, binsCount int) int64 {
	if len(differences) == 0 || binsCount <= 0 {
		return 0
	}

	stats := s.calculateDifferentialStatistics(differences)

	return (stats.Max - stats.Min) / int64(binsCount)
}
