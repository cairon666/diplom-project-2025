package models

import (
	"time"
)

// HRVMetrics представляет метрики вариабельности сердечного ритма.
type HRVMetrics struct {
	RMSSD           float64 // Root Mean Square of Successive Differences
	SDNN            float64 // Standard Deviation of NN intervals
	PNN50           float64 // Percentage of NN50 intervals
	TriangularIndex float64 // Triangular Index
	TINN            float64 // Triangular Interpolation of NN histogram
	// Частотный анализ
	VLFPower   float64 // Very Low Frequency power (0.003-0.04 Hz)
	LFPower    float64 // Low Frequency power (0.04-0.15 Hz)
	HFPower    float64 // High Frequency power (0.15-0.4 Hz)
	LFHFRatio  float64 // LF/HF ratio
	TotalPower float64 // Total spectral power
}

// TimeValue представляет значение с временной меткой для анализа трендов.
type TimeValue struct {
	Time  time.Time
	Value int64
}

// AggregatedRRData представляет агрегированные данные R-R интервалов по временным окнам.
type AggregatedRRData struct {
	Time   time.Time
	Mean   float64
	StdDev float64
	Min    int64
	Max    int64
	Count  int64
}

// RRStatisticalSummary представляет базовую статистику R-R интервалов.
type RRStatisticalSummary struct {
	Mean   float64
	StdDev float64
	Min    int64
	Max    int64
	Count  int64
}

// HistogramBin представляет один бин гистограммы.
type HistogramBin struct {
	RangeStart int64   // Начало диапазона в мс
	RangeEnd   int64   // Конец диапазона в мс
	Count      int64   // Количество значений в этом бине
	Frequency  float64 // Частота (Count / TotalCount)
}

// RRHistogramData представляет данные гистограммы R-R интервалов.
type RRHistogramData struct {
	Bins       []HistogramBin
	TotalCount int64
	BinWidth   int64 // Ширина бина в мс
	Statistics *RRStatisticalSummary
}

// TrendPoint представляет точку тренда.
type TrendPoint struct {
	Time      time.Time
	Value     float64
	Direction TrendDirection // "up", "down", "stable"
}

// RRTrendAnalysis представляет результат анализа трендов.
type RRTrendAnalysis struct {
	Period        string
	TrendPoints   []TrendPoint
	OverallTrend  OverallTrend // "increasing", "decreasing", "stable"
	Correlation   float64      // Коэффициент корреляции (-1 до 1)
	Seasonality   []float64    // Сезонные компоненты
	TrendStrength float64      // Сила тренда (0-1)
}

// DifferentialHistogramBin представляет один бин дифференциальной гистограммы.
type DifferentialHistogramBin struct {
	RangeStart int64   // Начало диапазона разности в мс
	RangeEnd   int64   // Конец диапазона разности в мс
	Count      int64   // Количество значений в этом бине
	Frequency  float64 // Частота (Count / TotalCount)
}

// DifferentialHistogramData представляет данные дифференциальной гистограммы.
type DifferentialHistogramData struct {
	Bins       []DifferentialHistogramBin
	TotalCount int64
	BinWidth   int64 // Ширина бина в мс
	Statistics *DifferentialStatistics
}

// DifferentialStatistics представляет статистику разностей между соседними RR интервалами.
type DifferentialStatistics struct {
	Mean   float64 // Среднее значение разностей
	StdDev float64 // Стандартное отклонение разностей
	Min    int64   // Минимальная разность
	Max    int64   // Максимальная разность
	Count  int64   // Количество разностей
	RMSSD  float64 // Root Mean Square of Successive Differences
}

// ScatterplotPoint представляет точку на скаттерограмме (диаграмме Пуанкаре).
type ScatterplotPoint struct {
	RRn  int64 // RR интервал n
	RRn1 int64 // RR интервал n+1
}

// ScatterplotData представляет данные скаттерограммы (диаграммы Пуанкаре).
type ScatterplotData struct {
	Points     []ScatterplotPoint
	TotalCount int64
	Statistics *ScatterplotStatistics
	Ellipse    *PoincarePlotEllipse
}

// ScatterplotStatistics представляет статистику скаттерограммы.
type ScatterplotStatistics struct {
	SD1         float64 // Стандартное отклонение по короткой оси эллипса
	SD2         float64 // Стандартное отклонение по длинной оси эллипса
	SD1SD2Ratio float64 // Отношение SD1/SD2
	CSI         float64 // Cardiac Sympathetic Index
	CVI         float64 // Cardiac Vagal Index
}

// PoincarePlotEllipse представляет параметры эллипса диаграммы Пуанкаре.
type PoincarePlotEllipse struct {
	CenterX float64 // Центр по X (среднее RR[n])
	CenterY float64 // Центр по Y (среднее RR[n+1])
	SD1     float64 // Полуось эллипса (короткая ось)
	SD2     float64 // Полуось эллипса (длинная ось)
	Area    float64 // Площадь эллипса
}

// CompleteAnalysisData представляет полный комплексный анализ RR интервалов
// Объединяет все виды анализа в одной структуре для оптимизированного получения данных.
type CompleteAnalysisData struct {
	// Базовые данные
	RawValues []int64   `json:"raw_values"`
	TimeRange TimeRange `json:"time_range"`

	// Статистические данные
	Statistics *RRStatisticalSummary `json:"statistics"`

	// HRV метрики
	HRVMetrics *HRVMetrics `json:"hrv_metrics"`

	// Агрегированные данные по интервалам
	AggregatedData []AggregatedRRData `json:"aggregated_data"`

	// Анализ трендов
	TrendAnalysis *RRTrendAnalysis `json:"trend_analysis"`

	// Гистограммы
	Histogram     *RRHistogramData           `json:"histogram"`
	DiffHistogram *DifferentialHistogramData `json:"differential_histogram"`

	// Скаттерограмма
	Scatterplot *ScatterplotData `json:"scatterplot"`

	// Метаданные
	ProcessingTime time.Duration       `json:"processing_time"`
	DataQuality    *DataQualityMetrics `json:"data_quality"`
}

// TimeRange представляет временной диапазон анализа.
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// DataQualityMetrics представляет метрики качества данных.
type DataQualityMetrics struct {
	TotalMeasurements   int64         `json:"total_measurements"`
	ValidMeasurements   int64         `json:"valid_measurements"`
	InvalidMeasurements int64         `json:"invalid_measurements"`
	QualityPercentage   float64       `json:"quality_percentage"`
	MissingDataGaps     int64         `json:"missing_data_gaps"`
	LargestGapDuration  time.Duration `json:"largest_gap_duration"`
	AverageSamplingRate float64       `json:"average_sampling_rate"` // Hz
}

// CompleteAnalysisOptions представляет опции для комплексного анализа.
type CompleteAnalysisOptions struct {
	// Опции для агрегации
	AggregationIntervalMinutes int `json:"aggregation_interval_minutes"`

	// Опции для анализа трендов
	TrendWindowSizeMinutes int `json:"trend_window_size_minutes"`

	// Опции для гистограмм
	HistogramBinsCount     int `json:"histogram_bins_count"`
	DiffHistogramBinsCount int `json:"diff_histogram_bins_count"`

	// Опции для HRV анализа
	EnableFrequencyDomainAnalysis bool `json:"enable_frequency_domain_analysis"`

	// Опции обработки
	IncludeRawData bool `json:"include_raw_data"`
	MaxDataPoints  int  `json:"max_data_points"` // Ограничение на количество точек
}
