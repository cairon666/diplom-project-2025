package health_usecase

import (
	"context"
	"strings"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

// isEmptyRangeError проверяет, является ли ошибка связанной с пустым диапазоном времени
func isEmptyRangeError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "empty range") || 
		   strings.Contains(errMsg, "cannot query an empty range") ||
		   strings.Contains(errMsg, "invalid range")
}

// DTO для CreateStep

type CreateStepRequest struct {
	ID        uuid.UUID
	StepCount int64
	DeviceID  *uuid.UUID
	CreatedAt time.Time
}

func NewCreateStepRequest(id uuid.UUID, stepCount int64, deviceID *uuid.UUID, createdAt time.Time) CreateStepRequest {
	return CreateStepRequest{
		ID:        id,
		StepCount: stepCount,
		DeviceID:  deviceID,
		CreatedAt: createdAt,
	}
}

type CreateStepResponse struct{}

// DTO для CreateSteps

type CreateStepsRequest struct {
	Steps []CreateStepRequest
}

func NewCreateStepsRequest(steps []CreateStepRequest) CreateStepsRequest {
	return CreateStepsRequest{
		Steps: steps,
	}
}

type CreateStepsResponse struct{}

// DTO для GetSteps

type GetStepsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetStepsRequest(from, to time.Time) GetStepsRequest {
	return GetStepsRequest{
		From: from,
		To:   to,
	}
}

type GetStepsResponse struct {
	Steps []models.Step
}

// DTO для GetHourlySteps

type GetHourlyStepsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetHourlyStepsRequest(from, to time.Time) GetHourlyStepsRequest {
	return GetHourlyStepsRequest{
		From: from,
		To:   to,
	}
}

type GetHourlyStepsResponse struct {
	Steps map[time.Time]int64
}

// DTO для GetDailySteps

type GetDailyStepsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetDailyStepsRequest(from, to time.Time) GetDailyStepsRequest {
	return GetDailyStepsRequest{
		From: from,
		To:   to,
	}
}

type GetDailyStepsResponse struct {
	Steps map[time.Time]int64
}

// CreateStep создает запись о шагах
func (u *HealthUsecase) CreateStep(ctx context.Context, dto CreateStepRequest) (CreateStepResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnSteps) {
		return CreateStepResponse{}, apperrors.Forbidden()
	}

	var deviceID uuid.UUID
	if dto.DeviceID != nil {
		deviceID = *dto.DeviceID
	}
	step := models.NewStep(dto.ID, authClaims.UserID, deviceID, dto.StepCount, dto.CreatedAt)

	err := u.healthService.CreateStep(ctx, step)
	if err != nil {
		u.logger.Error("failed to create step", logger.Error(err))
		return CreateStepResponse{}, err
	}

	return CreateStepResponse{}, nil
}

// CreateSteps создает множественные записи о шагах
func (u *HealthUsecase) CreateSteps(ctx context.Context, dto CreateStepsRequest) (CreateStepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnSteps) {
		return CreateStepsResponse{}, apperrors.Forbidden()
	}

	steps := make([]models.Step, len(dto.Steps))
	for i, stepReq := range dto.Steps {
		var deviceID uuid.UUID
		if stepReq.DeviceID != nil {
			deviceID = *stepReq.DeviceID
		}
		steps[i] = models.NewStep(stepReq.ID, authClaims.UserID, deviceID, stepReq.StepCount, stepReq.CreatedAt)
	}

	err := u.healthService.CreateSteps(ctx, steps)
	if err != nil {
		u.logger.Error("failed to create steps", logger.Error(err))
		return CreateStepsResponse{}, err
	}

	return CreateStepsResponse{}, nil
}

// GetSteps получает сырые данные о шагах
func (u *HealthUsecase) GetSteps(ctx context.Context, dto GetStepsRequest) (GetStepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSteps) {
		return GetStepsResponse{}, apperrors.Forbidden()
	}

	steps, err := u.healthService.GetSteps(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get steps", logger.Error(err))
		return GetStepsResponse{}, err
	}

	return GetStepsResponse{Steps: steps}, nil
}

// GetHourlySteps получает агрегированные данные по шагам по часам
func (u *HealthUsecase) GetHourlySteps(ctx context.Context, dto GetHourlyStepsRequest) (GetHourlyStepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSteps) {
		return GetHourlyStepsResponse{}, apperrors.Forbidden()
	}

	steps, err := u.aggregationService.GetHourlySteps(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		// Логируем как debug если это проблема с пустым диапазоном
		if isEmptyRangeError(err) {
			u.logger.Debug("empty range for hourly steps query", 
				logger.String("user_id", authClaims.UserID.String()),
				logger.String("from", dto.From.Format("2006-01-02T15:04:05Z07:00")),
				logger.String("to", dto.To.Format("2006-01-02T15:04:05Z07:00")),
			)
			// Возвращаем пустой результат вместо ошибки
			return GetHourlyStepsResponse{Steps: make(map[time.Time]int64)}, nil
		}
		u.logger.Error("failed to get hourly steps", logger.Error(err))
		return GetHourlyStepsResponse{}, err
	}

	return GetHourlyStepsResponse{Steps: steps}, nil
}

// GetDailySteps получает агрегированные данные по шагам по дням
func (u *HealthUsecase) GetDailySteps(ctx context.Context, dto GetDailyStepsRequest) (GetDailyStepsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSteps) {
		return GetDailyStepsResponse{}, apperrors.Forbidden()
	}

	steps, err := u.aggregationService.GetDailySteps(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		// Логируем как debug если это проблема с пустым диапазоном
		if isEmptyRangeError(err) {
			u.logger.Debug("empty range for daily steps query", 
				logger.String("user_id", authClaims.UserID.String()),
				logger.String("from", dto.From.Format("2006-01-02T15:04:05Z07:00")),
				logger.String("to", dto.To.Format("2006-01-02T15:04:05Z07:00")),
			)
			// Возвращаем пустой результат вместо ошибки
			return GetDailyStepsResponse{Steps: make(map[time.Time]int64)}, nil
		}
		u.logger.Error("failed to get daily steps", logger.Error(err))
		return GetDailyStepsResponse{}, err
	}

	return GetDailyStepsResponse{Steps: steps}, nil
} 