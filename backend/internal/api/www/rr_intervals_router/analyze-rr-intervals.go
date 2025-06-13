package rr_intervals_router

import (
	"time"

	"github.com/cairon666/vkr-backend/internal/models"
)

// RRHistogramDataDTO представляет данные гистограммы для HTTP ответа.
type RRHistogramDataDTO struct {
	Bins       []HistogramBinDTO        `json:"bins"`
	TotalCount int64                    `json:"total_count"`
	BinWidth   int64                    `json:"bin_width"`
	Statistics *RRStatisticalSummaryDTO `json:"statistics"`
}

// HistogramBinDTO представляет один бин гистограммы для HTTP ответа.
type HistogramBinDTO struct {
	RangeStart int64   `json:"range_start"`
	RangeEnd   int64   `json:"range_end"`
	Count      int64   `json:"count"`
	Frequency  float64 `json:"frequency"`
}

// RRTrendAnalysisDTO представляет результат анализа трендов для HTTP ответа.
type RRTrendAnalysisDTO struct {
	Period        string          `json:"period"`
	TrendPoints   []TrendPointDTO `json:"trend_points"`
	OverallTrend  string          `json:"overall_trend"`
	Correlation   float64         `json:"correlation"`
	Seasonality   []float64       `json:"seasonality"`
	TrendStrength float64         `json:"trend_strength"`
}

// TrendPointDTO represents a trend point for HTTP response.
type TrendPointDTO struct {
	Time      time.Time `json:"time"`
	Value     float64   `json:"value"`
	Direction string    `json:"direction"`
}

// RRStatisticalSummaryDTO represents basic statistics for HTTP response.
type RRStatisticalSummaryDTO struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    int64   `json:"min"`
	Max    int64   `json:"max"`
	Count  int64   `json:"count"`
}

// HRVMetricsDTO представляет метрики HRV для HTTP ответа.
type HRVMetricsDTO struct {
	RMSSD           float64 `json:"rmssd"`
	SDNN            float64 `json:"sdnn"`
	PNN50           float64 `json:"pnn50"`
	TriangularIndex float64 `json:"triangular_index"`
	TINN            float64 `json:"tinn"`
	VLFPower        float64 `json:"vlf_power"`
	LFPower         float64 `json:"lf_power"`
	HFPower         float64 `json:"hf_power"`
	LFHFRatio       float64 `json:"lf_hf_ratio"`
	TotalPower      float64 `json:"total_power"`
}

// DifferentialHistogramDataDTO представляет данные дифференциальной гистограммы для HTTP ответа.
type DifferentialHistogramDataDTO struct {
	Bins       []DifferentialHistogramBinDTO `json:"bins"`
	TotalCount int64                         `json:"total_count"`
	BinWidth   int64                         `json:"bin_width"`
	Statistics *DifferentialStatisticsDTO    `json:"statistics"`
}

// DifferentialHistogramBinDTO представляет один бин дифференциальной гистограммы для HTTP ответа.
type DifferentialHistogramBinDTO struct {
	RangeStart int64   `json:"range_start"`
	RangeEnd   int64   `json:"range_end"`
	Count      int64   `json:"count"`
	Frequency  float64 `json:"frequency"`
}

// DifferentialStatisticsDTO представляет статистику дифференциальной гистограммы для HTTP ответа.
type DifferentialStatisticsDTO struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    int64   `json:"min"`
	Max    int64   `json:"max"`
	Count  int64   `json:"count"`
	RMSSD  float64 `json:"rmssd"`
}

// ScatterplotDataDTO представляет данные скаттерограммы для HTTP ответа.
type ScatterplotDataDTO struct {
	Points     []ScatterplotPointDTO     `json:"points"`
	TotalCount int64                     `json:"total_count"`
	Statistics *ScatterplotStatisticsDTO `json:"statistics"`
	Ellipse    *PoincarePlotEllipseDTO   `json:"ellipse"`
}

// ScatterplotPointDTO представляет точку скаттерограммы для HTTP ответа.
type ScatterplotPointDTO struct {
	RRn  int64 `json:"rr_n"`
	RRn1 int64 `json:"rr_n1"`
}

// ScatterplotStatisticsDTO представляет статистику скаттерограммы для HTTP ответа.
type ScatterplotStatisticsDTO struct {
	SD1         float64 `json:"sd1"`
	SD2         float64 `json:"sd2"`
	SD1SD2Ratio float64 `json:"sd1_sd2_ratio"`
	CSI         float64 `json:"csi"`
	CVI         float64 `json:"cvi"`
}

// PoincarePlotEllipseDTO представляет параметры эллипса диаграммы Пуанкаре для HTTP ответа.
type PoincarePlotEllipseDTO struct {
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	SD1     float64 `json:"sd1"`
	SD2     float64 `json:"sd2"`
	Area    float64 `json:"area"`
}

func toRRHistogramDataDTO(histogram *models.RRHistogramData) *RRHistogramDataDTO {
	if histogram == nil {
		return nil
	}

	var bins []HistogramBinDTO
	for _, bin := range histogram.Bins {
		bins = append(bins, HistogramBinDTO{
			RangeStart: bin.RangeStart,
			RangeEnd:   bin.RangeEnd,
			Count:      bin.Count,
			Frequency:  bin.Frequency,
		})
	}

	return &RRHistogramDataDTO{
		Bins:       bins,
		TotalCount: histogram.TotalCount,
		BinWidth:   histogram.BinWidth,
		Statistics: toRRStatisticalSummaryDTO(histogram.Statistics),
	}
}

