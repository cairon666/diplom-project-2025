package health_router

import (
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Вспомогательные методы

func (r *HealthRouter) parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	now := time.Now()

	// По умолчанию: последние 7 дней
	defaultFrom := now.AddDate(0, 0, -7)
	defaultTo := now

	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		from = defaultFrom
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {	
			return time.Time{}, time.Time{}, err
		}
	} else {
		to = defaultTo
	}

	if from.After(to) {
		return time.Time{}, time.Time{}, apperrors.InvalidParams()
	}

	return from, to, nil
}

func (r *HealthRouter) parseUUID(idStr string) (uuid.UUID, error) {
	if idStr == "" {
		return uuid.New(), nil // Генерируем новый UUID если не указан
	}

	return uuid.Parse(idStr)
}

func (r *HealthRouter) parseOptionalUUID(idStr string) (*uuid.UUID, error) {
	if idStr == "" {
		return nil, nil
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func convertTimeMapToStringMapFloat(timeMap map[time.Time]float64) map[string]float64 {
	result := make(map[string]float64)
	for t, value := range timeMap {
		result[t.Format(time.RFC3339)] = value
	}

	return result
}
