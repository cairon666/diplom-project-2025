import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from 'ui/card';
import {
    BarChart3,
    TrendingUp,
    Activity,
    ArrowUpDown,
} from 'lucide-react';
import type { RRStatisticalSummary, HRVMetrics } from '../api/rrAnalyticsApi';

interface ComparisonStats {
    meanChange: number;
    stdDevChange: number;
    countChange: number;
    minChange: number | null;
    maxChange: number | null;
    rmssdChange: number | null;
    sdnnChange: number | null;
    pnn50Change: number | null;
    triangularChange: number | null;
    period1: RRStatisticalSummary;
    period2: RRStatisticalSummary;
    period1HRV?: HRVMetrics;
    period2HRV?: HRVMetrics;
}

interface ComparisonStatsProps {
    comparisonStats: ComparisonStats;
}

export function ComparisonStats({ comparisonStats }: ComparisonStatsProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <ArrowUpDown className="h-5 w-5" />
                    Сравнительный анализ изменений
                </CardTitle>
                <CardDescription>
                    Изменения от Периода 1 к Периоду 2 (слева направо)
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    {/* Изменение среднего RR */}
                    <Card className="border-slate-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                Среднее RR
                            </CardTitle>
                            <Activity className="h-4 w-4 text-muted-foreground" />
                        </CardHeader>
                        <CardContent>
                            <div
                                className={`text-2xl font-bold ${
                                    comparisonStats.meanChange >= 0
                                        ? 'text-green-600'
                                        : 'text-red-600'
                                }`}
                            >
                                {comparisonStats.meanChange >= 0 ? '+' : ''}
                                {comparisonStats.meanChange.toFixed(1)}%
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {comparisonStats.period1.mean.toFixed(1)} →{' '}
                                {comparisonStats.period2.mean.toFixed(1)} мс
                            </p>
                        </CardContent>
                    </Card>

                    {/* Изменение вариабельности */}
                    <Card className="border-slate-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                Вариабельность
                            </CardTitle>
                            <TrendingUp className="h-4 w-4 text-muted-foreground" />
                        </CardHeader>
                        <CardContent>
                            <div
                                className={`text-2xl font-bold ${
                                    comparisonStats.stdDevChange >= 0
                                        ? 'text-green-600'
                                        : 'text-red-600'
                                }`}
                            >
                                {comparisonStats.stdDevChange >= 0 ? '+' : ''}
                                {comparisonStats.stdDevChange.toFixed(1)}%
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {comparisonStats.period1.std_dev.toFixed(1)} →{' '}
                                {comparisonStats.period2.std_dev.toFixed(1)} мс
                            </p>
                        </CardContent>
                    </Card>

                    {/* Изменение диапазона (минимум) */}
                    {comparisonStats.minChange !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    Минимум RR
                                </CardTitle>
                                <TrendingUp className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                {(() => {
                                    const absoluteChange =
                                        comparisonStats.period2.min -
                                        comparisonStats.period1.min;
                                    return (
                                        <>
                                            <div
                                                className={`text-2xl font-bold ${
                                                    absoluteChange >= 0
                                                        ? 'text-green-600'
                                                        : 'text-red-600'
                                                }`}
                                            >
                                                {absoluteChange >= 0 ? '+' : ''}
                                                {absoluteChange.toFixed(0)} мс
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {comparisonStats.period1.min} →{' '}
                                                {comparisonStats.period2.min} мс
                                            </p>
                                        </>
                                    );
                                })()}
                            </CardContent>
                        </Card>
                    )}

                    {/* Изменение диапазона (максимум) */}
                    {comparisonStats.maxChange !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    Максимум RR
                                </CardTitle>
                                <TrendingUp className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                {(() => {
                                    const absoluteChange =
                                        comparisonStats.period2.max -
                                        comparisonStats.period1.max;
                                    return (
                                        <>
                                            <div
                                                className={`text-2xl font-bold ${
                                                    absoluteChange >= 0
                                                        ? 'text-green-600'
                                                        : 'text-red-600'
                                                }`}
                                            >
                                                {absoluteChange >= 0 ? '+' : ''}
                                                {absoluteChange.toFixed(0)} мс
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {comparisonStats.period1.max} →{' '}
                                                {comparisonStats.period2.max} мс
                                            </p>
                                        </>
                                    );
                                })()}
                            </CardContent>
                        </Card>
                    )}

                    {/* Изменение количества данных */}
                    <Card className="border-slate-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                Количество
                            </CardTitle>
                            <BarChart3 className="h-4 w-4 text-muted-foreground" />
                        </CardHeader>
                        <CardContent>
                            <div
                                className={`text-2xl font-bold ${
                                    comparisonStats.countChange >= 0
                                        ? 'text-green-600'
                                        : 'text-red-600'
                                }`}
                            >
                                {comparisonStats.countChange >= 0 ? '+' : ''}
                                {comparisonStats.countChange.toFixed(1)}%
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {comparisonStats.period1.count} →{' '}
                                {comparisonStats.period2.count} измерений
                            </p>
                        </CardContent>
                    </Card>

                    {/* Изменение BPM */}
                    <Card className="border-slate-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                ЧСС (BPM)
                            </CardTitle>
                            <Activity className="h-4 w-4 text-muted-foreground" />
                        </CardHeader>
                        <CardContent>
                            {(() => {
                                const currentBPM = Math.round(
                                    60000 / comparisonStats.period2.mean,
                                );
                                const previousBPM = Math.round(
                                    60000 / comparisonStats.period1.mean,
                                );
                                const bpmChange =
                                    ((currentBPM - previousBPM) / previousBPM) *
                                    100;

                                return (
                                    <>
                                        <div
                                            className={`text-2xl font-bold ${
                                                bpmChange >= 0
                                                    ? 'text-green-600'
                                                    : 'text-red-600'
                                            }`}
                                        >
                                            {bpmChange >= 0 ? '+' : ''}
                                            {bpmChange.toFixed(1)}%
                                        </div>
                                        <p className="text-xs text-muted-foreground">
                                            {previousBPM} → {currentBPM} уд/мин
                                        </p>
                                    </>
                                );
                            })()}
                        </CardContent>
                    </Card>

                    {/* HRV метрики */}
                    {comparisonStats.rmssdChange !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    RMSSD
                                </CardTitle>
                                <CardDescription className="text-xs">
                                    HRV
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div
                                    className={`text-2xl font-bold ${
                                        comparisonStats.rmssdChange >= 0
                                            ? 'text-green-600'
                                            : 'text-red-600'
                                    }`}
                                >
                                    {comparisonStats.rmssdChange >= 0 ? '+' : ''}
                                    {comparisonStats.rmssdChange.toFixed(1)}%
                                </div>
                                <p className="text-xs text-muted-foreground">
                                    {comparisonStats.period1HRV?.rmssd?.toFixed(1)} →{' '}
                                    {comparisonStats.period2HRV?.rmssd?.toFixed(1)} мс
                                </p>
                            </CardContent>
                        </Card>
                    )}

                    {comparisonStats.sdnnChange !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    SDNN
                                </CardTitle>
                                <CardDescription className="text-xs">
                                    HRV
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div
                                    className={`text-2xl font-bold ${
                                        comparisonStats.sdnnChange >= 0
                                            ? 'text-green-600'
                                            : 'text-red-600'
                                    }`}
                                >
                                    {comparisonStats.sdnnChange >= 0 ? '+' : ''}
                                    {comparisonStats.sdnnChange.toFixed(1)}%
                                </div>
                                <p className="text-xs text-muted-foreground">
                                    {comparisonStats.period1HRV?.sdnn?.toFixed(1)} →{' '}
                                    {comparisonStats.period2HRV?.sdnn?.toFixed(1)} мс
                                </p>
                            </CardContent>
                        </Card>
                    )}

                    {comparisonStats.pnn50Change !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    pNN50
                                </CardTitle>
                                <CardDescription className="text-xs">
                                    HRV
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                {(() => {
                                    const absoluteChange =
                                        comparisonStats.period2HRV?.pnn50 &&
                                        comparisonStats.period1HRV?.pnn50
                                            ? comparisonStats.period2HRV.pnn50 -
                                              comparisonStats.period1HRV.pnn50
                                            : 0;
                                    return (
                                        <>
                                            <div
                                                className={`text-2xl font-bold ${
                                                    absoluteChange >= 0
                                                        ? 'text-green-600'
                                                        : 'text-red-600'
                                                }`}
                                            >
                                                {absoluteChange >= 0 ? '+' : ''}
                                                {absoluteChange.toFixed(1)}%
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {comparisonStats.period1HRV?.pnn50?.toFixed(1)} →{' '}
                                                {comparisonStats.period2HRV?.pnn50?.toFixed(1)}%
                                            </p>
                                        </>
                                    );
                                })()}
                            </CardContent>
                        </Card>
                    )}

                    {comparisonStats.triangularChange !== null && (
                        <Card className="border-slate-200">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">
                                    Треугольный индекс
                                </CardTitle>
                                <CardDescription className="text-xs">
                                    HRV
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                {(() => {
                                    const absoluteChange =
                                        comparisonStats.period2HRV?.triangular_index &&
                                        comparisonStats.period1HRV?.triangular_index
                                            ? comparisonStats.period2HRV.triangular_index -
                                              comparisonStats.period1HRV.triangular_index
                                            : 0;
                                    return (
                                        <>
                                            <div
                                                className={`text-2xl font-bold ${
                                                    absoluteChange >= 0
                                                        ? 'text-green-600'
                                                        : 'text-red-600'
                                                }`}
                                            >
                                                {absoluteChange >= 0 ? '+' : ''}
                                                {absoluteChange.toFixed(2)}
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {comparisonStats.period1HRV?.triangular_index?.toFixed(2)} →{' '}
                                                {comparisonStats.period2HRV?.triangular_index?.toFixed(2)}
                                            </p>
                                        </>
                                    );
                                })()}
                            </CardContent>
                        </Card>
                    )}
                </div>
            </CardContent>
        </Card>
    );
} 