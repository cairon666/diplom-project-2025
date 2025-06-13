package rr_intervals_router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/rr_intervals_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRRIntervalsResponseDTO представляет HTTP ответ с R-R интервалами.
type GetRRIntervalsResponseDTO struct {
	Intervals  []RRIntervalResponseDTO `json:"intervals"`
	TotalCount int                     `json:"total_count"`
	ValidCount int                     `json:"valid_count"`
	TimeRange  TimeRangeDTO            `json:"time_range"`
}

func (r *RRIntervalsRouter) GetRRIntervals(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")
	deviceIDStr := c.Query("device_id")

	// Валидируем обязательные параметры
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to parameters are required (RFC3339 format)",
		})

		return
	}

	// Парсим время
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid from time format (expected RFC3339)",
		})

		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid to time format (expected RFC3339)",
		})

		return
	}

	// Парсим device_id если указан
	var deviceID *uuid.UUID
	if deviceIDStr != "" {
		parsedDeviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid device_id format",
			})

			return
		}
		deviceID = &parsedDeviceID
	}

	// Создаем usecase запрос
	usecaseReq := rr_intervals_usecase.NewGetRRIntervalsRequest(deviceID, from, to)

	// Выполняем usecase
	resp, err := r.rrIntervalsUsecase.GetRRIntervals(c.Request.Context(), usecaseReq)
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
	responseDTO := GetRRIntervalsResponseDTO{
		Intervals:  intervalDTOs,
		TotalCount: resp.TotalCount,
		ValidCount: resp.ValidCount,
		TimeRange: TimeRangeDTO{
			From: resp.TimeRange.From,
			To:   resp.TimeRange.To,
		},
	}

	c.JSON(http.StatusOK, responseDTO)
}
