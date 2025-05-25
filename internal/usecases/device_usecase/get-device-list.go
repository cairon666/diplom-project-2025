package device_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type DeviceListItem struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	DeviceName string
	CreatedAt  time.Time
}

type GetDeviceListResponse struct {
	Devices []DeviceListItem
}

func (du *DeviceUsecase) GetDeviceList(ctx context.Context) (GetDeviceListResponse, error) {
	authClaims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !authClaims.HasPermission(permission.ReadOwnDevices) {
		return GetDeviceListResponse{}, apperrors.ErrForbidden
	}

	devices, err := du.deviceService.ListByUserID(ctx, authClaims.UserID)
	if errors.Is(err, apperrors.ErrNotFound) {
		return GetDeviceListResponse{Devices: []DeviceListItem{}}, nil
	} else if err != nil {
		du.logger.Error("failed to get device list", logger.Error(err))
		return GetDeviceListResponse{}, err
	}

	listItems := make([]DeviceListItem, 0, len(devices))
	for _, device := range devices {
		listItems = append(listItems, DeviceListItem{
			ID:         device.ID,
			UserID:     device.UserID,
			DeviceName: device.DeviceName,
			CreatedAt:  device.CreatedAt,
		})
	}

	return GetDeviceListResponse{Devices: listItems}, nil
}
