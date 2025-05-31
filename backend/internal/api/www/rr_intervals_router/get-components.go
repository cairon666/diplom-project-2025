package rr_intervals_router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

// GetHistogramResponseDTO представляет HTTP ответ с гистограммой
type GetHistogramResponseDTO struct {
	Histogram *RRHistogramDataDTO `json:"histogram"`
	TimeRange TimeRangeDTO        `json:"time_range"`
}

// GetTrendsResponseDTO представляет HTTP ответ с анализом трендов
type GetTrendsResponseDTO struct {
	TrendAnalysis *RRTrendAnalysisDTO `json:"trend_analysis"`
	TimeRange     TimeRangeDTO        `json:"time_range"`
}

// GetHRVMetricsResponseDTO представляет HTTP ответ с HRV метриками
type GetHRVMetricsResponseDTO struct {
	HRVMetrics *HRVMetricsDTO `json:"hrv_metrics"`
	TimeRange  TimeRangeDTO   `json:"time_range"`
}

// GetHistogram возвращает гистограмму R-R интервалов
func (r *RRIntervalsRouter) GetHistogram(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")
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

	// Парсим количество bins
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

	// Выполняем usecase
	histogram, err := r.rrIntervalsUsecase.GetRRHistogram(c.Request.Context(), from, to, binsCount)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Формируем ответ
	responseDTO := GetHistogramResponseDTO{
		Histogram: toRRHistogramDataDTO(histogram),
		TimeRange: TimeRangeDTO{
			From: from,
			To:   to,
		},
	}

	c.JSON(http.StatusOK, responseDTO)
}

// GetTrends возвращает анализ трендов R-R интервалов
func (r *RRIntervalsRouter) GetTrends(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")
	windowSizeStr := c.Query("window_size")

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

	// Парсим размер окна (по умолчанию 5 минут)
	windowSize := 5
	if windowSizeStr != "" {
		windowSize, err = strconv.Atoi(windowSizeStr)
		if err != nil || windowSize < 1 || windowSize > 60 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid window_size parameter (expected 1-60 minutes)",
			})
			return
		}
	}

	// Выполняем usecase
	trends, err := r.rrIntervalsUsecase.GetRRTrends(c.Request.Context(), from, to, windowSize)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Формируем ответ
	responseDTO := GetTrendsResponseDTO{
		TrendAnalysis: toRRTrendAnalysisDTO(trends),
		TimeRange: TimeRangeDTO{
			From: from,
			To:   to,
		},
	}

	c.JSON(http.StatusOK, responseDTO)
}

// GetHRVMetrics возвращает метрики вариабельности сердечного ритма
func (r *RRIntervalsRouter) GetHRVMetrics(c *gin.Context) {
	// Извлекаем query параметры
	fromStr := c.Query("from")
	toStr := c.Query("to")

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

	// Выполняем usecase
	hrvMetrics, err := r.rrIntervalsUsecase.GetHRVMetrics(c.Request.Context(), from, to)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Формируем ответ
	responseDTO := GetHRVMetricsResponseDTO{
		HRVMetrics: toHRVMetricsDTO(hrvMetrics),
		TimeRange: TimeRangeDTO{
			From: from,
			To:   to,
		},
	}

	c.JSON(http.StatusOK, responseDTO)
}