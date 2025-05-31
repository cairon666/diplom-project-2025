import { useState, useCallback } from 'react';
import {
    Card,
    CardContent,
} from 'ui/card';
import { Button } from 'ui/button';
import { Alert, AlertDescription } from 'ui/alert';
import { LoadingSpinner } from 'ui/loading-spinner';
import { LoadingOverlay } from 'ui/loading-overlay';
import {
    RefreshCw,
    BarChart3,
    TrendingUp,
    Heart,
    Activity,
    ScatterChart,
    Play,
    AlertTriangle,
    Info,
    Eye,
    EyeOff,
    ChevronDown,
    ChevronRight,
} from 'lucide-react';
import { subMinutes, differenceInMinutes } from 'date-fns';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import type { SerializedError } from '@reduxjs/toolkit';
import { useNavigate, useSearch } from '@tanstack/react-router';

import { DateTimeRangeSelector, DateTimeRange } from './DateTimeRangeSelector';
import { RRStatistics } from './RRStatistics';
import { RRHistogramChart } from './RRHistogramChart';
import { RRDifferentialHistogramChart } from './RRDifferentialHistogramChart';
import { RRScatterplotChart } from './RRScatterplotChart';
import { RRTrendsChart } from './RRTrendsChart';
import { HeartRateChart } from './HeartRateChart';

import {
    useGetRRStatisticsQuery,
    useGetHistogramQuery,
    useGetDifferentialHistogramQuery,
    useGetScatterplotQuery,
    useGetTrendsQuery,
} from '../api/rrAnalyticsApi';
import { useGetRRIntervalsQuery } from '../api/rrIntervalsApi';
import { 
    ApiErrorCode, 
    isBaseError, 
    isApiErrorWithFields,
    type AppBaseError 
} from '../../../store/apiErrors';

// Импортируем новые хуки и компоненты из ui/
import { useChartVisibility } from '../hooks/useChartVisibility';
import { usePeriodValidation } from '../hooks/usePeriodValidation';
import { useApiErrorHandling } from '../hooks/useApiErrorHandling';
import { ChartVisibilityControls } from './ChartControls';
import { ChartSection } from './ChartSection';
import { ErrorAlert } from './AlertComponents';

