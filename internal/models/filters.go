package models

import (
	"time"

	"github.com/google/uuid"
)

type DateRangeFilter struct {
	UserID uuid.UUID
	From   time.Time
	To     time.Time
}