func toRRTrendAnalysisDTO(trends *models.RRTrendAnalysis) *RRTrendAnalysisDTO {
	if trends == nil {
		return nil
	}

	var trendPoints []TrendPointDTO
	for _, point := range trends.TrendPoints {
		trendPoints = append(trendPoints, TrendPointDTO{
			Time:      point.Time,
			Value:     point.Value,
			Direction: string(point.Direction),
		})
	}

	return &RRTrendAnalysisDTO{
		Period:        trends.Period,
		TrendPoints:   trendPoints,
		OverallTrend:  string(trends.OverallTrend),
		Correlation:   trends.Correlation,
		Seasonality:   trends.Seasonality,
		TrendStrength: trends.TrendStrength,
	}
}

func toRRStatisticalSummaryDTO(stats *models.RRStatisticalSummary) *RRStatisticalSummaryDTO {
	if stats == nil {
		return nil
	}

	return &RRStatisticalSummaryDTO{
		Mean:   stats.Mean,
		StdDev: stats.StdDev,
		Min:    stats.Min,
		Max:    stats.Max,
		Count:  stats.Count,
	}
}

func toHRVMetricsDTO(hrv *models.HRVMetrics) *HRVMetricsDTO {
	if hrv == nil {
		return nil
	}

	return &HRVMetricsDTO{
		RMSSD:           hrv.RMSSD,
		SDNN:            hrv.SDNN,
		PNN50:           hrv.PNN50,
		TriangularIndex: hrv.TriangularIndex,
		TINN:            hrv.TINN,
		VLFPower:        hrv.VLFPower,
		LFPower:         hrv.LFPower,
		HFPower:         hrv.HFPower,
		LFHFRatio:       hrv.LFHFRatio,
		TotalPower:      hrv.TotalPower,
	}
}

// toDifferentialHistogramDataDTO конвертирует данные дифференциальной гистограммы в DTO.
func toDifferentialHistogramDataDTO(histogram *models.DifferentialHistogramData) *DifferentialHistogramDataDTO {
	if histogram == nil {
		return nil
	}

	var bins []DifferentialHistogramBinDTO
	for _, bin := range histogram.Bins {
		bins = append(bins, DifferentialHistogramBinDTO{
			RangeStart: bin.RangeStart,
			RangeEnd:   bin.RangeEnd,
			Count:      bin.Count,
			Frequency:  bin.Frequency,
		})
	}

	return &DifferentialHistogramDataDTO{
		Bins:       bins,
		TotalCount: histogram.TotalCount,
		BinWidth:   histogram.BinWidth,
		Statistics: toDifferentialStatisticsDTO(histogram.Statistics),
	}
}

// toDifferentialStatisticsDTO конвертирует статистику дифференциальной гистограммы в DTO.
func toDifferentialStatisticsDTO(stats *models.DifferentialStatistics) *DifferentialStatisticsDTO {
	if stats == nil {
		return nil
	}

	return &DifferentialStatisticsDTO{
		Mean:   stats.Mean,
		StdDev: stats.StdDev,
		Min:    stats.Min,
		Max:    stats.Max,
		Count:  stats.Count,
		RMSSD:  stats.RMSSD,
	}
}

// toScatterplotDataDTO конвертирует данные скаттерограммы в DTO.
func toScatterplotDataDTO(scatterplot *models.ScatterplotData) *ScatterplotDataDTO {
	if scatterplot == nil {
		return nil
	}

	var points []ScatterplotPointDTO
	for _, point := range scatterplot.Points {
		points = append(points, ScatterplotPointDTO{
			RRn:  point.RRn,
			RRn1: point.RRn1,
		})
	}

	return &ScatterplotDataDTO{
		Points:     points,
		TotalCount: scatterplot.TotalCount,
		Statistics: toScatterplotStatisticsDTO(scatterplot.Statistics),
		Ellipse:    toPoincarePlotEllipseDTO(scatterplot.Ellipse),
	}
}

// toScatterplotStatisticsDTO конвертирует статистику скаттерограммы в DTO.
func toScatterplotStatisticsDTO(stats *models.ScatterplotStatistics) *ScatterplotStatisticsDTO {
	if stats == nil {
		return nil
	}

	return &ScatterplotStatisticsDTO{
		SD1:         stats.SD1,
		SD2:         stats.SD2,
		SD1SD2Ratio: stats.SD1SD2Ratio,
		CSI:         stats.CSI,
		CVI:         stats.CVI,
	}
}

// toPoincarePlotEllipseDTO конвертирует параметры эллипса Пуанкаре в DTO.
func toPoincarePlotEllipseDTO(ellipse *models.PoincarePlotEllipse) *PoincarePlotEllipseDTO {
	if ellipse == nil {
		return nil
	}

	return &PoincarePlotEllipseDTO{
		CenterX: ellipse.CenterX,
		CenterY: ellipse.CenterY,
		SD1:     ellipse.SD1,
		SD2:     ellipse.SD2,
		Area:    ellipse.Area,
	}
}
