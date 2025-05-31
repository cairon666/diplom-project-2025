package rr_intervals_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/rr_intervals_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBatchRRIntervalsRequestDTO представляет HTTP запрос на создание batch R-R интервалов
type CreateBatchRRIntervalsRequestDTO struct {
	DeviceID  string                       `json:"device_id" binding:"required"`
	Intervals []RRIntervalCreateRequestDTO `json:"intervals" binding:"required"`
}

// RRIntervalCreateRequestDTO представляет отдельный R-R интервал в batch запросе
type RRIntervalCreateRequestDTO struct {
	RRIntervalMs int64      `json:"rr_interval_ms" binding:"required,min=200,max=3000"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
}

// CreateBatchRRIntervalsResponseDTO представляет HTTP ответ на создание batch R-R интервалов
type CreateBatchRRIntervalsResponseDTO struct {
	ProcessedCount int                     `json:"processed_count"`
	ValidCount     int                     `json:"valid_count"`
	Intervals      []RRIntervalResponseDTO `json:"intervals"`
}

// RRIntervalResponseDTO представляет HTTP ответ с данными R-R интервала
type RRIntervalResponseDTO struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	DeviceID     string    `json:"device_id"`
	RRIntervalMs int64     `json:"rr_interval_ms"`
	BPM          int64     `json:"bpm"`
	CreatedAt    time.Time `json:"created_at"`
	IsValid      bool      `json:"is_valid"`
}

// fromUsecaseRRIntervalResponse конвертирует usecase ответ в HTTP DTO
func fromUsecaseRRIntervalResponse(data rr_intervals_usecase.RRIntervalResponseData) RRIntervalResponseDTO {
	return RRIntervalResponseDTO{
		ID:           data.ID,
		UserID:       data.UserID,
		DeviceID:     data.DeviceID,
		RRIntervalMs: data.RRIntervalMs,
		BPM:          data.BPM,
		CreatedAt:    data.CreatedAt,
		IsValid:      data.IsValid,
	}
}

func (r *RRIntervalsRouter) CreateBatchRRIntervals(c *gin.Context) {
	var reqDTO CreateBatchRRIntervalsRequestDTO
	if err := c.BindJSON(&reqDTO); err != nil {
		www.HandleError(c, err)
		return
	}

	// Парсим device_id
	deviceID, err := uuid.Parse(reqDTO.DeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid device_id format",
		})
		return
	}

	// Преобразуем DTO в usecase структуры
	var intervals []rr_intervals_usecase.RRIntervalCreateData
	for _, intervalDTO := range reqDTO.Intervals {
		interval := rr_intervals_usecase.RRIntervalCreateData{
			DeviceID:     deviceID,
			RRIntervalMs: intervalDTO.RRIntervalMs,
			Timestamp:    intervalDTO.Timestamp,
		}
		intervals = append(intervals, interval)
	}

	// Создаем usecase запрос
	usecaseReq := rr_intervals_usecase.NewCreateBatchRRIntervalsRequest(deviceID, intervals)

	// Выполняем usecase
	resp, err := r.rrIntervalsUsecase.CreateBatchRRIntervals(c.Request.Context(), usecaseReq)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Преобразуем usecase ответ в HTTP DTO
	var intervalDTOs []RRIntervalResponseDTO
	for _, interval := range resp.Intervals {
		intervalDTOs = append(intervalDTOs, fromUsecaseRRIntervalResponse(interval))
	}

	// Формируем ответ
	responseDTO := CreateBatchRRIntervalsResponseDTO{
		ProcessedCount: resp.ProcessedCount,
		ValidCount:     resp.ValidCount,
		Intervals:      intervalDTOs,
	}

	c.JSON(http.StatusCreated, responseDTO)
}
