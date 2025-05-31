package health_usecase

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

// DTO для CreateSleep

type CreateSleepRequest struct {
	ID        uuid.UUID
	StartedAt time.Time
	EndedAt   time.Time
	DeviceID  *uuid.UUID
	CreatedAt time.Time
}

func NewCreateSleepRequest(id uuid.UUID, startedAt, endedAt time.Time, deviceID *uuid.UUID, createdAt time.Time) CreateSleepRequest {
	return CreateSleepRequest{
		ID:        id,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		DeviceID:  deviceID,
		CreatedAt: createdAt,
	}
}

type CreateSleepResponse struct{}

// DTO для CreateSleeps

type CreateSleepsRequest struct {
	Sleeps []CreateSleepRequest
}

func NewCreateSleepsRequest(sleeps []CreateSleepRequest) CreateSleepsRequest {
	return CreateSleepsRequest{
		Sleeps: sleeps,
	}
}

type CreateSleepsResponse struct{}

// DTO для GetSleeps

type GetSleepsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetSleepsRequest(from, to time.Time) GetSleepsRequest {
	return GetSleepsRequest{
		From: from,
		To:   to,
	}
}

type GetSleepsResponse struct {
	Sleeps []models.Sleep
}

// DTO для GetDailySleepDuration

type GetDailySleepDurationRequest struct {
	From time.Time
	To   time.Time
}

func NewGetDailySleepDurationRequest(from, to time.Time) GetDailySleepDurationRequest {
	return GetDailySleepDurationRequest{
		From: from,
		To:   to,
	}
}

type GetDailySleepDurationResponse struct {
	SleepDurations map[time.Time]float64
}

// CreateSleep создает запись о сне
func (u *HealthUsecase) CreateSleep(ctx context.Context, dto CreateSleepRequest) (CreateSleepResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnSleeps) {
		return CreateSleepResponse{}, apperrors.Forbidden()
	}

	var deviceID uuid.UUID
	if dto.DeviceID != nil {
		deviceID = *dto.DeviceID
	}
	sleep := models.NewSleep(dto.ID, authClaims.UserID, deviceID, dto.StartedAt, dto.EndedAt)

	err := u.healthService.CreateSleep(ctx, sleep)
	if err != nil {
		u.logger.Error("failed to create sleep", logger.Error(err))
		return CreateSleepResponse{}, err
	}

	return CreateSleepResponse{}, nil
}

// CreateSleeps создает множественные записи о сне
func (u *HealthUsecase) CreateSleeps(ctx context.Context, dto CreateSleepsRequest) (CreateSleepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnSleeps) {
		return CreateSleepsResponse{}, apperrors.Forbidden()
	}

	sleeps := make([]models.Sleep, len(dto.Sleeps))
	for i, sleepReq := range dto.Sleeps {
		var deviceID uuid.UUID
		if sleepReq.DeviceID != nil {
			deviceID = *sleepReq.DeviceID
		}
		sleeps[i] = models.NewSleep(sleepReq.ID, authClaims.UserID, deviceID, sleepReq.StartedAt, sleepReq.EndedAt)
	}

	err := u.healthService.CreateSleeps(ctx, sleeps)
	if err != nil {
		u.logger.Error("failed to create sleeps", logger.Error(err))
		return CreateSleepsResponse{}, err
	}

	return CreateSleepsResponse{}, nil
}

// GetSleeps получает сырые данные о сне
func (u *HealthUsecase) GetSleeps(ctx context.Context, dto GetSleepsRequest) (GetSleepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSleeps) {
		return GetSleepsResponse{}, apperrors.Forbidden()
	}

	sleeps, err := u.healthService.GetSleeps(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get sleeps", logger.Error(err))
		return GetSleepsResponse{}, err
	}

	return GetSleepsResponse{Sleeps: sleeps}, nil
}

// GetDailySleepDuration получает продолжительность сна по дням
func (u *HealthUsecase) GetDailySleepDuration(ctx context.Context, dto GetDailySleepDurationRequest) (GetDailySleepDurationResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSleeps) {
		return GetDailySleepDurationResponse{}, apperrors.Forbidden()
	}

	sleepDurations, err := u.aggregationService.GetDailySleepDuration(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get daily sleep duration", logger.Error(err))
		return GetDailySleepDurationResponse{}, err
	}

	return GetDailySleepDurationResponse{SleepDurations: sleepDurations}, nil
} 