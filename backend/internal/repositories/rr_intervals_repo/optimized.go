package rr_intervals_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// GetAggregatedByInterval получает агрегированные данные по временным интервалам.
func (r *RRIntervalsRepo) GetAggregatedByInterval(ctx context.Context, userID uuid.UUID, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error) {
	// Простой InfluxDB запрос с aggregateWindow
	query := fmt.Sprintf(`
		from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r._measurement == "rr_intervals")
		  |> filter(fn: (r) => r.user_id == "%s")
		  |> filter(fn: (r) => r._field == "rr_interval_ms")
		  |> filter(fn: (r) => r._value >= 300 and r._value <= 2000)
		  |> aggregateWindow(every: %dm, fn: mean, createEmpty: false)
		  |> sort(columns: ["_time"])
	`,
		r.bucket, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano),
		userID.String(), intervalMinutes,
	)

	queryAPI := r.influxClient.QueryAPI(r.org)
	queryResult, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to query aggregated data: %v", err)
	}
	defer queryResult.Close()

	var aggregatedResult []models.AggregatedRRData
	for queryResult.Next() {
		record := queryResult.Record()

		var meanValue float64
		if floatVal, ok := record.Value().(float64); ok {
			meanValue = floatVal
		} else if intVal, ok := record.Value().(int64); ok {
			meanValue = float64(intVal)
		} else {
			continue
		}

		aggregatedResult = append(aggregatedResult, models.AggregatedRRData{
			Time:   record.Time(),
			Mean:   meanValue,
			StdDev: 0, // InfluxDB aggregateWindow doesn't calculate stddev
			Min:    0,
			Max:    0,
			Count:  1, // We don't get count from aggregateWindow mean
		})
	}

	if queryResult.Err() != nil {
		return nil, apperrors.DataProcessingErrorf("error reading aggregated data: %v", queryResult.Err())
	}

	return aggregatedResult, nil
}

// GetCompleteAnalysisData получает все данные для комплексного анализа
// Простая реализация для соответствия интерфейсу.
func (r *RRIntervalsRepo) GetCompleteAnalysisData(ctx context.Context, userID uuid.UUID, from, to time.Time, options models.CompleteAnalysisOptions) (*models.CompleteAnalysisData, error) {
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	startTime := time.Now()

	// Получаем сырые данные
	rawValues, err := r.GetRawValuesForAnalysis(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get raw values: %v", err)
	}

	// Получаем статистику
	statistics, err := r.GetStatisticalSummary(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get statistics: %v", err)
	}

	// Получаем агрегированные данные
	aggregatedData, err := r.GetAggregatedByInterval(ctx, userID, from, to, options.AggregationIntervalMinutes)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get aggregated data: %v", err)
	}

	// Если нет данных, возвращаем пустой результат с базовой структурой
	if len(rawValues) == 0 {
		return &models.CompleteAnalysisData{
			RawValues: []int64{},
			TimeRange: models.TimeRange{
				From: from,
				To:   to,
			},
			Statistics:     &models.RRStatisticalSummary{},
			HRVMetrics:     &models.HRVMetrics{},
			AggregatedData: []models.AggregatedRRData{},
			TrendAnalysis: &models.RRTrendAnalysis{
				Period:        fmt.Sprintf("%v to %v", from.Format("15:04:05"), to.Format("15:04:05")),
				TrendPoints:   []models.TrendPoint{},
				OverallTrend:  "stable",
				Correlation:   0,
				Seasonality:   []float64{},
				TrendStrength: 0,
			},
			Histogram: &models.RRHistogramData{
				Bins:       []models.HistogramBin{},
				TotalCount: 0,
				BinWidth:   0,
				Statistics: &models.RRStatisticalSummary{},
			},
			DiffHistogram: &models.DifferentialHistogramData{
				Bins:       []models.DifferentialHistogramBin{},
				TotalCount: 0,
				BinWidth:   0,
				Statistics: &models.DifferentialStatistics{},
			},
			Scatterplot: &models.ScatterplotData{
				Points:     []models.ScatterplotPoint{},
				TotalCount: 0,
				Statistics: &models.ScatterplotStatistics{},
				Ellipse:    &models.PoincarePlotEllipse{},
			},
			ProcessingTime: time.Since(startTime),
			DataQuality: &models.DataQualityMetrics{
				TotalMeasurements:   0,
				ValidMeasurements:   0,
				InvalidMeasurements: 0,
				QualityPercentage:   0,
				MissingDataGaps:     0,
				LargestGapDuration:  0,
				AverageSamplingRate: 0,
			},
		}, nil
	}

	// Подготавливаем данные для возврата
	var rawValuesForReturn []int64
	if options.IncludeRawData {
		rawValuesForReturn = rawValues
	}

	processingTime := time.Since(startTime)

	return &models.CompleteAnalysisData{
		RawValues: rawValuesForReturn,
		TimeRange: models.TimeRange{
			From: from,
			To:   to,
		},
		Statistics: statistics,
		HRVMetrics: &models.HRVMetrics{
			// Базовые HRV метрики будут вычислены в сервисном слое
		},
		AggregatedData: aggregatedData,
		TrendAnalysis: &models.RRTrendAnalysis{
			Period:        fmt.Sprintf("%v to %v", from.Format("15:04:05"), to.Format("15:04:05")),
			TrendPoints:   []models.TrendPoint{},
			OverallTrend:  "stable",
			Correlation:   0,
			Seasonality:   []float64{},
			TrendStrength: 0,
		},
		Histogram: &models.RRHistogramData{
			Bins:       []models.HistogramBin{},
			TotalCount: int64(len(rawValues)),
			BinWidth:   0,
			Statistics: statistics,
		},
		DiffHistogram: &models.DifferentialHistogramData{
			Bins:       []models.DifferentialHistogramBin{},
			TotalCount: int64(len(rawValues)),
			BinWidth:   0,
			Statistics: &models.DifferentialStatistics{},
		},
		Scatterplot: &models.ScatterplotData{
			Points:     []models.ScatterplotPoint{},
			TotalCount: int64(len(rawValues)),
			Statistics: &models.ScatterplotStatistics{},
			Ellipse:    &models.PoincarePlotEllipse{},
		},
		ProcessingTime: processingTime,
		DataQuality: &models.DataQualityMetrics{
			TotalMeasurements:   int64(len(rawValues)),
			ValidMeasurements:   int64(len(rawValues)),
			InvalidMeasurements: 0,
			QualityPercentage:   100.0,
			MissingDataGaps:     0,
			LargestGapDuration:  0,
			AverageSamplingRate: calculateSamplingRate(from, to, int64(len(rawValues))),
		},
	}, nil
}

// Вспомогательная функция для расчета частоты дискретизации.
func calculateSamplingRate(from, to time.Time, measurements int64) float64 {
	if measurements == 0 {
		return 0
	}
	duration := to.Sub(from).Seconds()
	if duration == 0 {
		return 0
	}

	return float64(measurements) / duration
}
