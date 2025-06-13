package models

import (
	"time"

	"github.com/google/uuid"
)

// RRInterval представляет R-R интервал между сердечными сокращениями.
type RRInterval struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	DeviceID     uuid.UUID
	RRIntervalMs int64     // R-R интервал в миллисекундах
	CreatedAt    time.Time // Точное время измерения
}

func NewRRInterval(id, userID, deviceID uuid.UUID, rrIntervalMs int64, createdAt time.Time) RRInterval {
	return RRInterval{
		ID:           id,
		UserID:       userID,
		DeviceID:     deviceID,
		RRIntervalMs: rrIntervalMs,
		CreatedAt:    createdAt,
	}
}

// ToBPM вычисляет мгновенную частоту пульса из R-R интервала.
func (rr RRInterval) ToBPM() int64 {
	if rr.RRIntervalMs <= 0 {
		return 0
	}

	return 60000 / rr.RRIntervalMs
}

// IsValid проверяет, находится ли R-R интервал в физиологических пределах.
func (rr RRInterval) IsValid() bool {
	return rr.RRIntervalMs >= 300 && rr.RRIntervalMs <= 2000
}
