package rr_intervals_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// GetRawValuesForAnalysis получает только значения R-R интервалов для статистического анализа.
func (r *RRIntervalsRepo) GetRawValuesForAnalysis(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]int64, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> filter(fn: (r) => r._value >= 300 and r._value <= 2000)
		|> sort(columns: ["_time"])
		|> keep(columns: ["_value"])`,
		r.bucket,
		from.Format(time.RFC3339Nano),
		to.Format(time.RFC3339Nano),
		userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query RR intervals for analysis: %v", err)
	}
	defer result.Close()

	var values []int64
	for result.Next() {
		record := result.Record()

		// Извлекаем значение R-R интервала
		rrIntervalMs, ok := record.Value().(int64)
		if !ok {
			// Попробуем преобразовать из float64
			if floatVal, ok := record.Value().(float64); ok {
				rrIntervalMs = int64(floatVal)
			} else {
				continue
			}
		}

		// Дополнительная проверка на уровне Go (защита от edge cases)
		if rrIntervalMs >= 300 && rrIntervalMs <= 2000 {
			values = append(values, rrIntervalMs)
		}
	}

	if result.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading RR intervals for analysis: %v", result.Err())
	}

	if len(values) == 0 {
		return nil, apperrors.NoValidDataf("no valid RR intervals found for the specified time range")
	}

	return values, nil
}

// GetValidRRIntervalsByUserID получает только валидные R-R интервалы пользователя (отфильтрованные).
func (r *RRIntervalsRepo) GetValidRRIntervalsByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.RRInterval, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> filter(fn: (r) => r._value >= 300 and r._value <= 2000)
		|> sort(columns: ["_time"])`,
		r.bucket,
		from.Format(time.RFC3339Nano),
		to.Format(time.RFC3339Nano),
		userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query valid RR intervals: %v", err)
	}
	defer result.Close()

	var rrIntervals []models.RRInterval
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		// Извлекаем device_id из tags
		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		// Извлекаем значение R-R интервала
		rrIntervalMs, ok := record.Value().(int64)
		if !ok {
			// Попробуем преобразовать из float64
			if floatVal, ok := record.Value().(float64); ok {
				rrIntervalMs = int64(floatVal)
			} else {
				continue
			}
		}

		// Дополнительная проверка диапазона на уровне Go
		if rrIntervalMs >= 300 && rrIntervalMs <= 2000 {
			// Создаем модель R-R интервала
			rrInterval := models.NewRRInterval(id, userID, deviceID, rrIntervalMs, record.Time())
			rrIntervals = append(rrIntervals, rrInterval)
		}
	}

	if result.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading valid RR intervals result: %v", result.Err())
	}

	if len(rrIntervals) == 0 {
		return nil, apperrors.NoValidDataf("no valid RR intervals found for the specified time range")
	}

	return rrIntervals, nil
}

// GetTimeSeriesForTrends получает временные ряды для анализа трендов.
func (r *RRIntervalsRepo) GetTimeSeriesForTrends(ctx context.Context, userID uuid.UUID, from, to time.Time, windowSizeMinutes int) ([]models.TrendPoint, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	if windowSizeMinutes < 1 || windowSizeMinutes > 60 {
		return nil, apperrors.ParameterOutOfRangef("window size must be between 1 and 60 minutes, got %d", windowSizeMinutes)
	}

	// Получаем сырые данные и агрегируем их на уровне Go
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> filter(fn: (r) => r._value >= 300 and r._value <= 2000)
		|> sort(columns: ["_time"])`,
		r.bucket,
		from.Format(time.RFC3339Nano),
		to.Format(time.RFC3339Nano),
		userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	queryResult, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query RR intervals for trends: %v", err)
	}
	defer queryResult.Close()

	// Собираем сырые данные
	var rawData []models.TimeValue
	for queryResult.Next() {
		record := queryResult.Record()

		var rrValue int64
		if intVal, ok := record.Value().(int64); ok {
			rrValue = intVal
		} else if floatVal, ok := record.Value().(float64); ok {
			rrValue = int64(floatVal)
		} else {
			continue
		}

		rawData = append(rawData, models.TimeValue{
			Time:  record.Time(),
			Value: rrValue,
		})
	}

	if queryResult.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading RR intervals for trends: %v", queryResult.Err())
	}

	if len(rawData) < 2 {
		return nil, apperrors.InsufficientDataf("need at least 2 data points for trend analysis, got %d", len(rawData))
	}

	// Агрегируем данные по временным интервалам в Go
	intervalDuration := time.Duration(windowSizeMinutes) * time.Minute
	aggregatedData := make(map[time.Time][]int64)

	// Группируем данные по временным окнам
	for _, dataPoint := range rawData {
		// Округляем время к началу интервала
		intervalStart := dataPoint.Time.Truncate(intervalDuration)
		aggregatedData[intervalStart] = append(aggregatedData[intervalStart], dataPoint.Value)
	}

	// Преобразуем агрегированные данные в тренд-точки
	var trendPoints []models.TrendPoint
	var prevValue float64

	// Сортируем временные интервалы
	var sortedTimes []time.Time
	for t := range aggregatedData {
		sortedTimes = append(sortedTimes, t)
	}

	// Простая сортировка пузырьком
	for i := range len(sortedTimes) - 1 {
		for j := i + 1; j < len(sortedTimes); j++ {
			if sortedTimes[i].After(sortedTimes[j]) {
				sortedTimes[i], sortedTimes[j] = sortedTimes[j], sortedTimes[i]
			}
		}
	}

	for _, intervalTime := range sortedTimes {
		values := aggregatedData[intervalTime]
		if len(values) == 0 {
			continue
		}

		// Вычисляем среднее значение
		var sum int64
		for _, val := range values {
			sum += val
		}
		meanValue := float64(sum) / float64(len(values))

		// Определяем направление тренда
		direction := models.TrendDirectionStable
		if len(trendPoints) > 0 {
			diff := meanValue - prevValue
			threshold := 10.0 // 10 мс порог для значимого изменения R-R интервала

			if diff > threshold {
				direction = models.TrendDirectionUp
			} else if diff < -threshold {
				direction = models.TrendDirectionDown
			}
		}

		trendPoints = append(trendPoints, models.TrendPoint{
			Time:      intervalTime,
			Value:     meanValue,
			Direction: direction,
		})

		prevValue = meanValue
	}

	return trendPoints, nil
}
