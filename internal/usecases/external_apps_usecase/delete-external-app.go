package external_apps_usecase

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type DeleteExternalAppRequest struct {
	ID uuid.UUID
}

func NewDeleteExternalAppRequest(id uuid.UUID) DeleteExternalAppRequest {
	return DeleteExternalAppRequest{
		ID: id,
	}
}

func (uc *ExternalAppsUsecase) DeleteExternalApp(ctx context.Context, dto DeleteExternalAppRequest) error {
	claims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !claims.HasPermission(permission.UpdateOwnExternalApps) {
		return apperrors.ErrForbidden
	}

	externalApp, err := uc.externalAppsService.GetByID(ctx, dto.ID)
	if err != nil {
		uc.logger.Error("failed to get external app", logger.Error(err))
		return err
	}

	if externalApp.OwnerID != claims.UserID {
		return apperrors.ErrForbidden
	}

	if err := uc.externalAppsService.Delete(ctx, dto.ID); errors.Is(err, apperrors.ErrNotFound) {
		return err
	} else if err != nil {
		uc.logger.Error("failed to delete external app", logger.Error(err))
		return err
	}

	return nil
}
