package device_service

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type DeviceRepo interface {
	Create(ctx context.Context, device models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Device, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Device, error)
	UpdateName(ctx context.Context, id uuid.UUID, name string) (models.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeviceService struct {
	deviceRepo DeviceRepo
}

func NewDeviceService(deviceRepo DeviceRepo) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

func (s *DeviceService) Create(ctx context.Context, device models.Device) error {
	return s.deviceRepo.Create(ctx, device)
}

func (s *DeviceService) GetByID(ctx context.Context, id uuid.UUID) (models.Device, error) {
	return s.deviceRepo.GetByID(ctx, id)
}

func (s *DeviceService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Device, error) {
	return s.deviceRepo.ListByUserID(ctx, userID)
}

func (s *DeviceService) UpdateName(ctx context.Context, id uuid.UUID, name string) (models.Device, error) {
	return s.deviceRepo.UpdateName(ctx, id, name)
}

func (s *DeviceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deviceRepo.Delete(ctx, id)
}
