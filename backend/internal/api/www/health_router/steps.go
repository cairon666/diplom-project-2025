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

// CreateStepRequest представляет запрос на создание шага
type CreateStepRequest struct {
	ID        *string `json:"id,omitempty"`
	StepCount int64   `json:"step_count" binding:"required,min=0"`
	DeviceID  *string `json:"device_id,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}

// CreateStepsRequest представляет запрос на создание множественных шагов
type CreateStepsRequest struct {
	Steps []CreateStepRequest `json:"steps" binding:"required,min=1"`
}

// Response DTOs для router
type StepResponse struct {
	ID        string    `json:"id"`
	StepCount int64     `json:"step_count"`
	DeviceID  *string   `json:"device_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type StepsResponse struct {
	Steps []StepResponse `json:"steps"`
}

type HourlyStepsResponse struct {
	Steps map[string]int64 `json:"steps"` // key: hour timestamp
}

type DailyStepsResponse struct {
	Steps map[string]int64 `json:"steps"` // key: day timestamp
}

// Функции конвертации
func convertToRouterStep(step models.Step) StepResponse {
	var deviceID *string
	if step.DeviceID != uuid.Nil {
		deviceIDStr := step.DeviceID.String()
		deviceID = &deviceIDStr
	}
	
	return StepResponse{
		ID:        step.ID.String(),
		StepCount: step.StepCount,
		DeviceID:  deviceID,
		CreatedAt: step.CreatedAt,
	}
}

func convertToRouterSteps(steps []models.Step) []StepResponse {
	result := make([]StepResponse, len(steps))
	for i, step := range steps {
		result[i] = convertToRouterStep(step)
	}
	return result
}

func convertTimeMapToStringMap(timeMap map[time.Time]int64) map[string]int64 {
	result := make(map[string]int64)
	for t, value := range timeMap {
		result[t.Format(time.RFC3339)] = value
	}
	return result
}

// CreateStep обрабатывает запрос на создание шага
func (r *HealthRouter) CreateStep(c *gin.Context) {
	var req CreateStepRequest
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

	// Парсим CreatedAt
	var createdAt time.Time
	if req.CreatedAt != nil {
		createdAt, err = time.Parse(time.RFC3339, *req.CreatedAt)
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())
			return
		}
	} else {
		createdAt = time.Now()
	}

	dto := health_usecase.NewCreateStepRequest(id, req.StepCount, deviceID, createdAt)
	_, err = r.healthUsecase.CreateStep(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// CreateSteps обрабатывает запрос на создание множественных шагов
func (r *HealthRouter) CreateSteps(c *gin.Context) {
	var req CreateStepsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, apperrors.InvalidParams())
		return
	}

	steps := make([]health_usecase.CreateStepRequest, len(req.Steps))
	for i, stepReq := range req.Steps {
		// Парсим ID
		var id uuid.UUID
		var err error
		if stepReq.ID != nil {
			id, err = r.parseUUID(*stepReq.ID)
		} else {
			id = uuid.New()
		}
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())
			return
		}

		// Парсим DeviceID
		var deviceID *uuid.UUID
		if stepReq.DeviceID != nil {
			deviceID, err = r.parseOptionalUUID(*stepReq.DeviceID)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())
				return
			}
		}

		// Парсим CreatedAt
		var createdAt time.Time
		if stepReq.CreatedAt != nil {
			createdAt, err = time.Parse(time.RFC3339, *stepReq.CreatedAt)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())
				return
			}
		} else {
			createdAt = time.Now()
		}

		steps[i] = health_usecase.NewCreateStepRequest(id, stepReq.StepCount, deviceID, createdAt)
	}

	dto := health_usecase.NewCreateStepsRequest(steps)
	_, err := r.healthUsecase.CreateSteps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// GetSteps обрабатывает запрос на получение шагов
func (r *HealthRouter) GetSteps(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())
		return
	}

	dto := health_usecase.NewGetStepsRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetSteps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Создаем router DTO
	routerResp := StepsResponse{
		Steps: convertToRouterSteps(usecaseResp.Steps),
	}
	
	c.JSON(http.StatusOK, routerResp)
}

// GetHourlySteps обрабатывает запрос на получение агрегированных шагов по часам
func (r *HealthRouter) GetHourlySteps(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())
		return
	}

	dto := health_usecase.NewGetHourlyStepsRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetHourlySteps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Создаем router DTO
	routerResp := HourlyStepsResponse{
		Steps: convertTimeMapToStringMap(usecaseResp.Steps),
	}
	
	c.JSON(http.StatusOK, routerResp)
}

// GetDailySteps обрабатывает запрос на получение агрегированных шагов по дням
func (r *HealthRouter) GetDailySteps(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())
		return
	}

	dto := health_usecase.NewGetDailyStepsRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetDailySteps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Создаем router DTO
	routerResp := DailyStepsResponse{
		Steps: convertTimeMapToStringMap(usecaseResp.Steps),
	}
	
	c.JSON(http.StatusOK, routerResp)
} 