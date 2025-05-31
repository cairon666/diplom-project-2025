package models

import (
	"time"

	"github.com/google/uuid"
)

type Step struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	DeviceID  uuid.UUID
	StepCount int64
	CreatedAt time.Time
}

func NewStep(id, userID, deviceID uuid.UUID, stepCount int64, createdAt time.Time) Step {
	return Step{
		ID:        id,
		UserID:    userID,
		DeviceID:  deviceID,
		StepCount: stepCount,
		CreatedAt: createdAt,
	}
}

type HeartRate struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	DeviceID  uuid.UUID
	BPM       int64
	CreatedAt time.Time
}

func NewHeartRate(id, userID, deviceID uuid.UUID, bpm int64, createdAt time.Time) HeartRate {
	return HeartRate{
		ID:        id,
		UserID:    userID,
		DeviceID:  deviceID,
		BPM:       bpm,
		CreatedAt: createdAt,
	}
}

type Temperature struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	DeviceID           uuid.UUID
	TemperatureCelsius float64
	CreatedAt          time.Time
}

func NewTemperature(id, userID, deviceID uuid.UUID, temperatureCelsius float64, createdAt time.Time) Temperature {
	return Temperature{
		ID:                 id,
		UserID:             userID,
		DeviceID:           deviceID,
		TemperatureCelsius: temperatureCelsius,
		CreatedAt:          createdAt,
	}
}

type Weight struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	DeviceID  uuid.UUID
	WeightKg  float64
	CreatedAt time.Time
}

func NewWeight(id, userID, deviceID uuid.UUID, weightKg float64, createdAt time.Time) Weight {
	return Weight{
		ID:        id,
		UserID:    userID,
		DeviceID:  deviceID,
		WeightKg:  weightKg,
		CreatedAt: createdAt,
	}
}

type Sleep struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	DeviceID  uuid.UUID
	StartedAt time.Time
	EndedAt   time.Time
}

func NewSleep(id, userID, deviceID uuid.UUID, startedAt, endedAt time.Time) Sleep {
	return Sleep{
		ID:        id,
		UserID:    userID,
		DeviceID:  deviceID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}
}
