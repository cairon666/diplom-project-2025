package rr_intervals_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// CreateBatch создает несколько R-R интервалов одновременно.
func (r *RRIntervalsRepo) CreateBatch(ctx context.Context, rrIntervals []models.RRInterval) error {
	if len(rrIntervals) == 0 {
		return apperrors.BatchEmpty()
	}

	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	// Записываем точки по одной в рамках одного batch
	for _, rrInterval := range rrIntervals {
		// Создаем теги для идентификации записи
		tags := map[string]string{
			"id":        rrInterval.ID.String(),
			"user_id":   rrInterval.UserID.String(),
			"device_id": rrInterval.DeviceID.String(),
		}

		// Создаем поля с данными измерения
		fields := map[string]interface{}{
			"rr_interval_ms": rrInterval.RRIntervalMs,
		}

		// Создаем точку с наносекундной точностью времени
		point := write.NewPoint("rr_intervals", tags, fields, rrInterval.CreatedAt)

		// Записываем точку
		if err := writeAPI.WritePoint(ctx, point); err != nil {
			return apperrors.DataProcessingErrorf("failed to write RR interval: %v", err)
		}
	}

	return nil
}

// GetByUserID получает R-R интервалы пользователя за указанный период.
func (r *RRIntervalsRepo) GetByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.RRInterval, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> sort(columns: ["_time"])`,
		r.bucket,
		from.Format(time.RFC3339Nano),
		to.Format(time.RFC3339Nano),
		userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query RR intervals: %v", err)
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

		// Создаем модель R-R интервала
		rrInterval := models.NewRRInterval(id, userID, deviceID, rrIntervalMs, record.Time())
		rrIntervals = append(rrIntervals, rrInterval)
	}

	if result.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading RR intervals result: %v", result.Err())
	}

	return rrIntervals, nil
}

// GetByDeviceID получает R-R интервалы устройства за указанный период.
func (r *RRIntervalsRepo) GetByDeviceID(ctx context.Context, deviceID uuid.UUID, from, to time.Time) ([]models.RRInterval, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "rr_intervals")
		|> filter(fn: (r) => r.device_id == "%s")
		|> filter(fn: (r) => r._field == "rr_interval_ms")
		|> sort(columns: ["_time"])`,
		r.bucket,
		from.Format(time.RFC3339Nano),
		to.Format(time.RFC3339Nano),
		deviceID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query RR intervals by device: %v", err)
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

		// Извлекаем user_id из tags
		userIDStr, ok := record.ValueByKey("user_id").(string)
		if !ok {
			continue
		}
		userID, err := uuid.Parse(userIDStr)
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

		// Создаем модель R-R интервала
		rrInterval := models.NewRRInterval(id, userID, deviceID, rrIntervalMs, record.Time())
		rrIntervals = append(rrIntervals, rrInterval)
	}

	if result.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading RR intervals by device result: %v", result.Err())
	}

	return rrIntervals, nil
}
