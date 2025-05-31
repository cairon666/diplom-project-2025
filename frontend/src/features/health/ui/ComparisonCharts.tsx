import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { LoadingOverlay } from 'ui/loading-overlay';
import {
    BarChart3,
    TrendingUp,
    Activity,
    ScatterChart,
} from 'lucide-react';
import { ChartSection } from './ChartSection';
import { RRHistogramChart } from './RRHistogramChart';
import { RRDifferentialHistogramChart } from './RRDifferentialHistogramChart';
import { RRScatterplotChart } from './RRScatterplotChart';
import { RRTrendsChart } from './RRTrendsChart';
import type { DateTimeRange } from './DateTimeRangeSelector';
import type { RRHistogramData, RRTrendAnalysis } from '../api/rrAnalyticsApi';
import type { ChartVisibilityState } from '../hooks/useChartVisibility';

interface ComparisonChartsProps {
    appliedPeriod1: DateTimeRange;
    appliedPeriod2: DateTimeRange;
    visibleCharts: ChartVisibilityState;
    onToggleChart: (chartName: keyof ChartVisibilityState) => void;
    histogramData1?: RRHistogramData;
    histogramData2?: RRHistogramData;
    trendsData1?: { trend_analysis: RRTrendAnalysis };
    trendsData2?: { trend_analysis: RRTrendAnalysis };
    loadingStates: {
        histogram1: boolean;
        histogram2: boolean;
        differentialHistogram1: boolean;
        differentialHistogram2: boolean;
        scatterplot1: boolean;
        scatterplot2: boolean;
        trends1: boolean;
        trends2: boolean;
    };
}

export function ComparisonCharts({
    appliedPeriod1,
    appliedPeriod2,
    visibleCharts,
    onToggleChart,
    histogramData1,
    histogramData2,
    trendsData1,
    trendsData2,
    loadingStates,
}: ComparisonChartsProps) {
    return (
        <div className="space-y-6">
            {/* Гистограммы */}
            <ChartSection
                chartName="histogram"
                title="Сравнение распределений RR интервалов"
                icon={BarChart3}
                isVisible={visibleCharts.histogram}
                isValid={true}
                onToggle={onToggleChart}
            >
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Гистограмма первого периода */}
                    <Card className="border-blue-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-blue-500" />
                                Период 1 - Распределение
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.histogram1}>
                                {histogramData1 ? (
                                    <RRHistogramChart
                                        data={histogramData1}
                                        className="w-full"
                                    />
                                ) : (
                                    <div className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для отображения
                                        </p>
                                    </div>
                                )}
                            </LoadingOverlay>
                        </CardContent>
                    </Card>

                    {/* Гистограмма второго периода */}
                    <Card className="border-red-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-red-500" />
                                Период 2 - Распределение
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.histogram2}>
                                {histogramData2 ? (
                                    <RRHistogramChart
                                        data={histogramData2}
                                        className="w-full"
                                    />
                                ) : (
                                    <div className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для отображения
                                        </p>
                                    </div>
                                )}
                            </LoadingOverlay>
                        </CardContent>
                    </Card>
                </div>
            </ChartSection>

            {/* Дифференциальные гистограммы */}
            <ChartSection
                chartName="differentialHistogram"
                title="Сравнение дифференциальных гистограмм"
                icon={Activity}
                isVisible={visibleCharts.differentialHistogram}
                isValid={true}
                onToggle={onToggleChart}
            >
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Дифференциальная гистограмма первого периода */}
                    <Card className="border-blue-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-blue-500" />
                                Период 1 - Дифференциальная гистограмма
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.differentialHistogram1}>
                                <RRDifferentialHistogramChart
                                    from={appliedPeriod1.from.toISOString()}
                                    to={appliedPeriod1.to.toISOString()}
                                    binsCount={25}
                                />
                            </LoadingOverlay>
                        </CardContent>
                    </Card>

                    {/* Дифференциальная гистограмма второго периода */}
                    <Card className="border-red-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-red-500" />
                                Период 2 - Дифференциальная гистограмма
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.differentialHistogram2}>
                                <RRDifferentialHistogramChart
                                    from={appliedPeriod2.from.toISOString()}
                                    to={appliedPeriod2.to.toISOString()}
                                    binsCount={25}
                                />
                            </LoadingOverlay>
                        </CardContent>
                    </Card>
                </div>
            </ChartSection>

            {/* Скаттерограммы (диаграммы Пуанкаре) */}
            <ChartSection
                chartName="scatterplot"
                title="Сравнение скаттерограмм (диаграммы Пуанкаре)"
                icon={ScatterChart}
                isVisible={visibleCharts.scatterplot}
                isValid={true}
                onToggle={onToggleChart}
            >
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Скаттерограмма первого периода */}
                    <Card className="border-blue-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-blue-500" />
                                Период 1 - Скаттерограмма
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.scatterplot1}>
                                <RRScatterplotChart
                                    from={appliedPeriod1.from.toISOString()}
                                    to={appliedPeriod1.to.toISOString()}
                                />
                            </LoadingOverlay>
                        </CardContent>
                    </Card>

                    {/* Скаттерограмма второго периода */}
                    <Card className="border-red-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-red-500" />
                                Период 2 - Скаттерограмма
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.scatterplot2}>
                                <RRScatterplotChart
                                    from={appliedPeriod2.from.toISOString()}
                                    to={appliedPeriod2.to.toISOString()}
                                />
                            </LoadingOverlay>
                        </CardContent>
                    </Card>
                </div>
            </ChartSection>

            {/* Тренды */}
            <ChartSection
                chartName="trends"
                title="Сравнение трендов изменений"
                icon={TrendingUp}
                isVisible={visibleCharts.trends}
                isValid={true}
                onToggle={onToggleChart}
            >
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Тренды первого периода */}
                    <Card className="border-blue-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-blue-500" />
                                Период 1 - Тренды
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.trends1}>
                                {trendsData1?.trend_analysis ? (
                                    <RRTrendsChart
                                        data={trendsData1.trend_analysis}
                                        className="w-full"
                                    />
                                ) : (
                                    <div className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для анализа трендов
                                        </p>
                                    </div>
                                )}
                            </LoadingOverlay>
                        </CardContent>
                    </Card>

                    {/* Тренды второго периода */}
                    <Card className="border-red-200">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-lg">
                                <div className="w-4 h-4 rounded-full bg-red-500" />
                                Период 2 - Тренды
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LoadingOverlay isLoading={loadingStates.trends2}>
                                {trendsData2?.trend_analysis ? (
                                    <RRTrendsChart
                                        data={trendsData2.trend_analysis}
                                        className="w-full"
                                    />
                                ) : (
                                    <div className="py-8 text-center">
                                        <p className="text-muted-foreground">
                                            Нет данных для анализа трендов
                                        </p>
                                    </div>
                                )}
                            </LoadingOverlay>
                        </CardContent>
                    </Card>
                </div>
            </ChartSection>
        </div>
    );
} 