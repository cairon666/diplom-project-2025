package device_usecase

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

type CreateDeviceRequest struct {
	DeviceName string
}

func NewCreateDeviceRequest(deviceName string) CreateDeviceRequest {
	return CreateDeviceRequest{
		DeviceName: deviceName,
	}
}

type CreateDeviceResponse struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	DeviceName string
	CreatedAt  time.Time
}

func (ud *DeviceUsecase) CreateDevice(ctx context.Context, dto CreateDeviceRequest) (CreateDeviceResponse, error) {
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.RegisterDevice) {
		return CreateDeviceResponse{}, apperrors.Forbidden()
	}

	id := uuid.New()
	device := models.NewDevice(id, authClaims.UserID, dto.DeviceName, time.Now())
	if err := ud.deviceService.Create(ctx, device); err != nil {
		ud.logger.Error("failed to create device", logger.Error(err))

		return CreateDeviceResponse{}, err
	}

	return CreateDeviceResponse{
		ID:         device.ID,
		UserID:     device.UserID,
		DeviceName: device.DeviceName,
		CreatedAt:  device.CreatedAt,
	}, nil
}
