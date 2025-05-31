import { useMemo } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ReferenceLine } from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from 'ui/card';
import { Badge } from 'ui/badge';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import { RRTrendAnalysis } from '../api/rrAnalyticsApi';

interface RRTrendsChartProps {
    data: RRTrendAnalysis | null | undefined;
    className?: string;
}

export function RRTrendsChart({ data, className }: RRTrendsChartProps) {
    // Проверяем данные и предоставляем fallback
    if (!data || !data.trend_points || !Array.isArray(data.trend_points) || data.trend_points.length === 0) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle>Анализ трендов R-R интервалов</CardTitle>
                    <CardDescription>
                        Нет данных для анализа трендов
                    </CardDescription>
                </CardHeader>
                <CardContent className="py-8 text-center">
                    <p className="text-muted-foreground">
                        Недостаточно данных для построения трендов. Попробуйте выбрать больший временной диапазон (минимум 15-30 минут).
                    </p>
                </CardContent>
            </Card>
        );
    }

    const chartData = useMemo(() => {
        if (!data?.trend_points) return [];
        
        const result = data.trend_points.map(point => ({
            time: new Date(point.time).toLocaleTimeString('ru-RU', { 
                hour: '2-digit', 
                minute: '2-digit' 
            }),
            fullTime: point.time,
            value: Math.round(point.value * 10) / 10,
            direction: point.direction,
        }));
        
        return result;
    }, [data.trend_points]);

    const getTrendIcon = (trend: string) => {
        switch (trend) {
            case 'increasing':
                return <TrendingUp className="h-4 w-4" />;
            case 'decreasing':
                return <TrendingDown className="h-4 w-4" />;
            default:
                return <Minus className="h-4 w-4" />;
        }
    };

    const getTrendBadge = (trend: string) => {
        switch (trend) {
            case 'increasing':
                return { label: 'Возрастающий', variant: 'default' as const, color: 'text-green-600' };
            case 'decreasing':
                return { label: 'Убывающий', variant: 'destructive' as const, color: 'text-red-600' };
            default:
                return { label: 'Стабильный', variant: 'secondary' as const, color: 'text-blue-600' };
        }
    };

    const trendBadge = getTrendBadge(data.overall_trend || 'stable');
    const averageValue = chartData.length > 0 
        ? chartData.reduce((sum, point) => sum + point.value, 0) / chartData.length 
        : 0;

    const formatTooltip = (value: number, name: string) => {
        if (name === 'value') {
            return [`${value} мс`, 'RR интервал'];
        }
        return [value, name];
    };

    const formatLabel = (label: string) => {
        return `Время: ${label}`;
    };

    return (
        <Card className={className}>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    Анализ трендов R-R интервалов
                    <div className={`${trendBadge.color}`}>
                        {getTrendIcon(data.overall_trend || 'stable')}
                    </div>
                    <span className="text-sm font-normal text-muted-foreground">
                        ({chartData.length} точек)
                    </span>
                </CardTitle>
                <CardDescription>
                    Динамика изменения R-R интервалов во времени. Период: {data.period || 'Не указан'}
                    {chartData.length > 0 && (
                        <span className="block text-xs mt-1">
                            Диапазон значений: {Math.min(...chartData.map(p => p.value)).toFixed(1)} - {Math.max(...chartData.map(p => p.value)).toFixed(1)} мс
                        </span>
                    )}
                </CardDescription>
            </CardHeader>
            <CardContent>
                {chartData.length < 2 && (
                    <div className="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                        <p className="text-sm text-yellow-800">
                            ⚠️ Недостаточно точек для построения линии тренда. Нужно минимум 2 точки. 
                            Увеличьте временной диапазон для получения более детального анализа.
                        </p>
                    </div>
                )}
                
                <div className="h-80">
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart
                            data={chartData}
                            margin={{
                                top: 20,
                                right: 30,
                                left: 20,
                                bottom: 20,
                            }}
                        >
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis 
                                dataKey="time"
                                fontSize={12}
                                angle={-45}
                                textAnchor="end"
                                height={60}
                            />
                            <YAxis 
                                domain={['dataMin - 10', 'dataMax + 10']}
                                label={{ value: 'RR (мс)', angle: -90, position: 'insideLeft' }}
                            />
                            <Tooltip 
                                formatter={formatTooltip}
                                labelFormatter={formatLabel}
                                contentStyle={{
                                    backgroundColor: 'hsl(var(--card))',
                                    border: '1px solid hsl(var(--border))',
                                    borderRadius: '6px',
                                }}
                            />
                            {/* Средняя линия */}
                            {averageValue > 0 && (
                                <ReferenceLine 
                                    y={averageValue} 
                                    stroke="hsl(var(--muted-foreground))" 
                                    strokeDasharray="5 5" 
                                    label="Среднее"
                                />
                            )}
                            <Line 
                                type="monotone" 
                                dataKey="value" 
                                stroke="#3b82f6"
                                strokeWidth={3}
                                dot={{ r: 4, fill: "#3b82f6", stroke: "#1d4ed8", strokeWidth: 1 }}
                                activeDot={{ r: 6, fill: "#1d4ed8", stroke: "#1e40af", strokeWidth: 2 }}
                                connectNulls={false}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
                
                {/* Статистика по трендам */}
                <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="text-center">
                        <Badge variant={trendBadge.variant} className="mb-2">
                            {trendBadge.label}
                        </Badge>
                        <div className="text-sm text-muted-foreground">Общий тренд</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold text-lg">{((data.correlation || 0) * 100).toFixed(1)}%</div>
                        <div className="text-sm text-muted-foreground">Корреляция</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold text-lg">{((data.trend_strength || 0) * 100).toFixed(1)}%</div>
                        <div className="text-sm text-muted-foreground">Сила тренда</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold text-lg">{averageValue.toFixed(1)} мс</div>
                        <div className="text-sm text-muted-foreground">Среднее значение</div>
                    </div>
                </div>

                {/* Сезонность */}
                {data.seasonality && Array.isArray(data.seasonality) && data.seasonality.length > 0 && (
                    <div className="mt-4">
                        <h4 className="text-sm font-medium mb-2">Сезонность (по часам)</h4>
                        <div className="grid grid-cols-12 gap-1">
                            {data.seasonality.map((value, index) => (
                                <div key={index} className="text-center">
                                    <div className="text-xs font-mono">{(value || 0).toFixed(0)}</div>
                                    <div className="text-xs text-muted-foreground">{index}ч</div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
