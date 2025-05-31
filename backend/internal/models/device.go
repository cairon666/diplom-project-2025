package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	DeviceName string
	CreatedAt  time.Time
}

func NewDevice(ID uuid.UUID, userID uuid.UUID, deviceName string, createdAt time.Time) Device {
	return Device{
		ID:         ID,
		UserID:     userID,
		DeviceName: deviceName,
		CreatedAt:  createdAt,
	}
}
