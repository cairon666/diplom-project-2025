import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

// Типы анализа
export type AnalysisType = 'full' | 'histogram' | 'differential_histogram' | 'scatterplot' | 'trends' | 'statistics' | 'hrv';
export type FilterType = 'artifacts' | 'outliers';

// Базовые интерфейсы
export interface AnalyticsTimeRange {
    from: string;
    to: string;
}

export interface RRStatisticalSummary {
    mean: number;
    std_dev: number;
    min: number;
    max: number;
    count: number;
}

export interface HistogramBin {
    range_start: number;
    range_end: number;
    count: number;
    frequency: number;
}

export interface RRHistogramData {
    bins: HistogramBin[];
    total_count: number;
    bin_width: number;
    statistics: RRStatisticalSummary;
}

export interface TrendPoint {
    time: string;
    value: number;
    direction: 'up' | 'down' | 'stable';
}

export interface RRTrendAnalysis {
    period: string;
    trend_points: TrendPoint[];
    overall_trend: 'increasing' | 'decreasing' | 'stable';
    correlation: number;
    seasonality: number[];
    trend_strength: number;
}

export interface HRVMetrics {
    rmssd: number;
    sdnn: number;
    pnn50: number;
    triangular_index: number;
    tinn: number;
    vlf_power: number;
    lf_power: number;
    hf_power: number;
    lf_hf_ratio: number;
    total_power: number;
}

export interface DataQualityReport {
    total_measurements: number;
    valid_measurements: number;
    artifact_count: number;
    outlier_count: number;
    quality_score: number;
    recommended_actions: string[];
    data_completeness: number;
}

// Запросы и ответы

export interface AnalyzeRRIntervalsRequest {
    from: string;
    to: string;
    analysis_type: AnalysisType;
    bins_count?: number;
    window_size?: number;
    filters?: FilterType[];
}

export interface AnalyzeRRIntervalsResponse {
    user_id: string;
    period: string;
    generated_at: string;
    analysis_type: AnalysisType;
    histogram?: RRHistogramData;
    differential_histogram?: DifferentialHistogramData;
    scatterplot?: ScatterplotData;
    trend_analysis?: RRTrendAnalysis;
    statistics?: RRStatisticalSummary;
    hrv_metrics?: HRVMetrics;
    quality_report?: DataQualityReport;
    recommendations?: string[];
}

export interface GetRRStatisticsParams {
    from: string;
    to: string;
    include_histogram?: boolean;
    include_hrv?: boolean;
    bins_count?: number;
}

export interface GetRRStatisticsResponse {
    summary: RRStatisticalSummary;
    histogram?: RRHistogramData;
    hrv_metrics?: HRVMetrics;
    time_range: AnalyticsTimeRange;
}

export interface GetAggregatedDataParams {
    from: string;
    to: string;
    interval: number; // минуты
}

export interface AggregatedRRData {
    time: string;
    mean: number;
    std_dev: number;
    min: number;
    max: number;
    count: number;
}

export interface GetAggregatedDataResponse {
    data: AggregatedRRData[];
    time_range: AnalyticsTimeRange;
    interval_minutes: number;
}

export interface GetHistogramParams {
    from: string;
    to: string;
    bins_count?: number;
}

export interface GetHistogramResponse {
    histogram: RRHistogramData;
    time_range: AnalyticsTimeRange;
}

export interface GetTrendsParams {
    from: string;
    to: string;
    window_size?: number;
}

export interface GetTrendsResponse {
    trend_analysis: RRTrendAnalysis;
    time_range: AnalyticsTimeRange;
}

export interface GetHRVMetricsParams {
    from: string;
    to: string;
}

export interface GetHRVMetricsResponse {
    hrv_metrics: HRVMetrics;
    time_range: AnalyticsTimeRange;
}

// Дифференциальная гистограмма
export interface DifferentialHistogramBin {
    range_start: number;
    range_end: number;
    count: number;
    frequency: number;
}

export interface DifferentialStatistics {
    mean: number;
    std_dev: number;
    min: number;
    max: number;
    count: number;
    rmssd: number;
}

export interface DifferentialHistogramData {
    bins: DifferentialHistogramBin[];
    total_count: number;
    bin_width: number;
    statistics: DifferentialStatistics;
}

// Скаттерограмма (диаграмма Пуанкаре)
export interface ScatterplotPoint {
    rr_n: number;
    rr_n1: number;
}

export interface ScatterplotStatistics {
    sd1: number;
    sd2: number;
    sd1_sd2_ratio: number;
    csi: number;
    cvi: number;
}

export interface PoincarePlotEllipse {
    center_x: number;
    center_y: number;
    sd1: number;
    sd2: number;
    area: number;
}

