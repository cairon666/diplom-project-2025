package rr_intervals_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

// GetRRIntervalsRequest представляет запрос на получение R-R интервалов
type GetRRIntervalsRequest struct {
	DeviceID *uuid.UUID
	From     time.Time
	To       time.Time
}

func NewGetRRIntervalsRequest(deviceID *uuid.UUID, from, to time.Time) GetRRIntervalsRequest {
	return GetRRIntervalsRequest{
		DeviceID: deviceID,
		From:     from,
		To:       to,
	}
}

// GetRRIntervalsResponse представляет ответ с R-R интервалами
type GetRRIntervalsResponse struct {
	Intervals   []RRIntervalResponseData
	TotalCount  int
	ValidCount  int
	TimeRange   TimeRange
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

func (uc *RRIntervalsUsecase) GetRRIntervals(ctx context.Context, req GetRRIntervalsRequest) (GetRRIntervalsResponse, error) {
	// Проверяем права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnRRIntervals) {
		return GetRRIntervalsResponse{}, apperrors.Forbidden()
	}

	// Валидируем временной диапазон
	if req.From.After(req.To) {
		return GetRRIntervalsResponse{}, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", req.From, req.To)
	}

	// Проверяем максимальный размер запрашиваемого периода (защита от слишком больших запросов)
	maxDuration := 7 * 24 * time.Hour // 7 дней
	if req.To.Sub(req.From) > maxDuration {
		return GetRRIntervalsResponse{}, apperrors.TimeRangeTooLargef("time range too large (max %v), got %v", maxDuration, req.To.Sub(req.From))
	}

	var rrIntervals []RRIntervalResponseData

	// Получаем данные в зависимости от того, указан ли конкретный девайс
	if req.DeviceID != nil {
		// Проверяем права на устройство
		device, err := uc.deviceService.GetByID(ctx, *req.DeviceID)
		if err != nil {
			uc.logger.Error("failed to get device for RR intervals query",
				logger.String("device_id", req.DeviceID.String()),
				logger.Error(err))
			return GetRRIntervalsResponse{}, apperrors.DeviceNotFoundf("device not found: %v", err)
		}

		if device.UserID != authClaims.UserID {
			uc.logger.Warn("user attempted to get RR intervals for device they don't own",
				logger.String("user_id", authClaims.UserID.String()),
				logger.String("device_id", req.DeviceID.String()),
				logger.String("device_owner", device.UserID.String()))
			return GetRRIntervalsResponse{}, apperrors.DeviceAccessDenied()
		}

		// Получаем интервалы для конкретного устройства
		rawIntervals, err := uc.rrIntervalsService.GetByDeviceID(ctx, *req.DeviceID, req.From, req.To)
		if err != nil {
			uc.logger.Error("failed to get RR intervals by device",
				logger.String("device_id", req.DeviceID.String()),
				logger.Error(err))
			return GetRRIntervalsResponse{}, apperrors.DataProcessingErrorf("failed to get RR intervals: %v", err)
		}

		for _, rr := range rawIntervals {
			rrIntervals = append(rrIntervals, FromRRInterval(rr))
		}
	} else {
		// Получаем интервалы для всех устройств пользователя
		rawIntervals, err := uc.rrIntervalsService.GetByUserID(ctx, authClaims.UserID, req.From, req.To)
		if err != nil {
			uc.logger.Error("failed to get RR intervals by user",
				logger.String("user_id", authClaims.UserID.String()),
				logger.Error(err))
			return GetRRIntervalsResponse{}, apperrors.DataProcessingErrorf("failed to get RR intervals: %v", err)
		}

		for _, rr := range rawIntervals {
			rrIntervals = append(rrIntervals, FromRRInterval(rr))
		}
	}

	// Подсчитываем статистику
	totalCount := len(rrIntervals)
	validCount := 0

	for _, interval := range rrIntervals {
		if interval.IsValid {
			validCount++
		}
	}

	return GetRRIntervalsResponse{
		Intervals:   rrIntervals,
		TotalCount:  totalCount,
		ValidCount:  validCount,
		TimeRange: TimeRange{
			From: req.From,
			To:   req.To,
		},
	}, nil
} 