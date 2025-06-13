package health_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/usecases/health_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSleepRequest представляет запрос на создание сна.
type CreateSleepRequest struct {
	ID        *string `json:"id,omitempty"`
	StartedAt string  `binding:"required"         json:"started_at"`
	EndedAt   string  `binding:"required"         json:"ended_at"`
	DeviceID  *string `json:"device_id,omitempty"`
}

// CreateSleepsRequest представляет запрос на создание множественных снов.
type CreateSleepsRequest struct {
	Sleeps []CreateSleepRequest `binding:"required,min=1" json:"sleeps"`
}

// Response DTOs для router.
type SleepResponse struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	DeviceID  *string   `json:"device_id,omitempty"`
}

type SleepsResponse struct {
	Sleeps []SleepResponse `json:"sleeps"`
}

type DailySleepResponse struct {
	Sleeps map[string]float64 `json:"sleeps"` // key: day timestamp, value: duration in hours
}

// Функции конвертации.
func convertToRouterSleep(sleep models.Sleep) SleepResponse {
	var deviceID *string
	if sleep.DeviceID != uuid.Nil {
		deviceIDStr := sleep.DeviceID.String()
		deviceID = &deviceIDStr
	}

	return SleepResponse{
		ID:        sleep.ID.String(),
		StartedAt: sleep.StartedAt,
		EndedAt:   sleep.EndedAt,
		DeviceID:  deviceID,
	}
}

func convertToRouterSleeps(sleeps []models.Sleep) []SleepResponse {
	result := make([]SleepResponse, len(sleeps))
	for i, sleep := range sleeps {
		result[i] = convertToRouterSleep(sleep)
	}

	return result
}

// CreateSleep обрабатывает запрос на создание сна.
func (r *HealthRouter) CreateSleep(c *gin.Context) {
	var req CreateSleepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	// Парсим ID
	var id uuid.UUID
	var err error
	if req.ID != nil {
		id, err = r.parseUUID(*req.ID)
	} else {
		id = uuid.New()
	}
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	// Парсим DeviceID
	var deviceID *uuid.UUID
	if req.DeviceID != nil {
		deviceID, err = r.parseOptionalUUID(*req.DeviceID)
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}
	}

	// Парсим StartedAt
	startedAt, err := time.Parse(time.RFC3339, req.StartedAt)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	// Парсим EndedAt
	endedAt, err := time.Parse(time.RFC3339, req.EndedAt)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewCreateSleepRequest(id, startedAt, endedAt, deviceID, time.Now())
	_, err = r.healthUsecase.CreateSleep(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// CreateSleeps обрабатывает запрос на создание множественных снов.
func (r *HealthRouter) CreateSleeps(c *gin.Context) {
	var req CreateSleepsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	sleeps := make([]health_usecase.CreateSleepRequest, len(req.Sleeps))
	for i, sleepReq := range req.Sleeps {
		// Парсим ID
		var id uuid.UUID
		var err error
		if sleepReq.ID != nil {
			id, err = r.parseUUID(*sleepReq.ID)
		} else {
			id = uuid.New()
		}
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}

		// Парсим DeviceID
		var deviceID *uuid.UUID
		if sleepReq.DeviceID != nil {
			deviceID, err = r.parseOptionalUUID(*sleepReq.DeviceID)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())

				return
			}
		}

		// Парсим StartedAt
		startedAt, err := time.Parse(time.RFC3339, sleepReq.StartedAt)
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}

		// Парсим EndedAt
		endedAt, err := time.Parse(time.RFC3339, sleepReq.EndedAt)
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}

		sleeps[i] = health_usecase.NewCreateSleepRequest(id, startedAt, endedAt, deviceID, time.Now())
	}

	dto := health_usecase.NewCreateSleepsRequest(sleeps)
	_, err := r.healthUsecase.CreateSleeps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// GetSleeps обрабатывает запрос на получение снов.
func (r *HealthRouter) GetSleeps(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetSleepsRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetSleeps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := SleepsResponse{
		Sleeps: convertToRouterSleeps(usecaseResp.Sleeps),
	}

	c.JSON(http.StatusOK, routerResp)
}

// GetDailySleepDuration обрабатывает запрос на получение агрегированной продолжительности сна по дням.
func (r *HealthRouter) GetDailySleepDuration(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetDailySleepDurationRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetDailySleepDuration(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := DailySleepResponse{
		Sleeps: convertTimeMapToStringMapFloat(usecaseResp.SleepDurations),
	}

	c.JSON(http.StatusOK, routerResp)
}
