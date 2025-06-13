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

// CreateTemperatureRequest представляет запрос на создание температуры.
type CreateTemperatureRequest struct {
	ID                 *string `json:"id,omitempty"`
	TemperatureCelsius float64 `binding:"required,min=30,max=50" json:"temperature_celsius"`
	DeviceID           *string `json:"device_id,omitempty"`
	CreatedAt          *string `json:"created_at,omitempty"`
}

// CreateTemperaturesRequest представляет запрос на создание множественных температур.
type CreateTemperaturesRequest struct {
	Temperatures []CreateTemperatureRequest `binding:"required,min=1" json:"temperatures"`
}

// Response DTOs для router.
type TemperatureResponse struct {
	ID                 string    `json:"id"`
	TemperatureCelsius float64   `json:"temperature_celsius"`
	DeviceID           *string   `json:"device_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type TemperaturesResponse struct {
	Temperatures []TemperatureResponse `json:"temperatures"`
}

type HourlyTemperatureResponse struct {
	Temperatures map[string]float64 `json:"temperatures"` // key: hour timestamp
}

type DailyTemperatureResponse struct {
	Temperatures map[string]float64 `json:"temperatures"` // key: day timestamp
}

// Функции конвертации.
func convertToRouterTemperature(temperature models.Temperature) TemperatureResponse {
	var deviceID *string
	if temperature.DeviceID != uuid.Nil {
		deviceIDStr := temperature.DeviceID.String()
		deviceID = &deviceIDStr
	}

	return TemperatureResponse{
		ID:                 temperature.ID.String(),
		TemperatureCelsius: temperature.TemperatureCelsius,
		DeviceID:           deviceID,
		CreatedAt:          temperature.CreatedAt,
	}
}

func convertToRouterTemperatures(temperatures []models.Temperature) []TemperatureResponse {
	result := make([]TemperatureResponse, len(temperatures))
	for i, temperature := range temperatures {
		result[i] = convertToRouterTemperature(temperature)
	}

	return result
}

// CreateTemperature обрабатывает запрос на создание температуры.
func (r *HealthRouter) CreateTemperature(c *gin.Context) {
	var req CreateTemperatureRequest
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

	dto := health_usecase.NewCreateTemperatureRequest(id, req.TemperatureCelsius, deviceID, createdAt)
	_, err = r.healthUsecase.CreateTemperature(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// CreateTemperatures обрабатывает запрос на создание множественных температур.
func (r *HealthRouter) CreateTemperatures(c *gin.Context) {
	var req CreateTemperaturesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	temperatures := make([]health_usecase.CreateTemperatureRequest, len(req.Temperatures))
	for i, tempReq := range req.Temperatures {
		// Парсим ID
		var id uuid.UUID
		var err error
		if tempReq.ID != nil {
			id, err = r.parseUUID(*tempReq.ID)
		} else {
			id = uuid.New()
		}
		if err != nil {
			www.HandleError(c, apperrors.InvalidParams())

			return
		}

		// Парсим DeviceID
		var deviceID *uuid.UUID
		if tempReq.DeviceID != nil {
			deviceID, err = r.parseOptionalUUID(*tempReq.DeviceID)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())

				return
			}
		}

		// Парсим CreatedAt
		var createdAt time.Time
		if tempReq.CreatedAt != nil {
			createdAt, err = time.Parse(time.RFC3339, *tempReq.CreatedAt)
			if err != nil {
				www.HandleError(c, apperrors.InvalidParams())

				return
			}
		} else {
			createdAt = time.Now()
		}

		temperatures[i] = health_usecase.NewCreateTemperatureRequest(id, tempReq.TemperatureCelsius, deviceID, createdAt)
	}

	dto := health_usecase.NewCreateTemperaturesRequest(temperatures)
	_, err := r.healthUsecase.CreateTemperatures(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusCreated)
}

// GetTemperatures обрабатывает запрос на получение температур.
func (r *HealthRouter) GetTemperatures(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetTemperaturesRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetTemperatures(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := TemperaturesResponse{
		Temperatures: convertToRouterTemperatures(usecaseResp.Temperatures),
	}

	c.JSON(http.StatusOK, routerResp)
}

// GetHourlyTemperatureAvg обрабатывает запрос на получение агрегированных температур по часам.
func (r *HealthRouter) GetHourlyTemperatureAvg(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetHourlyTemperatureAvgRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetHourlyTemperatureAvg(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := HourlyTemperatureResponse{
		Temperatures: convertTimeMapToStringMapFloat(usecaseResp.Temperatures),
	}

	c.JSON(http.StatusOK, routerResp)
}

// GetDailyTemperatureAvg обрабатывает запрос на получение агрегированных температур по дням.
func (r *HealthRouter) GetDailyTemperatureAvg(c *gin.Context) {
	from, to, err := r.parseDateRange(c)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())

		return
	}

	dto := health_usecase.NewGetDailyTemperatureAvgRequest(from, to)
	usecaseResp, err := r.healthUsecase.GetDailyTemperatureAvg(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Создаем router DTO
	routerResp := DailyTemperatureResponse{
		Temperatures: convertTimeMapToStringMapFloat(usecaseResp.Temperatures),
	}

	c.JSON(http.StatusOK, routerResp)
}
