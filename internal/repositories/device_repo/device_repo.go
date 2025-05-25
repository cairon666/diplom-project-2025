package device_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DeviceRepo struct {
	query *dbqueries.Queries
}

func NewDeviceRepo(query *dbqueries.Queries) *DeviceRepo {
	return &DeviceRepo{
		query: query,
	}
}

func (r *DeviceRepo) Create(ctx context.Context, device models.Device) error {
	_, err := r.query.CreateDevice(ctx, dbqueries.CreateDeviceParams{
		ID:         device.ID,
		UserID:     device.UserID,
		DeviceName: device.DeviceName,
		CreatedAt:  device.CreatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *DeviceRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Device, error) {
	dbDevice, err := r.query.GetDeviceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Device{}, apperrors.ErrNotFound
		}
		return models.Device{}, err
	}
	return mapToModel(dbDevice), nil
}

func (r *DeviceRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Device, error) {
	dbDevices, err := r.query.ListDevicesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}
	var result []models.Device
	for _, d := range dbDevices {
		result = append(result, mapToModel(d))
	}
	return result, nil
}

func (r *DeviceRepo) UpdateName(ctx context.Context, id uuid.UUID, name string) (models.Device, error) {
	dbDevice, err := r.query.UpdateDeviceName(ctx, dbqueries.UpdateDeviceNameParams{
		ID:         id,
		DeviceName: name,
	})
	if err != nil {
		return models.Device{}, err
	}
	return mapToModel(dbDevice), nil
}

func (r *DeviceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.query.DeleteDevice(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}
	return nil
}

// Вспомогательная функция маппинга
func mapToModel(d dbqueries.DEVICE) models.Device {
	return models.Device{
		ID:         d.ID,
		UserID:     d.UserID,
		DeviceName: d.DeviceName,
		CreatedAt:  d.CreatedAt,
	}
}
