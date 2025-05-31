package rr_intervals_router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// GetAggregatedDataResponseDTO представляет HTTP ответ с агрегированными данными
type GetAggregatedDataResponseDTO struct {
	Data      []AggregatedRRDataDTO `json:"data"`
	TimeRange TimeRangeDTO          `json:"time_range"`
	Interval  int                   `json:"interval_minutes"`
}

// AggregatedRRDataDTO представляет агрегированные данные R-R интервалов для HTTP ответа
type AggregatedRRDataDTO struct {
	Time   time.Time `json:"time"`
	Mean   float64   `json:"mean"`
	StdDev float64   `json:"std_dev"`
	Min    int64     `json:"min"`
	Max    int64     `json:"max"`
	Count  int64     `json:"count"`
}

func (r *RRIntervalsRouter) GetAggregatedData(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")
	intervalStr := c.Query("interval")

	// Валидируем обязательные параметры
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to parameters are required (RFC3339 format)",
		})
		return
	}

	if intervalStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "interval parameter is required (minutes)",
		})
		return
	}

	// Парсим время
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
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

	// Парсим интервал
	intervalMinutes, err := strconv.Atoi(intervalStr)
	if err != nil || intervalMinutes < 1 || intervalMinutes > 60 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid interval parameter (expected 1-60 minutes)",
		})
		return
	}

	// Выполняем usecase
	data, err := r.rrIntervalsUsecase.GetAggregatedRRData(c.Request.Context(), from, to, intervalMinutes)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Преобразуем usecase ответ в HTTP DTO
	var dataDTOs []AggregatedRRDataDTO
	for _, item := range data {
		dataDTOs = append(dataDTOs, toAggregatedRRDataDTO(item))
	}

	responseDTO := GetAggregatedDataResponseDTO{
		Data: dataDTOs,
		TimeRange: TimeRangeDTO{
			From: from,
			To:   to,
		},
		Interval: intervalMinutes,
	}

	c.JSON(http.StatusOK, responseDTO)
}

// toAggregatedRRDataDTO конвертирует модель в DTO
func toAggregatedRRDataDTO(data models.AggregatedRRData) AggregatedRRDataDTO {
	return AggregatedRRDataDTO{
		Time:   data.Time,
		Mean:   data.Mean,
		StdDev: data.StdDev,
		Min:    data.Min,
		Max:    data.Max,
		Count:  data.Count,
	}
} 