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

// DTO для CreateWeight

type CreateWeightRequest struct {
	ID       uuid.UUID
	WeightKg float64
	DeviceID *uuid.UUID
	CreatedAt time.Time
}

func NewCreateWeightRequest(id uuid.UUID, weightKg float64, deviceID *uuid.UUID, createdAt time.Time) CreateWeightRequest {
	return CreateWeightRequest{
		ID:       id,
		WeightKg: weightKg,
		DeviceID: deviceID,
		CreatedAt: createdAt,
	}
}

type CreateWeightResponse struct{}

// DTO для CreateWeights

type CreateWeightsRequest struct {
	Weights []CreateWeightRequest
}

func NewCreateWeightsRequest(weights []CreateWeightRequest) CreateWeightsRequest {
	return CreateWeightsRequest{
		Weights: weights,
	}
}

type CreateWeightsResponse struct{}

// DTO для GetWeights

type GetWeightsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetWeightsRequest(from, to time.Time) GetWeightsRequest {
	return GetWeightsRequest{
		From: from,
		To:   to,
	}
}

type GetWeightsResponse struct {
	Weights []models.Weight
}

// DTO для GetDailyWeightAvg

type GetDailyWeightAvgRequest struct {
	From time.Time
	To   time.Time
}

func NewGetDailyWeightAvgRequest(from, to time.Time) GetDailyWeightAvgRequest {
	return GetDailyWeightAvgRequest{
		From: from,
		To:   to,
	}
}

type GetDailyWeightAvgResponse struct {
	Weights map[time.Time]float64
}

// CreateWeight создает запись о весе
func (u *HealthUsecase) CreateWeight(ctx context.Context, dto CreateWeightRequest) (CreateWeightResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnWeights) {
		return CreateWeightResponse{}, apperrors.Forbidden()
	}

	var deviceID uuid.UUID
	if dto.DeviceID != nil {
		deviceID = *dto.DeviceID
	}
	weight := models.NewWeight(dto.ID, authClaims.UserID, deviceID, dto.WeightKg, dto.CreatedAt)

	err := u.healthService.CreateWeight(ctx, weight)
	if err != nil {
		u.logger.Error("failed to create weight", logger.Error(err))
		return CreateWeightResponse{}, err
	}

	return CreateWeightResponse{}, nil
}

// CreateWeights создает множественные записи о весе
func (u *HealthUsecase) CreateWeights(ctx context.Context, dto CreateWeightsRequest) (CreateWeightsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnWeights) {
		return CreateWeightsResponse{}, apperrors.Forbidden()
	}

	weights := make([]models.Weight, len(dto.Weights))
	for i, weightReq := range dto.Weights {
		var deviceID uuid.UUID
		if weightReq.DeviceID != nil {
			deviceID = *weightReq.DeviceID
		}
		weights[i] = models.NewWeight(weightReq.ID, authClaims.UserID, deviceID, weightReq.WeightKg, weightReq.CreatedAt)
	}

	err := u.healthService.CreateWeights(ctx, weights)
	if err != nil {
		u.logger.Error("failed to create weights", logger.Error(err))
		return CreateWeightsResponse{}, err
	}

	return CreateWeightsResponse{}, nil
}

// GetWeights получает сырые данные о весе
func (u *HealthUsecase) GetWeights(ctx context.Context, dto GetWeightsRequest) (GetWeightsResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnWeights) {
		return GetWeightsResponse{}, apperrors.Forbidden()
	}

	weights, err := u.healthService.GetWeights(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get weights", logger.Error(err))
		return GetWeightsResponse{}, err
	}

	return GetWeightsResponse{Weights: weights}, nil
}

// GetDailyWeightAvg получает средний вес по дням
func (u *HealthUsecase) GetDailyWeightAvg(ctx context.Context, dto GetDailyWeightAvgRequest) (GetDailyWeightAvgResponse, error) {
	// Проверяем авторизацию и права доступа
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnWeights) {
		return GetDailyWeightAvgResponse{}, apperrors.Forbidden()
	}

	weights, err := u.aggregationService.GetDailyWeightAvg(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		u.logger.Error("failed to get daily weight avg", logger.Error(err))
		return GetDailyWeightAvgResponse{}, err
	}

	return GetDailyWeightAvgResponse{Weights: weights}, nil
} 