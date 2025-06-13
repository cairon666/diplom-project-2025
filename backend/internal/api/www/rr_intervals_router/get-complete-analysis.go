package rr_intervals_router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/usecases/rr_intervals_usecase"
	"github.com/gin-gonic/gin"
)

// CompleteAnalysisDataDTO представляет полный комплексный анализ для HTTP ответа.
type CompleteAnalysisDataDTO struct {
	// Базовые данные
	RawValues []int64      `json:"raw_values,omitempty"`
	TimeRange TimeRangeDTO `json:"time_range"`

	// Статистические данные
	Statistics *RRStatisticalSummaryDTO `json:"statistics"`

	// HRV метрики
	HRVMetrics *HRVMetricsDTO `json:"hrv_metrics"`

	// Агрегированные данные по интервалам (используем существующий тип)
	AggregatedData []AggregatedRRDataDTO `json:"aggregated_data"`

	// Анализ трендов
	TrendAnalysis *RRTrendAnalysisDTO `json:"trend_analysis"`

	// Гистограммы
	Histogram     *RRHistogramDataDTO           `json:"histogram"`
	DiffHistogram *DifferentialHistogramDataDTO `json:"differential_histogram"`

	// Скаттерограмма
	Scatterplot *ScatterplotDataDTO `json:"scatterplot"`

	// Метаданные
	ProcessingTime string                 `json:"processing_time"` // Duration как строка
	DataQuality    *DataQualityMetricsDTO `json:"data_quality"`
}

// DataQualityMetricsDTO представляет метрики качества данных для HTTP ответа.
type DataQualityMetricsDTO struct {
	TotalMeasurements   int64   `json:"total_measurements"`
	ValidMeasurements   int64   `json:"valid_measurements"`
	InvalidMeasurements int64   `json:"invalid_measurements"`
	QualityPercentage   float64 `json:"quality_percentage"`
	MissingDataGaps     int64   `json:"missing_data_gaps"`
	LargestGapDuration  string  `json:"largest_gap_duration"` // Duration как строка
	AverageSamplingRate float64 `json:"average_sampling_rate"`
}

// GetCompleteAnalysisResponse представляет HTTP ответ комплексного анализа.
type GetCompleteAnalysisResponse struct {
	Data    *CompleteAnalysisDataDTO `json:"data"`
	Success bool                     `json:"success"`
	Message string                   `json:"message,omitempty"`
}

// GetCompleteAnalysis godoc
// @Summary Получение комплексного анализа RR интервалов
// @Description Выполняет полный комплексный анализ RR интервалов за один оптимизированный запрос. Заменяет множественные вызовы отдельных аналитических методов.
// @Tags RR Intervals Analytics
// @Accept json
// @Produce json
// @Param from query string true "Начало временного диапазона (RFC3339)" example(2024-01-01T10:00:00Z)
// @Param to query string true "Конец временного диапазона (RFC3339)" example(2024-01-01T11:00:00Z)
// @Param aggregation_interval query int false "Интервал агрегации в минутах (1-60)" default(5)
// @Param trend_window_size query int false "Размер окна для анализа трендов в минутах (5-120)" default(15)
// @Param histogram_bins query int false "Количество bins для гистограммы (5-100)" default(25)
// @Param diff_histogram_bins query int false "Количество bins для дифференциальной гистограммы (5-100)" default(20)
// @Param include_raw_data query bool false "Включать сырые данные в ответ" default(true)
// @Param max_data_points query int false "Максимальное количество точек данных (100-100000)" default(10000)
// @Param quick query bool false "Быстрый режим - без сырых данных и с упрощенными опциями" default(false)
// @Success 200 {object} GetCompleteAnalysisResponse
// @Failure 400 {object} www.ErrorResponse
// @Failure 401 {object} www.ErrorResponse
// @Failure 403 {object} www.ErrorResponse
// @Failure 500 {object} www.ErrorResponse
// @Router /v1/rr-intervals/analytics/complete [get].
func (r *RRIntervalsRouter) GetCompleteAnalysis(c *gin.Context) {
	// Парсинг временных параметров
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to parameters are required (RFC3339 format)",
		})

		return
	}

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

	// Парсинг опций анализа
	options, err := r.parseCompleteAnalysisOptions(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Формируем запрос
	req := rr_intervals_usecase.GetCompleteAnalysisRequest{
		From:    from,
		To:      to,
		Options: options,
	}

	// Выполняем анализ
	response, err := r.rrIntervalsUsecase.GetCompleteAnalysis(c.Request.Context(), req)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Преобразуем в DTO для ответа
	responseDTO := GetCompleteAnalysisResponse{
		Data:    toCompleteAnalysisDataDTO(response.Data),
		Success: true,
		Message: "Complete analysis performed successfully",
	}

	c.JSON(http.StatusOK, responseDTO)
}

