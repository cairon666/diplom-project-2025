import { Card, CardContent, CardDescription, CardHeader, CardTitle } from 'ui/card';
import { Badge } from 'ui/badge';
import { Skeleton } from 'ui/skeleton';
import { Activity, Heart, TrendingUp, BarChart } from 'lucide-react';
import { RRStatisticalSummary, HRVMetrics } from '../api/rrAnalyticsApi';

interface RRStatisticsProps {
    statistics: RRStatisticalSummary | null | undefined;
    hrvMetrics?: HRVMetrics | null | undefined;
    className?: string;
    isLoading?: boolean;
    compact?: boolean; // Компактный режим: максимум 2 колонки для сравнительного анализа
}

export function RRStatistics({ statistics, hrvMetrics, className, isLoading = false, compact = false }: RRStatisticsProps) {
    // Компонент скелета для загрузки
    const SkeletonCard = () => (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <Skeleton className="h-4 w-20" />
                <Skeleton className="h-4 w-4" />
            </CardHeader>
            <CardContent>
                <Skeleton className="h-8 w-16 mb-1" />
                <Skeleton className="h-3 w-24" />
            </CardContent>
        </Card>
    );

    // Определяем CSS классы для сетки в зависимости от режима
    const gridClasses = compact 
        ? "grid grid-cols-1 md:grid-cols-2 gap-4" 
        : "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4";

    // Показываем скелет при загрузке
    if (isLoading) {
        return (
            <div className={`${gridClasses} ${className}`}>
                <SkeletonCard />
                <SkeletonCard />
                <SkeletonCard />
                <SkeletonCard />
                {hrvMetrics && (
                    <>
                        <SkeletonCard />
                        <SkeletonCard />
                        <SkeletonCard />
                        <SkeletonCard />
                    </>
                )}
            </div>
        );
    }

    // Проверяем, есть ли валидные данные в statistics
    const hasValidStatistics = statistics && 
        (statistics.mean !== undefined && statistics.mean !== null && statistics.mean > 0) &&
        (statistics.count !== undefined && statistics.count !== null && statistics.count > 0);

    // Проверяем, есть ли валидные HRV метрики
    const hasValidHRVMetrics = hrvMetrics && 
        (hrvMetrics.rmssd !== undefined && hrvMetrics.rmssd !== null && hrvMetrics.rmssd >= 0);

    // Если нет данных вообще
    if (!statistics) {
        return (
            <div className={`text-center py-8 ${className}`}>
                <p className="text-muted-foreground">Нет данных для отображения статистики</p>
            </div>
        );
    }

    // Если данные пустые
    if (!hasValidStatistics) {
        return (
            <div className={`text-center py-6 ${className}`}>
                <p className="text-muted-foreground">Статистические данные пусты</p>
                <p className="text-xs text-muted-foreground mt-1">
                    Загрузите R-R интервалы для получения статистики
                </p>
            </div>
        );
    }

    const formatNumber = (value: number | undefined, decimals = 1) => {
        return (value || 0).toFixed(decimals);
    };

    const getQualityBadge = (count: number) => {
        if (count < 100) return { label: 'Мало данных', variant: 'secondary' as const };
        if (count < 500) return { label: 'Достаточно', variant: 'default' as const };
        return { label: 'Отлично', variant: 'default' as const };
    };

    const qualityBadge = getQualityBadge(statistics.count || 0);

    return (
        <div className={`${gridClasses} ${className}`}>
            {/* Основная статистика */}
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Среднее RR</CardTitle>
                    <Heart className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(statistics.mean)} мс</div>
                    <p className="text-xs text-muted-foreground">
                        ≈ {Math.round(60000 / (statistics.mean || 1))} BPM
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Вариабельность</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(statistics.std_dev)} мс</div>
                    <p className="text-xs text-muted-foreground">
                        Стандартное отклонение
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Диапазон</CardTitle>
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">
                        {statistics.min || 0} - {statistics.max || 0}
                    </div>
                    <p className="text-xs text-muted-foreground">
                        мс (мин - макс)
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Количество</CardTitle>
                    <BarChart className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{(statistics.count || 0).toLocaleString()}</div>
                    <div className="flex items-center gap-2 mt-1">
                        <Badge variant={qualityBadge.variant} className="text-xs">
                            {qualityBadge.label}
                        </Badge>
                    </div>
                </CardContent>
            </Card>

            {/* HRV метрики */}
            {hasValidHRVMetrics && (
                <>
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">RMSSD</CardTitle>
                            <CardDescription className="text-xs">HRV</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{formatNumber(hrvMetrics.rmssd)} мс</div>
                            <p className="text-xs text-muted-foreground">
                                Парасимпатическая активность
                            </p>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">SDNN</CardTitle>
                            <CardDescription className="text-xs">HRV</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{formatNumber(hrvMetrics.sdnn)} мс</div>
                            <p className="text-xs text-muted-foreground">
                                Общая вариабельность
                            </p>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">pNN50</CardTitle>
                            <CardDescription className="text-xs">HRV</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{formatNumber(hrvMetrics.pnn50)}%</div>
                            <p className="text-xs text-muted-foreground">
                                Процент соседних интервалов
                            </p>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">Треугольный индекс</CardTitle>
                            <CardDescription className="text-xs">HRV</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{formatNumber(hrvMetrics.triangular_index, 2)}</div>
                            <p className="text-xs text-muted-foreground">
                                Геометрический показатель
                            </p>
                        </CardContent>
                    </Card>
                </>
            )}
        </div>
    );
} 