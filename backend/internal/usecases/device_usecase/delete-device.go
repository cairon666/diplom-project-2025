package device_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type DeleteDeviceRequest struct {
	ID uuid.UUID
}

func NewDeleteDeviceRequest(id uuid.UUID) DeleteDeviceRequest {
	return DeleteDeviceRequest{
		ID: id,
	}
}

func (du *DeviceUsecase) DeleteDevice(ctx context.Context, req DeleteDeviceRequest) error {
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.RemoveDevice) {
		return apperrors.Forbidden()
	}

	device, err := du.deviceService.GetByID(ctx, req.ID)
	if err != nil {
		du.logger.Error("failed to get device", logger.Error(err))

		return err
	}

	if device.UserID != authClaims.UserID {
		return apperrors.Forbidden()
	}

	if err := du.deviceService.Delete(ctx, req.ID); err != nil {
		du.logger.Error("failed to delete device", logger.Error(err))

		return err
	}

	return nil
}
