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

// DTO для CreateTemperature

type CreateTemperatureRequest struct {
	ID                uuid.UUID
	TemperatureCelsius float64
	DeviceID          *uuid.UUID
	CreatedAt         time.Time
}

func NewCreateTemperatureRequest(id uuid.UUID, temperatureCelsius float64, deviceID *uuid.UUID, createdAt time.Time) CreateTemperatureRequest {
	return CreateTemperatureRequest{
		ID:                id,
		TemperatureCelsius: temperatureCelsius,
		DeviceID:          deviceID,
		CreatedAt:         createdAt,
	}
}

type CreateTemperatureResponse struct{}

// DTO для CreateTemperatures

type CreateTemperaturesRequest struct {
	Temperatures []CreateTemperatureRequest
}

func NewCreateTemperaturesRequest(temperatures []CreateTemperatureRequest) CreateTemperaturesRequest {
	return CreateTemperaturesRequest{
		Temperatures: temperatures,
	}
}

type CreateTemperaturesResponse struct{}

// DTO для GetTemperatures

type GetTemperaturesRequest struct {
	From time.Time
	To   time.Time
}

func NewGetTemperaturesRequest(from, to time.Time) GetTemperaturesRequest {
	return GetTemperaturesRequest{
		From: from,
		To:   to,
	}
}

type GetTemperaturesResponse struct {
	Temperatures []models.Temperature
}

// DTO для GetHourlyTemperatureAvg

type GetHourlyTemperatureAvgRequest struct {
	From time.Time
	To   time.Time
}

func NewGetHourlyTemperatureAvgRequest(from, to time.Time) GetHourlyTemperatureAvgRequest {
	return GetHourlyTemperatureAvgRequest{
		From: from,
		To:   to,
	}
}

type GetHourlyTemperatureAvgResponse struct {
	Temperatures map[time.Time]float64
}

// DTO для GetDailyTemperatureAvg

type GetDailyTemperatureAvgRequest struct {
	From time.Time
	To   time.Time
}

func NewGetDailyTemperatureAvgRequest(from, to time.Time) GetDailyTemperatureAvgRequest {
	return GetDailyTemperatureAvgRequest{
		From: from,
		To:   to,
	}
}

type GetDailyTemperatureAvgResponse struct {
	Temperatures map[time.Time]float64
}

// CreateTemperature создает запись о температуре
func (u *HealthUsecase) CreateTemperature(ctx context.Context, dto CreateTemperatureRequest) (CreateTemperatureResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnTemperatures) {
		return CreateTemperatureResponse{}, apperrors.Forbidden()
	}

	var deviceID uuid.UUID
	if dto.DeviceID != nil {
		deviceID = *dto.DeviceID
	}
	temperature := models.NewTemperature(dto.ID, authClaims.UserID, deviceID, dto.TemperatureCelsius, dto.CreatedAt)

	err := u.healthService.CreateTemperature(ctx, temperature)
	if err != nil {
		u.logger.Error("failed to create temperature", logger.Error(err))
		return CreateTemperatureResponse{}, err
	}

	return CreateTemperatureResponse{}, nil
}

// CreateTemperatures создает множественные записи о температуре
func (u *HealthUsecase) CreateTemperatures(ctx context.Context, dto CreateTemperaturesRequest) (CreateTemperaturesResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnTemperatures) {
		return CreateTemperaturesResponse{}, apperrors.Forbidden()
	}

	temperatures := make([]models.Temperature, len(dto.Temperatures))
	for i, tempReq := range dto.Temperatures {
		var deviceID uuid.UUID
		if tempReq.DeviceID != nil {
			deviceID = *tempReq.DeviceID
		}
		temperatures[i] = models.NewTemperature(tempReq.ID, authClaims.UserID, deviceID, tempReq.TemperatureCelsius, tempReq.CreatedAt)
	}

	err := u.healthService.CreateTemperatures(ctx, temperatures)
	if err != nil {
		u.logger.Error("failed to create temperatures", logger.Error(err))
		return CreateTemperaturesResponse{}, err
	}

	return CreateTemperaturesResponse{}, nil
}

// GetTemperatures получает сырые данные о температуре
func (u *HealthUsecase) GetTemperatures(ctx context.Context, dto GetTemperaturesRequest) (GetTemperaturesResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnTemperatures) {
		return GetTemperaturesResponse{}, apperrors.Forbidden()
	}

	temperatures, err := u.healthService.GetTemperatures(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get temperatures", logger.Error(err))
		return GetTemperaturesResponse{}, err
	}

	return GetTemperaturesResponse{Temperatures: temperatures}, nil
}

// GetHourlyTemperatureAvg получает среднюю температуру по часам
func (u *HealthUsecase) GetHourlyTemperatureAvg(ctx context.Context, dto GetHourlyTemperatureAvgRequest) (GetHourlyTemperatureAvgResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnTemperatures) {
		return GetHourlyTemperatureAvgResponse{}, apperrors.Forbidden()
	}

	temperatures, err := u.aggregationService.GetHourlyTemperatureAvg(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get hourly temperature avg", logger.Error(err))
		return GetHourlyTemperatureAvgResponse{}, err
	}

	return GetHourlyTemperatureAvgResponse{Temperatures: temperatures}, nil
}

// GetDailyTemperatureAvg получает среднюю температуру по дням
func (u *HealthUsecase) GetDailyTemperatureAvg(ctx context.Context, dto GetDailyTemperatureAvgRequest) (GetDailyTemperatureAvgResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnTemperatures) {
		return GetDailyTemperatureAvgResponse{}, apperrors.Forbidden()
	}

	temperatures, err := u.aggregationService.GetDailyTemperatureAvg(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get daily temperature avg", logger.Error(err))
		return GetDailyTemperatureAvgResponse{}, err
	}

	return GetDailyTemperatureAvgResponse{Temperatures: temperatures}, nil
} 