export function RRAnalyticsDashboard() {
    const navigate = useNavigate();
    const search = useSearch({ strict: false }) as {
        from?: string;
        to?: string;
    };

    // Функция для получения дат из URL или возврата дефолтных значений
    const getInitialDateRange = useCallback((): DateTimeRange => {
        const defaultDateRange = {
            from: subMinutes(new Date(), 5), // последние 5 минут для детального анализа
            to: new Date(),
        };

        // Попытка восстановить dateRange из URL
        if (search.from && search.to) {
            try {
                const fromDate = new Date(search.from);
                const toDate = new Date(search.to);
                
                // Проверяем, что даты валидны
                if (!isNaN(fromDate.getTime()) && !isNaN(toDate.getTime())) {
                    return { from: fromDate, to: toDate };
                }
            } catch {
                // Если ошибка парсинга, возвращаем дефолтные значения
            }
        }

        return defaultDateRange;
    }, [search]);

    // Инициализируем dateRange из URL или дефолтными значениями
    const [dateRange, setDateRange] = useState<DateTimeRange>(getInitialDateRange());
    // Состояние для примененного диапазона дат (для API запросов)
    const [appliedDateRange, setAppliedDateRange] = useState<DateTimeRange>(getInitialDateRange());
    // Состояние для контроля загрузки данных
    const [shouldFetchData, setShouldFetchData] = useState(true);

    // Используем новые хуки
    const { visibleCharts, toggleChart, showAllCharts, hideAllCharts } = useChartVisibility();
    const validation = usePeriodValidation(appliedDateRange);
    const { getErrorMessage, combineErrors } = useApiErrorHandling();

    // Обновляем URL при изменении dateRange
    const updateURL = useCallback((newDateRange: DateTimeRange) => {
        navigate({
            to: '/panel',
            search: {
                from: newDateRange.from.toISOString(),
                to: newDateRange.to.toISOString(),
            },
            replace: true, // Заменяем текущую запись в истории
        });
    }, [navigate]);

    // Обработчик изменения dateRange с обновлением URL (но без запроса данных)
    const handleDateRangeChange = useCallback((newDateRange: DateTimeRange) => {
        setDateRange(newDateRange);
        updateURL(newDateRange);
    }, [updateURL]);

    // Обработчик применения изменений (загрузка данных)
    const handleApplyChanges = useCallback(() => {
        setAppliedDateRange(dateRange);
        setShouldFetchData(true);
    }, [dateRange]);

    // Проверяем, есть ли неприменённые изменения
    const hasUnappliedChanges = dateRange.from.getTime() !== appliedDateRange.from.getTime() || 
                               dateRange.to.getTime() !== appliedDateRange.to.getTime();

    // Функция для расчета продолжительности в минутах
    const getDurationInMinutes = (): number => {
        return differenceInMinutes(appliedDateRange.to, appliedDateRange.from);
    };

    // API запросы с условиями skip, включающими валидацию и видимость
    const {
        data: statisticsData,
        isLoading: statisticsLoading,
        isFetching: statisticsFetching,
        error: statisticsError,
        refetch: refetchStatistics,
    } = useGetRRStatisticsQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
        include_histogram: false,
        include_hrv: true,
    }, { skip: !shouldFetchData || !validation.isValidForBasicAnalysis() });

    const {
        data: histogramData,
        isLoading: histogramLoading,
        isFetching: histogramFetching,
        error: histogramError,
        refetch: refetchHistogram,
    } = useGetHistogramQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
        bins_count: 25,
    }, { skip: !shouldFetchData || !validation.isValidForHistogram() || !visibleCharts.histogram });

    const {
        data: differentialHistogramData,
        isLoading: differentialHistogramLoading,
        isFetching: differentialHistogramFetching,
        error: differentialHistogramError,
        refetch: refetchDifferentialHistogram,
    } = useGetDifferentialHistogramQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
        bins_count: 25,
    }, { skip: !shouldFetchData || !validation.isValidForHistogram() || !visibleCharts.differentialHistogram });

    const {
        data: scatterplotData,
        isLoading: scatterplotLoading,
        isFetching: scatterplotFetching,
        error: scatterplotError,
        refetch: refetchScatterplot,
    } = useGetScatterplotQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
    }, { skip: !shouldFetchData || !validation.isValidForScatterplot() || !visibleCharts.scatterplot });

    const {
        data: trendsData,
        isLoading: trendsLoading,
        isFetching: trendsFetching,
        error: trendsError,
        refetch: refetchTrends,
    } = useGetTrendsQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
        window_size: 5,
    }, { skip: !shouldFetchData || !validation.isValidForTrends() || !visibleCharts.trends });

    const {
        data: rrIntervalsData,
        isLoading: rrIntervalsLoading,
        isFetching: rrIntervalsFetching,
        error: rrIntervalsError,
        refetch: refetchRRIntervals,
    } = useGetRRIntervalsQuery({
        from: appliedDateRange.from.toISOString(),
        to: appliedDateRange.to.toISOString(),
    }, { skip: !shouldFetchData || !validation.isValidForDetailedHeartRate() || !visibleCharts.heartRate });

    const isLoading = statisticsLoading || histogramLoading || differentialHistogramLoading || scatterplotLoading || trendsLoading || rrIntervalsLoading;
    const isFetching = statisticsFetching ||
        histogramFetching ||
        differentialHistogramFetching ||
        scatterplotFetching ||
        trendsFetching ||
        rrIntervalsFetching;

    // Объединяем все ошибки
    const currentError = combineErrors(
        statisticsError, 
        histogramError, 
        differentialHistogramError, 
        scatterplotError, 
        trendsError, 
        rrIntervalsError
    );

    return (
        <div className="space-y-6">
            {/* Заголовок и управление */}
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div>
                    <div className="flex items-center gap-2">
                        <h2 className="text-2xl font-bold tracking-tight">
                            Аналитика R-R интервалов
                        </h2>
                    </div>
                    <p className="text-muted-foreground">
                        Анализ вариабельности сердечного ритма и трендов
                    </p>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                    {/* Управление графиками */}
                    <ChartVisibilityControls
                        showAllCharts={showAllCharts}
                        hideAllCharts={hideAllCharts}
                    />
                </div>
            </div>

            {/* Селектор диапазона дат */}
            <DateTimeRangeSelector
                value={dateRange}
                onChange={handleDateRangeChange}
                className="w-full"
            />

            {/* Кнопка применения изменений */}
            <div className="flex justify-center">
                <Button
                    onClick={handleApplyChanges}
                    disabled={!hasUnappliedChanges || isFetching}
                    size="lg"
                    className="px-8"
                >
                    {isFetching ? (
                        <>
                            <LoadingSpinner className="h-4 w-4 mr-2" />
                            Загрузка...
                        </>
                    ) : (
                        <>
                            <Play className="h-4 w-4 mr-2" />
                            {hasUnappliedChanges ? 'Применить изменения' : 'Данные актуальны'}
                        </>
                    )}
                </Button>
            </div>

            {/* Состояния загрузки - только для первоначальной загрузки */}
            {isLoading && !isFetching && (
                <div className="flex justify-center py-8">
                    <LoadingSpinner />
                    <span className="ml-2">Загрузка данных...</span>
                </div>
            )}

            {/* Показываем контент только если данные должны загружаться */}
            {shouldFetchData && (
                <div className="space-y-6">
                    {/* Основная статистика */}
                    <LoadingOverlay
                        isLoading={statisticsLoading || statisticsFetching}
                    >
                            <RRStatistics
                                statistics={statisticsData?.summary}
                                hrvMetrics={statisticsData?.hrv_metrics}
                                isLoading={statisticsLoading}
                            />
                    </LoadingOverlay>

                    {/* Ошибки */}
                    {currentError && (
                        <ErrorAlert 
                            error={currentError as any} 
                            getErrorMessage={getErrorMessage} 
                        />
                    )}

                    {/* Детальный анализ - График пульса */}
                    <ChartSection
                        chartName="heartRate"
                        title="Детальный анализ пульса"
                        icon={Heart}
                        isVisible={visibleCharts.heartRate}
                        isValid={validation.isValidForDetailedHeartRate()}
                        onToggle={toggleChart}
                        invalidMessage={!validation.isValidForDetailedHeartRate() ? validation.getInvalidMessage('detailedHeartRate') : undefined}
                        warningMessage={validation.shouldShowHeartRateWarning() ? (
                            <div>
                                <strong>Большой временной диапазон ({validation.getDurationInMinutes()} мин)</strong>
                                <br />
                                Для оптимальной производительности детальный анализ пульса рекомендуется проводить 
                                для периодов менее 15 минут. Данные будут агрегированы для лучшей визуализации.
                            </div>
                        ) : undefined}
                    >
                        <LoadingOverlay
                            isLoading={rrIntervalsFetching}
                            loadingText="Обновление данных пульса..."
                        >
                            <HeartRateChart 
                                dateRange={appliedDateRange} 
                                onChange={handleDateRangeChange}
                                externalData={rrIntervalsData}
                                externalLoading={rrIntervalsLoading}
                                externalError={rrIntervalsError as FetchBaseQueryError | SerializedError | undefined}
                            />
                        </LoadingOverlay>
                    </ChartSection>

                    {/* Гистограмма */}
                    <ChartSection
                        chartName="histogram"
                        title="Гистограмма R-R интервалов"
                        icon={BarChart3}
                        isVisible={visibleCharts.histogram}
                        isValid={validation.isValidForHistogram()}
                        onToggle={toggleChart}
                        invalidMessage={!validation.isValidForHistogram() ? validation.getInvalidMessage('histogram') : undefined}
                        infoMessage={validation.getDurationInMinutes() < 2 ? (
                            <div>
                                <strong>Рекомендация:</strong> Для более информативной гистограммы 
                                рекомендуется использовать периоды от 2 минут и более.
                            </div>
                        ) : undefined}
                    >
                        <LoadingOverlay
                            isLoading={histogramFetching}
                            loadingText="Обновление гистограммы..."
                        >
                            {histogramLoading ? (
                                <div className="flex justify-center py-8">
                                    <LoadingSpinner />
                                </div>
                            ) : histogramData ? (
                                <RRHistogramChart data={histogramData.histogram} />
                            ) : histogramError ? (
                                <Card>
                                    <CardContent>
                                        <ErrorAlert 
                                            title="Ошибка загрузки гистограммы"
                                            error={histogramError} 
                                            getErrorMessage={getErrorMessage} 
                                        />
                                    </CardContent>
                                </Card>
                            ) : (
                                <Card>
                                    <CardContent className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для отображения гистограммы
                                        </p>
                                    </CardContent>
                                </Card>
                            )}
                        </LoadingOverlay>
                    </ChartSection>

                    {/* Дифференциальная гистограмма */}
                    <ChartSection
                        chartName="differentialHistogram"
                        title="Дифференциальная гистограмма R-R интервалов"
                        icon={Activity}
                        isVisible={visibleCharts.differentialHistogram}
                        isValid={validation.isValidForHistogram()}
                        onToggle={toggleChart}
                        invalidMessage={!validation.isValidForHistogram() ? validation.getInvalidMessage('histogram') : undefined}
                    >
                        <LoadingOverlay
                            isLoading={differentialHistogramFetching}
                            loadingText="Обновление дифференциальной гистограммы..."
                        >
                            <RRDifferentialHistogramChart
                                from={appliedDateRange.from.toISOString()}
                                to={appliedDateRange.to.toISOString()}
                                binsCount={25}
                            />
                        </LoadingOverlay>
                    </ChartSection>

                    {/* Скаттерограмма (диаграмма Пуанкаре) */}
                    <ChartSection
                        chartName="scatterplot"
                        title="Скаттерограмма R-R интервалов (диаграмма Пуанкаре)"
                        icon={ScatterChart}
                        isVisible={visibleCharts.scatterplot}
                        isValid={validation.isValidForScatterplot()}
                        onToggle={toggleChart}
                        invalidMessage={!validation.isValidForScatterplot() ? validation.getInvalidMessage('scatterplot') : undefined}
                        infoMessage={validation.getDurationInMinutes() < 5 ? (
                            <div>
                                <strong>Рекомендация:</strong> Диаграмма Пуанкаре наиболее информативна 
                                для периодов от 5 минут, когда доступно достаточное количество R-R интервалов 
                                для анализа вариабельности.
                            </div>
                        ) : undefined}
                    >
                        <LoadingOverlay
                            isLoading={scatterplotFetching}
                            loadingText="Обновление скаттерограммы..."
                        >
                            <RRScatterplotChart
                                from={appliedDateRange.from.toISOString()}
                                to={appliedDateRange.to.toISOString()}
                            />
                        </LoadingOverlay>
                    </ChartSection>

                    {/* Анализ трендов */}
                    <ChartSection
                        chartName="trends"
                        title="Анализ трендов R-R интервалов"
                        icon={TrendingUp}
                        isVisible={visibleCharts.trends}
                        isValid={validation.isValidForTrends()}
                        onToggle={toggleChart}
                        invalidMessage={!validation.isValidForTrends() ? validation.getInvalidMessage('trends') : undefined}
                        warningMessage={validation.shouldShowTrendsWarning() ? (
                            <div>
                                <strong>Малый временной диапазон ({validation.getDurationInMinutes()} мин)</strong>
                                <br />
                                Для достоверного анализа трендов рекомендуется использовать периоды более 10 минут. 
                                Текущий диапазон может не обеспечить статистически значимые результаты.
                            </div>
                        ) : undefined}
                    >
                        <LoadingOverlay
                            isLoading={trendsFetching}
                            loadingText="Обновление трендов..."
                        >
                            {trendsLoading ? (
                                <div className="flex justify-center py-8">
                                    <LoadingSpinner />
                                </div>
                            ) : trendsData ? (
                                <RRTrendsChart data={trendsData.trend_analysis} />
                            ) : trendsError ? (
                                <Card>
                                    <CardContent>
                                        <ErrorAlert 
                                            title="Ошибка загрузки трендов"
                                            error={trendsError} 
                                            getErrorMessage={getErrorMessage} 
                                        />
                                    </CardContent>
                                </Card>
                            ) : (
                                <Card>
                                    <CardContent className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для анализа трендов
                                        </p>
                                    </CardContent>
                                </Card>
                            )}
                        </LoadingOverlay>
                    </ChartSection>
                </div>
            )}

            {/* Сообщение если данные не загружены */}
            {!shouldFetchData && (
                <Card>
                    <CardContent className="py-12 text-center">
                        <p className="text-muted-foreground text-lg">
                            Выберите диапазон дат и нажмите "Применить изменения" для загрузки данных
                        </p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