export interface ScatterplotData {
    points: ScatterplotPoint[];
    total_count: number;
    statistics: ScatterplotStatistics;
    ellipse: PoincarePlotEllipse;
}

// Дифференциальная гистограмма
export interface GetDifferentialHistogramParams {
    from: string;
    to: string;
    bins_count?: number;
}

export interface GetDifferentialHistogramResponse {
    bins: DifferentialHistogramBin[];
    total_count: number;
    bin_width: number;
    statistics: DifferentialStatistics;
}

// Скаттерограмма
export interface GetScatterplotParams {
    from: string;
    to: string;
}

export interface GetScatterplotResponse {
    points: ScatterplotPoint[];
    total_count: number;
    statistics: ScatterplotStatistics;
    ellipse: PoincarePlotEllipse;
}

// Метрики качества данных
export interface DataQualityMetrics {
    total_measurements: number;
    valid_measurements: number;
    invalid_measurements: number;
    quality_percentage: number;
    missing_data_gaps: number;
    largest_gap_duration: string;
    average_sampling_rate: number;
}

// Комплексный анализ - новый оптимизированный метод
export interface CompleteAnalysisOptions {
    aggregation_interval_minutes?: number;
    trend_window_size_minutes?: number;
    histogram_bins_count?: number;
    diff_histogram_bins_count?: number;
    enable_frequency_domain_analysis?: boolean;
    include_raw_data?: boolean;
    max_data_points?: number;
}

export interface CompleteAnalysisData {
    raw_values?: number[];
    time_range: AnalyticsTimeRange;
    statistics: RRStatisticalSummary;
    hrv_metrics: HRVMetrics;
    aggregated_data: AggregatedRRData[];
    trend_analysis: RRTrendAnalysis;
    histogram: RRHistogramData;
    differential_histogram: DifferentialHistogramData;
    scatterplot: ScatterplotData;
    processing_time: string;
    data_quality: DataQualityMetrics;
}

export interface GetCompleteAnalysisParams {
    from: string;
    to: string;
    aggregation_interval?: number;
    trend_window_size?: number;
    histogram_bins?: number;
    diff_histogram_bins?: number;
    include_raw_data?: boolean;
    max_data_points?: number;
    quick?: boolean;
}

export interface GetCompleteAnalysisResponse {
    data: CompleteAnalysisData;
    success: boolean;
    message?: string;
}

export const rrAnalyticsApi = createApi({
    reducerPath: 'rrAnalyticsApi',
    baseQuery: baseQueryWithReauth,
    tagTypes: ['RRAnalytics'],
    keepUnusedDataFor: 0,
    endpoints: (builder) => ({
        // Полный анализ
        analyzeRRIntervals: builder.mutation<AnalyzeRRIntervalsResponse, AnalyzeRRIntervalsRequest>({
            query: (body) => ({
                url: '/v1/rr-intervals/analytics/analyze',
                method: 'POST',
                body,
            }),
            invalidatesTags: ['RRAnalytics'],
        }),

        // Статистика
        getRRStatistics: builder.query<GetRRStatisticsResponse, GetRRStatisticsParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/statistics',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Агрегированные данные
        getAggregatedData: builder.query<GetAggregatedDataResponse, GetAggregatedDataParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/aggregated',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Гистограмма
        getHistogram: builder.query<GetHistogramResponse, GetHistogramParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/histogram',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Тренды
        getTrends: builder.query<GetTrendsResponse, GetTrendsParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/trends',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // HRV метрики
        getHRVMetrics: builder.query<GetHRVMetricsResponse, GetHRVMetricsParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/hrv',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Дифференциальная гистограмма
        getDifferentialHistogram: builder.query<GetDifferentialHistogramResponse, GetDifferentialHistogramParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/differential-histogram',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Скаттерограмма
        getScatterplot: builder.query<GetScatterplotResponse, GetScatterplotParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/scatterplot',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),

        // Комплексный анализ (новый оптимизированный метод)
        getCompleteAnalysis: builder.query<GetCompleteAnalysisResponse, GetCompleteAnalysisParams>({
            query: (params) => ({
                url: '/v1/rr-intervals/analytics/complete',
                params,
            }),
            providesTags: ['RRAnalytics'],
            keepUnusedDataFor: 0,
        }),
    }),
});

export const {
    useAnalyzeRRIntervalsMutation,
    useGetRRStatisticsQuery,
    useGetAggregatedDataQuery,
    useGetHistogramQuery,
    useGetTrendsQuery,
    useGetHRVMetricsQuery,
    useGetDifferentialHistogramQuery,
    useGetScatterplotQuery,
    useGetCompleteAnalysisQuery,
} = rrAnalyticsApi; 