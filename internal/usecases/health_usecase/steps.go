package health_usecase

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
)

type CreateStepsRequest struct {
	Steps []models.Step
}

func NewCreateStepsRequest(steps []models.Step) CreateStepsRequest {
	return CreateStepsRequest{
		Steps: steps,
	}
}

func (hu *HealthUsecase) CreateSteps(ctx context.Context, dto CreateStepsRequest) error {
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.WriteOwnSteps) {
		return apperrors.ErrForbidden
	}

	if err := hu.healthService.CreateSteps(ctx, dto.Steps); err != nil {
		return err
	}

	return nil
}

type GetStepsRequest struct {
	From time.Time
	To   time.Time
}

func NewGetStepsRequest(from time.Time, to time.Time) GetStepsRequest {
	return GetStepsRequest{
		From: from,
		To:   to,
	}
}

type GetStepsResponse struct {
	Steps []models.Step
}

func (hu *HealthUsecase) GetSteps(ctx context.Context, dto GetStepsRequest) (GetStepsResponse, error) {
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnSteps) {
		return GetStepsResponse{}, apperrors.ErrForbidden
	}

	steps, err := hu.healthService.GetSteps(ctx, authClaims.UserID, dto.From, dto.To)
	if err != nil {
		hu.logger.Error("failed to get steps", logger.Error(err))
		return GetStepsResponse{}, err
	}

	return GetStepsResponse{Steps: steps}, nil
}
