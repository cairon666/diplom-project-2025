import { useState, useMemo, useCallback } from 'react';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { Button } from 'ui/button';
import { LoadingSpinner } from 'ui/loading-spinner';
import { Play } from 'lucide-react';
import { subDays, subHours } from 'date-fns';

import type { DateTimeRange } from './DateTimeRangeSelector';
import {
    useGetRRStatisticsQuery,
    useGetHistogramQuery,
    useGetDifferentialHistogramQuery,
    useGetScatterplotQuery,
    useGetTrendsQuery,
} from '../api/rrAnalyticsApi';

// Импортируем новые компоненты
import { PeriodSelector } from './PeriodSelector';
import { PeriodStatistics } from './PeriodStatistics';
import { ComparisonStats } from './ComparisonStats';
import { ComparisonCharts } from './ComparisonCharts';

// Импортируем хуки и компоненты из ui/
import { useChartVisibility } from '../hooks/useChartVisibility';
import { useApiErrorHandling } from '../hooks/useApiErrorHandling';
import { ChartVisibilityControls } from './ChartControls';
import { ErrorAlert } from './AlertComponents';

export function ComparativeAnalysisDashboard() {
    const navigate = useNavigate();
    const search = useSearch({ strict: false }) as {
        period1_from?: string;
        period1_to?: string;
        period2_from?: string;
        period2_to?: string;
    };

    // Функция для получения дат из URL или возврата дефолтных значений
    const getInitialPeriods = useCallback(() => {
        const defaultPeriod1 = {
            from: subHours(new Date(), 4), // последние 4 часа
            to: new Date(),
        };

        const defaultPeriod2 = {
            from: subHours(subDays(new Date(), 1), 4), // вчера, те же 4 часа
            to: subDays(new Date(), 1),
        };

        // Попытка восстановить период 1 из URL
        let period1 = defaultPeriod1;
        if (search.period1_from && search.period1_to) {
            try {
                period1 = {
                    from: new Date(search.period1_from),
                    to: new Date(search.period1_to),
                };
                // Проверяем, что даты валидны
                if (
                    isNaN(period1.from.getTime()) ||
                    isNaN(period1.to.getTime())
                ) {
                    period1 = defaultPeriod1;
                }
            } catch {
                period1 = defaultPeriod1;
            }
        }

        // Попытка восстановить период 2 из URL
        let period2 = defaultPeriod2;
        if (search.period2_from && search.period2_to) {
            try {
                period2 = {
                    from: new Date(search.period2_from),
                    to: new Date(search.period2_to),
                };
                // Проверяем, что даты валидны
                if (
                    isNaN(period2.from.getTime()) ||
                    isNaN(period2.to.getTime())
                ) {
                    period2 = defaultPeriod2;
                }
            } catch {
                period2 = defaultPeriod2;
            }
        }

        return { period1, period2 };
    }, [search]);

    // Инициализируем периоды из URL или дефолтными значениями
    const initialPeriods = getInitialPeriods();

    // Состояние для текущих выбранных периодов (не вызывают API запросы)
    const [period1, setPeriod1] = useState<DateTimeRange>(
        initialPeriods.period1,
    );
    const [period2, setPeriod2] = useState<DateTimeRange>(
        initialPeriods.period2,
    );

    // Состояние для примененных периодов (используются в API запросах)
    const [appliedPeriod1, setAppliedPeriod1] = useState<DateTimeRange>(
        initialPeriods.period1,
    );
    const [appliedPeriod2, setAppliedPeriod2] = useState<DateTimeRange>(
        initialPeriods.period2,
    );

    // Состояние для контроля загрузки данных
    const [shouldFetchData, setShouldFetchData] = useState(true);

    // Используем новые хуки
    const { visibleCharts, toggleChart, showAllCharts, hideAllCharts } =
        useChartVisibility();
    const { getErrorMessage } = useApiErrorHandling();

    // Обновляем URL при изменении периодов
    const updateURL = useCallback(
        (newPeriod1: DateTimeRange, newPeriod2: DateTimeRange) => {
            navigate({
                to: '/panel/comparative-analysis',
                search: {
                    period1_from: newPeriod1.from.toISOString(),
                    period1_to: newPeriod1.to.toISOString(),
                    period2_from: newPeriod2.from.toISOString(),
                    period2_to: newPeriod2.to.toISOString(),
                },
                replace: true, // Заменяем текущую запись в истории
            });
        },
        [navigate],
    );

    // Обработчики изменения периодов с обновлением URL (но без запроса данных)
    const handlePeriod1Change = useCallback(
        (newPeriod: DateTimeRange) => {
            setPeriod1(newPeriod);
            updateURL(newPeriod, period2);
        },
        [period2, updateURL],
    );

    const handlePeriod2Change = useCallback(
        (newPeriod: DateTimeRange) => {
            setPeriod2(newPeriod);
            updateURL(period1, newPeriod);
        },
        [period1, updateURL],
    );

    // Обработчик применения изменений (загрузка данных)
    const handleCompare = useCallback(() => {
        setAppliedPeriod1(period1);
        setAppliedPeriod2(period2);
        setShouldFetchData(true);
    }, [period1, period2]);

    // Проверяем, есть ли неприменённые изменения
    const hasUnappliedChanges =
        period1.from.getTime() !== appliedPeriod1.from.getTime() ||
        period1.to.getTime() !== appliedPeriod1.to.getTime() ||
        period2.from.getTime() !== appliedPeriod2.from.getTime() ||
        period2.to.getTime() !== appliedPeriod2.to.getTime();

    const [isComparing, setIsComparing] = useState(false);

    // API запросы для первого периода
    const {
        data: statisticsData1,
        isLoading: statisticsLoading1,
        isFetching: statisticsFetching1,
        error: statisticsError1,
    } = useGetRRStatisticsQuery(
        {
            from: appliedPeriod1.from.toISOString(),
            to: appliedPeriod1.to.toISOString(),
            include_histogram: false,
            include_hrv: true,
        },
        { skip: !shouldFetchData },
    );

    const {
        data: histogramData1,
        isLoading: histogramLoading1,
        isFetching: histogramFetching1,
        error: histogramError1,
    } = useGetHistogramQuery(
        {
            from: appliedPeriod1.from.toISOString(),
            to: appliedPeriod1.to.toISOString(),
            bins_count: 25,
        },
        { skip: !shouldFetchData || !visibleCharts.histogram },
    );

    const {
        isLoading: differentialHistogramLoading1,
        isFetching: differentialHistogramFetching1,
        error: differentialHistogramError1,
    } = useGetDifferentialHistogramQuery(
        {
            from: appliedPeriod1.from.toISOString(),
            to: appliedPeriod1.to.toISOString(),
            bins_count: 25,
        },
        { skip: !shouldFetchData || !visibleCharts.differentialHistogram },
    );

    const {
        isLoading: scatterplotLoading1,
        isFetching: scatterplotFetching1,
        error: scatterplotError1,
    } = useGetScatterplotQuery(
        {
            from: appliedPeriod1.from.toISOString(),
            to: appliedPeriod1.to.toISOString(),
        },
        { skip: !shouldFetchData || !visibleCharts.scatterplot },
    );

    const {
        data: trendsData1,
        isLoading: trendsLoading1,
        isFetching: trendsFetching1,
        error: trendsError1,
    } = useGetTrendsQuery(
        {
            from: appliedPeriod1.from.toISOString(),
            to: appliedPeriod1.to.toISOString(),
            window_size: 5,
        },
        { skip: !shouldFetchData || !visibleCharts.trends },
    );

    // API запросы для второго периода
    const {
        data: statisticsData2,
        isLoading: statisticsLoading2,
        isFetching: statisticsFetching2,
        error: statisticsError2,
    } = useGetRRStatisticsQuery(
        {
            from: appliedPeriod2.from.toISOString(),
            to: appliedPeriod2.to.toISOString(),
            include_histogram: false,
            include_hrv: true,
        },
        { skip: !shouldFetchData },
    );

    const {
        data: histogramData2,
        isLoading: histogramLoading2,
        isFetching: histogramFetching2,
        error: histogramError2,
    } = useGetHistogramQuery(
        {
            from: appliedPeriod2.from.toISOString(),
            to: appliedPeriod2.to.toISOString(),
            bins_count: 25,
        },
        { skip: !shouldFetchData || !visibleCharts.histogram },
    );

    const {
        isLoading: differentialHistogramLoading2,
        isFetching: differentialHistogramFetching2,
        error: differentialHistogramError2,
    } = useGetDifferentialHistogramQuery(
        {
            from: appliedPeriod2.from.toISOString(),
            to: appliedPeriod2.to.toISOString(),
            bins_count: 25,
        },
        { skip: !shouldFetchData || !visibleCharts.differentialHistogram },
    );

    const {
        isLoading: scatterplotLoading2,
        isFetching: scatterplotFetching2,
        error: scatterplotError2,
    } = useGetScatterplotQuery(
        {
            from: appliedPeriod2.from.toISOString(),
            to: appliedPeriod2.to.toISOString(),
        },
        { skip: !shouldFetchData || !visibleCharts.scatterplot },
    );

    const {
        data: trendsData2,
        isLoading: trendsLoading2,
        isFetching: trendsFetching2,
        error: trendsError2,
    } = useGetTrendsQuery(
        {
            from: appliedPeriod2.from.toISOString(),
            to: appliedPeriod2.to.toISOString(),
            window_size: 5,
        },
        { skip: !shouldFetchData || !visibleCharts.trends },
    );

    const isLoading =
        statisticsLoading1 ||
        histogramLoading1 ||
        differentialHistogramLoading1 ||
        scatterplotLoading1 ||
        trendsLoading1 ||
        statisticsLoading2 ||
        histogramLoading2 ||
        differentialHistogramLoading2 ||
        scatterplotLoading2 ||
        trendsLoading2;

    const isFetching =
        statisticsFetching1 ||
        histogramFetching1 ||
        differentialHistogramFetching1 ||
        scatterplotFetching1 ||
        trendsFetching1 ||
        statisticsFetching2 ||
        histogramFetching2 ||
        differentialHistogramFetching2 ||
        scatterplotFetching2 ||
        trendsFetching2;

    // Вычисляем сравнительную статистику
    const comparisonStats = useMemo(() => {
        if (!statisticsData1?.summary || !statisticsData2?.summary) {
            return null;
        }

        const period1Stats = statisticsData1.summary;
        const period2Stats = statisticsData2.summary;
        const period1HRV = statisticsData1.hrv_metrics;
        const period2HRV = statisticsData2.hrv_metrics;

        return {
            meanChange:
                ((period2Stats.mean - period1Stats.mean) / period1Stats.mean) *
                100,
            stdDevChange:
                ((period2Stats.std_dev - period1Stats.std_dev) /
                    period1Stats.std_dev) *
                100,
            countChange:
                ((period2Stats.count - period1Stats.count) /
                    period1Stats.count) *
                100,
            minChange:
                period1Stats.min && period2Stats.min
                    ? ((period2Stats.min - period1Stats.min) /
                          period1Stats.min) *
                      100
                    : null,
            maxChange:
                period1Stats.max && period2Stats.max
                    ? ((period2Stats.max - period1Stats.max) /
                          period1Stats.max) *
                      100
                    : null,
            rmssdChange:
                period1HRV?.rmssd && period2HRV?.rmssd
                    ? ((period2HRV.rmssd - period1HRV.rmssd) /
                          period1HRV.rmssd) *
                      100
                    : null,
            sdnnChange:
                period1HRV?.sdnn && period2HRV?.sdnn
                    ? ((period2HRV.sdnn - period1HRV.sdnn) / period1HRV.sdnn) *
                      100
                    : null,
            pnn50Change:
                period1HRV?.pnn50 !== undefined &&
                period2HRV?.pnn50 !== undefined
                    ? ((period2HRV.pnn50 - period1HRV.pnn50) /
                          (period1HRV.pnn50 || 1)) *
                      100
                    : null,
            triangularChange:
                period1HRV?.triangular_index && period2HRV?.triangular_index
                    ? ((period2HRV.triangular_index -
                          period1HRV.triangular_index) /
                          period1HRV.triangular_index) *
                      100
                    : null,
            period1: period1Stats,
            period2: period2Stats,
            period1HRV,
            period2HRV,
        };
    }, [statisticsData1, statisticsData2]);

    const handleCompareWithLoading = () => {
        setIsComparing(true);
        handleCompare();
        // Имитируем процесс сравнения
        setTimeout(() => setIsComparing(false), 1000);
    };

    return (
        <div className="space-y-6">
            {/* Заголовок */}
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">
                        Сравнительный анализ параметров
                    </h2>
                    <p className="text-muted-foreground">
                        Выберите два периода времени для сравнения показателей
                        RR интервалов
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

            {/* Выбор периодов для сравнения */}
            <PeriodSelector
                period1={period1}
                period2={period2}
                onPeriod1Change={handlePeriod1Change}
                onPeriod2Change={handlePeriod2Change}
            />

            {/* Кнопка сравнения */}
            <div className="flex justify-center gap-4">
                <Button
                    onClick={handleCompareWithLoading}
                    disabled={isComparing || isFetching}
                    size="lg"
                    className="px-8"
                >
                    {isComparing || isFetching ? (
                        <>
                            <LoadingSpinner className="h-4 w-4 mr-2" />
                            Сравнение...
                        </>
                    ) : (
                        <>
                            <Play className="h-4 w-4 mr-2" />
                            {hasUnappliedChanges
                                ? 'Сравнить периоды'
                                : 'Данные актуальны'}
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

            {/* Результаты сравнения - показываем только если данные должны загружаться */}
            {shouldFetchData && (
                <div className="space-y-6">
                    {/* Ошибки */}
                    {(statisticsError1 ||
                        statisticsError2 ||
                        histogramError1 ||
                        histogramError2 ||
                        differentialHistogramError1 ||
                        differentialHistogramError2 ||
                        scatterplotError1 ||
                        scatterplotError2 ||
                        trendsError1 ||
                        trendsError2) && (
                        <ErrorAlert
                            title="Ошибки при загрузке данных"
                            error="Произошли ошибки при загрузке одного или нескольких наборов данных"
                            getErrorMessage={getErrorMessage}
                        />
                    )}

                    {/* Статистика периодов */}
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        <PeriodStatistics
                            title="Статистика периода 1"
                            colorClass="border-blue-200"
                            statistics={statisticsData1?.summary}
                            hrvMetrics={statisticsData1?.hrv_metrics}
                            isLoading={statisticsLoading1}
                        />
                        <PeriodStatistics
                            title="Статистика периода 2"
                            colorClass="border-red-200"
                            statistics={statisticsData2?.summary}
                            hrvMetrics={statisticsData2?.hrv_metrics}
                            isLoading={statisticsLoading2}
                        />
                    </div>

                    {/* Сравнительная статистика */}
                    {comparisonStats && (
                        <ComparisonStats comparisonStats={comparisonStats} />
                    )}

                    {/* Сравнительные графики */}
                    <ComparisonCharts
                        appliedPeriod1={appliedPeriod1}
                        appliedPeriod2={appliedPeriod2}
                        visibleCharts={visibleCharts}
                        onToggleChart={toggleChart}
                        histogramData1={histogramData1?.histogram}
                        histogramData2={histogramData2?.histogram}
                        trendsData1={trendsData1}
                        trendsData2={trendsData2}
                        loadingStates={{
                            histogram1: histogramLoading1,
                            histogram2: histogramLoading2,
                            differentialHistogram1: differentialHistogramLoading1,
                            differentialHistogram2: differentialHistogramLoading2,
                            scatterplot1: scatterplotLoading1,
                            scatterplot2: scatterplotLoading2,
                            trends1: trendsLoading1,
                            trends2: trendsLoading2,
                        }}
                    />
                </div>
            )}
        </div>
    );
}
