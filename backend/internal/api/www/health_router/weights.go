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

// CreateWeightRequest представляет запрос на создание веса.
type CreateWeightRequest struct {
	ID        *string `json:"id,omitempty"`
	WeightKg  float64 `binding:"required,min=0"    json:"weight_kg"`
	DeviceID  *string `json:"device_id,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}

// CreateWeightsRequest представляет запрос на создание множественных весов.
type CreateWeightsRequest struct {
	Weights []CreateWeightRequest `binding:"required,min=1" json:"weights"`
}

// Response DTOs для router.
type WeightResponse struct {
	ID        string    `json:"id"`
	WeightKg  float64   `json:"weight_kg"`
	DeviceID  *string   `json:"device_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type WeightsResponse struct {
	Weights []WeightResponse `json:"weights"`
}

type DailyWeightResponse struct {
	Weights map[string]float64 `json:"weights"` // key: day timestamp
}

// Функции конвертации.
func convertToRouterWeight(weight models.Weight) WeightResponse {
	var deviceID *string
	if weight.DeviceID != uuid.Nil {
		deviceIDStr := weight.DeviceID.String()
		deviceID = &deviceIDStr
	}

	return WeightResponse{
		ID:        weight.ID.String(),
		WeightKg:  weight.WeightKg,
		DeviceID:  deviceID,
		CreatedAt: weight.CreatedAt,
	}
}

func convertToRouterWeights(weights []models.Weight) []WeightResponse {
	result := make([]WeightResponse, len(weights))
	for i, weight := range weights {
		result[i] = convertToRouterWeight(weight)
	}

	return result
}

// CreateWeight обрабатывает запрос на создание веса.
func (r *HealthRouter) CreateWeight(c *gin.Context) {
	var req CreateWeightRequest
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

	dto := health_usecase.NewCreateWeightRequest(id, req.WeightKg, deviceID, createdAt)
	_, err = r.healthUsecase.CreateWeight(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// CreateWeights обрабатывает запрос на создание множественных весов.
func (r *HealthRouter) CreateWeights(c *gin.Context) {
	var req CreateWeightsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	weights := make([]health_usecase.CreateWeightRequest, len(req.Weights))
	for i, weightReq := range req.Weights {
		// Парсим ID
		var id uuid.UUID
		var err error
		if weightReq.ID != nil {
			id, err = r.parseUUID(*weightReq.ID)
		} else {
			id = uuid.New()
		}
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}

		// Парсим DeviceID
		var deviceID *uuid.UUID
		if weightReq.DeviceID != nil {
			deviceID, err = r.parseOptionalUUID(*weightReq.DeviceID)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())

				return
			}
		}

		// Парсим CreatedAt
		var createdAt time.Time
		if weightReq.CreatedAt != nil {
			createdAt, err = time.Parse(time.RFC3339, *weightReq.CreatedAt)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())

				return
			}
		} else {
			createdAt = time.Now()
		}

		weights[i] = health_usecase.NewCreateWeightRequest(id, weightReq.WeightKg, deviceID, createdAt)
	}

	dto := health_usecase.NewCreateWeightsRequest(weights)
	_, err := r.healthUsecase.CreateWeights(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// GetWeights обрабатывает запрос на получение весов.
func (r *HealthRouter) GetWeights(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetWeightsRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetWeights(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := WeightsResponse{
		Weights: convertToRouterWeights(usecaseResp.Weights),
	}

	c.JSON(http.StatusOK, routerResp)
}

// GetDailyWeightAvg обрабатывает запрос на получение агрегированных весов по дням.
func (r *HealthRouter) GetDailyWeightAvg(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetDailyWeightAvgRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetDailyWeightAvg(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := DailyWeightResponse{
		Weights: convertTimeMapToStringMapFloat(usecaseResp.Weights),
	}

	c.JSON(http.StatusOK, routerResp)
}
