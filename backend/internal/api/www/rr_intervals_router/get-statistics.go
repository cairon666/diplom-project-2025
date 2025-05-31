package rr_intervals_router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/rr_intervals_usecase"
	"github.com/gin-gonic/gin"
)

// GetRRStatisticsResponseDTO представляет HTTP ответ со статистикой R-R интервалов
type GetRRStatisticsResponseDTO struct {
	Summary    *RRStatisticalSummaryDTO `json:"summary"`
	Histogram  *RRHistogramDataDTO      `json:"histogram,omitempty"`
	HRVMetrics *HRVMetricsDTO           `json:"hrv_metrics,omitempty"`
	TimeRange  TimeRangeDTO             `json:"time_range"`
}

func (r *RRIntervalsRouter) GetRRStatistics(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")
	includeHistogramStr := c.Query("include_histogram")
	includeHRVStr := c.Query("include_hrv")
	binsCountStr := c.Query("bins_count")

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

	// Парсим опциональные параметры
	includeHistogram := false
	if includeHistogramStr != "" {
		includeHistogram, err = strconv.ParseBool(includeHistogramStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid include_histogram parameter (expected boolean)",
			})
			return
		}
	}

	includeHRV := false
	if includeHRVStr != "" {
		includeHRV, err = strconv.ParseBool(includeHRVStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid include_hrv parameter (expected boolean)",
			})
			return
		}
	}

	binsCount := 0
	if binsCountStr != "" {
		binsCount, err = strconv.Atoi(binsCountStr)
		if err != nil || binsCount < 0 || binsCount > 50 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid bins_count parameter (expected 0-50)",
			})
			return
		}
	}

	// Создаем usecase запрос
	usecaseReq := rr_intervals_usecase.NewGetRRStatisticsRequest(from, to)
	usecaseReq.IncludeHistogram = includeHistogram
	usecaseReq.IncludeHRV = includeHRV
	usecaseReq.BinsCount = binsCount

	// Выполняем usecase
	resp, err := r.rrIntervalsUsecase.GetRRStatistics(c.Request.Context(), usecaseReq)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Преобразуем usecase ответ в HTTP DTO
	responseDTO := GetRRStatisticsResponseDTO{
		Summary: toRRStatisticalSummaryDTO(resp.Summary),
		TimeRange: TimeRangeDTO{
			From: resp.TimeRange.From,
			To:   resp.TimeRange.To,
		},
	}

	// Добавляем опциональные компоненты
	if resp.Histogram != nil {
		responseDTO.Histogram = toRRHistogramDataDTO(resp.Histogram)
	}

	if resp.HRVMetrics != nil {
		responseDTO.HRVMetrics = toHRVMetricsDTO(resp.HRVMetrics)
	}

	c.JSON(http.StatusOK, responseDTO)
} 