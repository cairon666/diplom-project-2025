package rr_intervals_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

// RRIntervalCreateData представляет данные для создания R-R интервала
type RRIntervalCreateData struct {
	DeviceID     uuid.UUID
	RRIntervalMs int64
	Timestamp    *time.Time
}

// RRIntervalResponseData представляет данные R-R интервала для ответа
type RRIntervalResponseData struct {
	ID           string
	UserID       string
	DeviceID     string
	RRIntervalMs int64
	BPM          int64
	CreatedAt    time.Time
	IsValid      bool
}

// ToRRInterval конвертирует данные в модель R-R интервала
func (data RRIntervalCreateData) ToRRInterval(userID uuid.UUID) models.RRInterval {
	timestamp := time.Now()
	if data.Timestamp != nil {
		timestamp = *data.Timestamp
	}

	return models.NewRRInterval(uuid.New(), userID, data.DeviceID, data.RRIntervalMs, timestamp)
}

// FromRRInterval создает данные ответа из модели R-R интервала
func FromRRInterval(rr models.RRInterval) RRIntervalResponseData {
	return RRIntervalResponseData{
		ID:           rr.ID.String(),
		UserID:       rr.UserID.String(),
		DeviceID:     rr.DeviceID.String(),
		RRIntervalMs: rr.RRIntervalMs,
		BPM:          rr.ToBPM(),
		CreatedAt:    rr.CreatedAt,
		IsValid:      rr.IsValid(),
	}
}

// CreateBatchRRIntervalsRequest представляет запрос на создание batch R-R интервалов
type CreateBatchRRIntervalsRequest struct {
	DeviceID  uuid.UUID
	Intervals []RRIntervalCreateData
}

func NewCreateBatchRRIntervalsRequest(deviceID uuid.UUID, intervals []RRIntervalCreateData) CreateBatchRRIntervalsRequest {
	return CreateBatchRRIntervalsRequest{
		DeviceID:  deviceID,
		Intervals: intervals,
	}
}

// CreateBatchRRIntervalsResponse представляет ответ на создание batch R-R интервалов
type CreateBatchRRIntervalsResponse struct {
	ProcessedCount int
	ValidCount     int
	Intervals      []RRIntervalResponseData
}

func (uc *RRIntervalsUsecase) CreateBatchRRIntervals(ctx context.Context, req CreateBatchRRIntervalsRequest) (CreateBatchRRIntervalsResponse, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnRRIntervals) {
		return CreateBatchRRIntervalsResponse{}, apperrors.Forbidden()
	}

	// Проверяем существование устройства и права на него
	device, err := uc.deviceService.GetByID(ctx, req.DeviceID)
	if err != nil {
		uc.logger.Error("failed to get device for RR intervals batch",
			logger.String("device_id", req.DeviceID.String()),
			logger.Error(err))
		return CreateBatchRRIntervalsResponse{}, apperrors.DeviceNotFoundf("device not found: %v", err)
	}

	// Проверяем, что устройство принадлежит пользователю
	if device.UserID != authClaims.UserID {
		uc.logger.Warn("user attempted to create RR intervals for device they don't own",
			logger.String("user_id", authClaims.UserID.String()),
			logger.String("device_id", req.DeviceID.String()),
			logger.String("device_owner", device.UserID.String()))
		return CreateBatchRRIntervalsResponse{}, apperrors.DeviceAccessDenied()
	}

	// Валидируем количество интервалов
	if len(req.Intervals) == 0 {
		return CreateBatchRRIntervalsResponse{}, apperrors.BatchEmpty()
	}

	// Проверяем размер батча (добавим ограничение)
	maxBatchSize := 1000
	if len(req.Intervals) > maxBatchSize {
		return CreateBatchRRIntervalsResponse{}, apperrors.BatchTooLargef("batch size exceeds maximum allowed (%d), got %d", maxBatchSize, len(req.Intervals))
	}

	// Преобразуем запросы в модели RR интервалов
	var rrIntervals []models.RRInterval
	var validIntervals []RRIntervalResponseData
	validCount := 0
	invalidCount := 0

	for _, intervalData := range req.Intervals {
		// Проверяем корректность DeviceID в запросе
		if intervalData.DeviceID != req.DeviceID {
			uc.logger.Warn("interval device_id doesn't match batch device_id",
				logger.String("batch_device_id", req.DeviceID.String()),
				logger.String("interval_device_id", intervalData.DeviceID.String()))
			invalidCount++
			continue // Пропускаем некорректные интервалы
		}

		// Проверяем валидность значения RR интервала
		if intervalData.RRIntervalMs < 300 || intervalData.RRIntervalMs > 2000 {
			uc.logger.Warn("invalid RR interval value",
				logger.Int64("rr_interval_ms", intervalData.RRIntervalMs))
			invalidCount++
			continue
		}

		// Создаем модель R-R интервала
		rrInterval := intervalData.ToRRInterval(authClaims.UserID)

		// Проверяем валидность интервала
		if rrInterval.IsValid() {
			validCount++
		}

		rrIntervals = append(rrIntervals, rrInterval)
		validIntervals = append(validIntervals, FromRRInterval(rrInterval))
	}

	// Если нет валидных интервалов, возвращаем ошибку
	if len(rrIntervals) == 0 {
		return CreateBatchRRIntervalsResponse{}, apperrors.NoValidDataf("no valid intervals to process, %d invalid intervals found", invalidCount)
	}

	// Сохраняем batch в InfluxDB
	if err := uc.rrIntervalsService.CreateBatch(ctx, rrIntervals); err != nil {
		uc.logger.Error("failed to create RR intervals batch",
			logger.String("user_id", authClaims.UserID.String()),
			logger.String("device_id", req.DeviceID.String()),
			logger.Int("intervals_count", len(rrIntervals)),
			logger.Error(err))
		return CreateBatchRRIntervalsResponse{}, apperrors.DataProcessingErrorf("failed to save RR intervals: %v", err)
	}

	// Логируем только успешные batch операции - это важно для мониторинга
	uc.logger.Info("RR intervals batch created",
		logger.String("user_id", authClaims.UserID.String()),
		logger.String("device_id", req.DeviceID.String()),
		logger.Int("processed", len(rrIntervals)),
		logger.Int("valid", validCount))

	return CreateBatchRRIntervalsResponse{
		ProcessedCount: len(rrIntervals),
		ValidCount:     validCount,
		Intervals:      validIntervals,
	}, nil
}
