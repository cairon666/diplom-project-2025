package rr_intervals_repo

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// GetStatisticalSummary получает базовую статистику R-R интервалов
func (r *RRIntervalsRepo) GetStatisticalSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.RRStatisticalSummary, error) {
	// Всегда используем сырые данные - просто и надежно
	return r.getStatisticsFromRawData(ctx, userID, from, to)
}


// getStatisticsFromRawData - fallback метод для исходных данных (медленно, но всегда работает)
func (r *RRIntervalsRepo) getStatisticsFromRawData(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.RRStatisticalSummary, error) {
	// Упрощенный запрос - сначала получаем все данные и считаем статистику в Go
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> filter(fn: (r) => r._value >= 300 and r._value <= 2000)
		|> sort(columns: ["_time"])
		|> keep(columns: ["_value"])
	`, r.bucket, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query raw RR intervals: %v", err)
	}
	defer result.Close()

	// Собираем все значения
	var values []float64
	for result.Next() {
		record := result.Record()
		
		// Извлекаем значение R-R интервала
		rrIntervalMs, ok := record.Value().(float64)
		if !ok {
			if intVal, ok := record.Value().(int64); ok {
				rrIntervalMs = float64(intVal)
			} else {
				continue
			}
		}
		
		values = append(values, rrIntervalMs)
	}

	if result.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading raw RR intervals: %v", result.Err())
	}

	// Вычисляем статистику в Go коде
	return r.calculateStatistics(values), nil
}

// calculateStatistics вычисляет статистику из массива значений R-R интервалов
func (r *RRIntervalsRepo) calculateStatistics(values []float64) *models.RRStatisticalSummary {
	if len(values) == 0 {
		return &models.RRStatisticalSummary{
			Mean:   0,
			StdDev: 0,
			Min:    0,
			Max:    0,
			Count:  0,
		}
	}

	// Подсчет количества
	count := int64(len(values))
	
	// Поиск минимума и максимума
	min := values[0]
	max := values[0]
	sum := 0.0
	
	for _, value := range values {
		sum += value
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	
	// Среднее значение
	mean := sum / float64(len(values))
	
	// Стандартное отклонение
	var stdDev float64
	if len(values) > 1 {
		// Для n > 1 можем вычислять стандартное отклонение
		variance := 0.0
		for _, value := range values {
			diff := value - mean
			variance += diff * diff
		}
		variance = variance / float64(len(values)-1) // Используем выборочную дисперсию (n-1)
		stdDev = math.Sqrt(variance)
	} else {
		// Для одного значения стандартное отклонение = 0
		stdDev = 0.0
	}

	return &models.RRStatisticalSummary{
		Mean:   mean,
		StdDev: stdDev,
		Min:    int64(min),
		Max:    int64(max),
		Count:  count,
	}
} 