// parseCompleteAnalysisOptions парсит опции комплексного анализа из query параметров.
func (r *RRIntervalsRouter) parseCompleteAnalysisOptions(c *gin.Context) (*models.CompleteAnalysisOptions, error) {
	// Проверяем быстрый режим
	quickMode := c.Query("quick") == "true"

	if quickMode {
		// Быстрый режим - оптимизированные опции для дашбордов
		return &models.CompleteAnalysisOptions{
			AggregationIntervalMinutes:    10, // Больший интервал
			TrendWindowSizeMinutes:        30, // Больше окно
			HistogramBinsCount:            15, // Меньше bins
			DiffHistogramBinsCount:        15, // Меньше bins
			EnableFrequencyDomainAnalysis: false,
			IncludeRawData:                false, // Без сырых данных
			MaxDataPoints:                 1000,  // Меньше точек
		}, nil
	}

	// Парсинг индивидуальных параметров
	options := &models.CompleteAnalysisOptions{
		AggregationIntervalMinutes:    5,  // По умолчанию
		TrendWindowSizeMinutes:        15, // По умолчанию
		HistogramBinsCount:            25, // По умолчанию
		DiffHistogramBinsCount:        20, // По умолчанию
		EnableFrequencyDomainAnalysis: false,
		IncludeRawData:                true,  // По умолчанию
		MaxDataPoints:                 10000, // По умолчанию
	}

	// Парсинг интервала агрегации
	if aggregationStr := c.Query("aggregation_interval"); aggregationStr != "" {
		if val, err := strconv.Atoi(aggregationStr); err == nil {
			options.AggregationIntervalMinutes = val
		} else {
			return nil, fmt.Errorf("invalid aggregation_interval parameter: %w", err)
		}
	}

	// Парсинг размера окна для трендов
	if trendWindowStr := c.Query("trend_window_size"); trendWindowStr != "" {
		if val, err := strconv.Atoi(trendWindowStr); err == nil {
			options.TrendWindowSizeMinutes = val
		} else {
			return nil, fmt.Errorf("invalid trend_window_size parameter: %w", err)
		}
	}

	// Парсинг количества bins для гистограммы
	if histogramBinsStr := c.Query("histogram_bins"); histogramBinsStr != "" {
		if val, err := strconv.Atoi(histogramBinsStr); err == nil {
			options.HistogramBinsCount = val
		} else {
			return nil, fmt.Errorf("invalid histogram_bins parameter: %w", err)
		}
	}

	// Парсинг количества bins для дифференциальной гистограммы
	if diffHistogramBinsStr := c.Query("diff_histogram_bins"); diffHistogramBinsStr != "" {
		if val, err := strconv.Atoi(diffHistogramBinsStr); err == nil {
			options.DiffHistogramBinsCount = val
		} else {
			return nil, fmt.Errorf("invalid diff_histogram_bins parameter: %w", err)
		}
	}

	// Парсинг флага включения сырых данных
	if includeRawStr := c.Query("include_raw_data"); includeRawStr != "" {
		options.IncludeRawData = includeRawStr == "true"
	}

	// Парсинг максимального количества точек данных
	if maxDataPointsStr := c.Query("max_data_points"); maxDataPointsStr != "" {
		if val, err := strconv.Atoi(maxDataPointsStr); err == nil {
			options.MaxDataPoints = val
		} else {
			return nil, fmt.Errorf("invalid max_data_points parameter: %w", err)
		}
	}

	return options, nil
}

// toCompleteAnalysisDataDTO преобразует CompleteAnalysisData в DTO.
func toCompleteAnalysisDataDTO(data *models.CompleteAnalysisData) *CompleteAnalysisDataDTO {
	if data == nil {
		return nil
	}

	dto := &CompleteAnalysisDataDTO{
		TimeRange:      TimeRangeDTO{From: data.TimeRange.From, To: data.TimeRange.To},
		Statistics:     toRRStatisticalSummaryDTO(data.Statistics),
		HRVMetrics:     toHRVMetricsDTO(data.HRVMetrics),
		TrendAnalysis:  toRRTrendAnalysisDTO(data.TrendAnalysis),
		Histogram:      toRRHistogramDataDTO(data.Histogram),
		DiffHistogram:  toDifferentialHistogramDataDTO(data.DiffHistogram),
		Scatterplot:    toScatterplotDataDTO(data.Scatterplot),
		ProcessingTime: data.ProcessingTime.String(),
		DataQuality:    toDataQualityMetricsDTO(data.DataQuality),
	}

	// Включаем сырые данные только если они есть
	if len(data.RawValues) > 0 {
		dto.RawValues = data.RawValues
	}

	// Преобразуем агрегированные данные используя существующую функцию
	for _, aggData := range data.AggregatedData {
		dto.AggregatedData = append(dto.AggregatedData, toAggregatedRRDataDTO(aggData))
	}

	return dto
}

// toDataQualityMetricsDTO преобразует DataQualityMetrics в DTO.
func toDataQualityMetricsDTO(metrics *models.DataQualityMetrics) *DataQualityMetricsDTO {
	if metrics == nil {
		return nil
	}

	return &DataQualityMetricsDTO{
		TotalMeasurements:   metrics.TotalMeasurements,
		ValidMeasurements:   metrics.ValidMeasurements,
		InvalidMeasurements: metrics.InvalidMeasurements,
		QualityPercentage:   metrics.QualityPercentage,
		MissingDataGaps:     metrics.MissingDataGaps,
		LargestGapDuration:  metrics.LargestGapDuration.String(),
		AverageSamplingRate: metrics.AverageSamplingRate,
	}
}
