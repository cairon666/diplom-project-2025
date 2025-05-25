package device_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type DeviceService interface {
	Create(ctx context.Context, device models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Device, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeviceUsecase struct {
	deviceService DeviceService
	logger        *logger.Logger
}

func NewDeviceUsecase(deviceService DeviceService, logger *logger.Logger) *DeviceUsecase {
	return &DeviceUsecase{
		deviceService: deviceService,
		logger:        logger,
	}
}
