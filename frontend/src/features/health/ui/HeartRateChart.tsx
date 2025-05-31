import { useState, useMemo } from 'react';
import { 
    LineChart, 
    Line, 
    XAxis, 
    YAxis, 
    CartesianGrid, 
    Tooltip, 
    ResponsiveContainer,
    ReferenceLine,
    TooltipProps
} from 'recharts';
import { format, subMinutes } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { LoadingSpinner } from 'ui/loading-spinner';
import { LuHeart, LuCalendar } from 'react-icons/lu';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import type { SerializedError } from '@reduxjs/toolkit';

import { useGetRRIntervalsQuery, RRIntervalsResponse, RRInterval } from '../api/rrIntervalsApi';
import { DateTimeRange } from './DateTimeRangeSelector';

interface ChartDataPoint {
    timestamp: string;
    bpm: number;
    formattedTime: string;
    rr_interval_ms: number;
    is_valid: boolean;
}

interface HeartRateChartProps {
    dateRange?: DateTimeRange;
    onChange?: (dateRange: DateTimeRange) => void;
    externalData?: RRIntervalsResponse;
    externalLoading?: boolean;
    externalError?: FetchBaseQueryError | SerializedError | undefined;
}

export function HeartRateChart({ 
    dateRange: externalDateRange, 
    onChange, 
    externalData,
    externalLoading,
    externalError,
}: HeartRateChartProps = {}) {
    // По умолчанию последние 5 минут для детального анализа
    const [internalDateRange, setInternalDateRange] = useState<DateTimeRange>({
        from: subMinutes(new Date(), 5),
        to: new Date(),
    });

    // Используем внешний dateRange если он передан, иначе внутренний
    const dateRange = externalDateRange || internalDateRange;
    
    // Функция для обновления dateRange
    const handleDateRangeChange = (newDateRange: DateTimeRange) => {
        if (onChange) {
            onChange(newDateRange);
        } else {
            setInternalDateRange(newDateRange);
        }
    };

    // Преобразуем диапазон дат в RFC3339 формат для API
    const queryParams = useMemo(() => ({
        from: dateRange.from.toISOString(),
        to: dateRange.to.toISOString(),
        // device_id можно добавить позже для фильтрации по устройству
    }), [dateRange]);

    const { data: internalData, isLoading: internalLoading, error: internalError } = useGetRRIntervalsQuery(queryParams, {
        skip: !!externalData // Пропускаем запрос если данные переданы извне
    });

    // Используем внешние данные если они есть, иначе внутренние
    const data = externalData || internalData;
    const isLoading = externalLoading !== undefined ? externalLoading : internalLoading;
    const error = externalError || internalError;

    // Преобразуем RR интервалы в данные для графика
    const chartData: ChartDataPoint[] = useMemo(() => {
        if (!data?.intervals) return [];

        return data.intervals
            .filter((interval: RRInterval) => interval.is_valid) // Показываем только валидные интервалы
            .map((interval: RRInterval) => ({
                timestamp: interval.created_at,
                bpm: interval.bpm,
                rr_interval_ms: interval.rr_interval_ms,
                is_valid: interval.is_valid,
                formattedTime: format(new Date(interval.created_at), 'HH:mm:ss', { locale: ru }),
            }))
            .sort((a: ChartDataPoint, b: ChartDataPoint) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
    }, [data]);

    // Вычисляем статистику
    const stats = useMemo(() => {
        if (chartData.length === 0) return { avg: 0, max: 0, min: 0 };

        const bpmValues = chartData.map(point => point.bpm);
        return {
            avg: Math.round(bpmValues.reduce((sum, bpm) => sum + bpm, 0) / bpmValues.length),
            max: Math.max(...bpmValues),
            min: Math.min(...bpmValues),
        };
    }, [chartData]);

    // Определяем интервал для оси X на основе количества данных (для детального просмотра показываем больше меток)
    const getTickInterval = (): number => {
        const dataLength = chartData.length;
        if (dataLength <= 10) return 1;
        if (dataLength <= 50) return Math.ceil(dataLength / 10);
        if (dataLength <= 100) return Math.ceil(dataLength / 15);
        return Math.ceil(dataLength / 20);
    };

    const CustomTooltip = ({ active, payload }: TooltipProps<number, string>) => {
        if (active && payload && payload.length > 0) {
            const data = payload[0].payload as ChartDataPoint;
            return (
                <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-lg">
                    <div className="flex items-center gap-2 mb-1">
                        <LuHeart className="w-4 h-4 text-red-500" />
                        <p className="font-medium">
                            {format(new Date(data.timestamp), 'dd MMM yyyy, HH:mm:ss.SSS', { locale: ru })}
                        </p>
                    </div>
                    <p className="text-red-600 font-semibold">
                        {data.bpm} уд/мин
                    </p>
                    <p className="text-gray-600 text-sm">
                        R-R интервал: {data.rr_interval_ms} мс
                    </p>
                </div>
            );
        }
        return null;
    };

    if (isLoading) {
        return (
            <div className="space-y-4">
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <LuHeart className="w-5 h-5 text-red-500" />
                            Детальный анализ пульса (R-R интервалы)
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-center h-64">
                            <LoadingSpinner />
                            <span className="ml-2">Загрузка данных...</span>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    if (error) {
        return (
            <div className="space-y-4">
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <LuHeart className="w-5 h-5 text-red-500" />
                            Детальный анализ пульса (R-R интервалы)
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-center h-64 text-muted-foreground">
                            <div className="text-center">
                                <LuCalendar className="w-8 h-8 mx-auto mb-2 opacity-50" />
                                <p>Ошибка загрузки данных</p>
                                <p className="text-sm">Проверьте подключение и попробуйте еще раз</p>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    if (chartData.length === 0) {
        return (
            <div className="space-y-4">
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <LuHeart className="w-5 h-5 text-red-500" />
                            Детальный анализ пульса (R-R интервалы)
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-center h-64 text-muted-foreground">
                            <div className="text-center">
                                <LuHeart className="w-8 h-8 mx-auto mb-2 opacity-50" />
                                <p>Нет данных за выбранный период</p>
                                <p className="text-sm">Загрузите R-R интервалы для анализа</p>
                                <p className="text-xs text-muted-foreground mt-2">
                                    Используйте компонент загрузки для тестовых данных
                                </p>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle className="flex items-center gap-2">
                            <LuHeart className="w-5 h-5 text-red-500" />
                            Детальный анализ пульса (R-R интервалы)
                            <span className="text-sm font-normal text-muted-foreground">
                                {data?.total_count} измерений, {data?.valid_count} валидных
                            </span>
                        </CardTitle>
                    </div>
                    
                    {/* Статистика */}
                    <div className="grid grid-cols-3 gap-4 mt-4">
                        <div className="text-center">
                            <p className="text-2xl font-bold text-red-600">{stats.avg}</p>
                            <p className="text-xs text-muted-foreground">Средний</p>
                        </div>
                        <div className="text-center">
                            <p className="text-2xl font-bold text-orange-600">{stats.max}</p>
                            <p className="text-xs text-muted-foreground">Максимум</p>
                        </div>
                        <div className="text-center">
                            <p className="text-2xl font-bold text-blue-600">{stats.min}</p>
                            <p className="text-xs text-muted-foreground">Минимум</p>
                        </div>
                    </div>
                </CardHeader>
                
                <CardContent>
                    <div className="h-96">
                        <ResponsiveContainer width="100%" height="100%">
                            <LineChart data={chartData}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                                <XAxis 
                                    dataKey="formattedTime"
                                    tick={{ fontSize: 10 }}
                                    tickLine={{ stroke: '#e0e0e0' }}
                                    axisLine={{ stroke: '#e0e0e0' }}
                                    interval={getTickInterval()}
                                    angle={-45}
                                    textAnchor="end"
                                    height={80}
                                    label={{ 
                                        value: 'Время (ЧЧ:ММ:СС)', 
                                        position: 'insideBottom', 
                                        offset: -5,
                                        style: { fontSize: '12px' }
                                    }}
                                />
                                <YAxis 
                                    domain={['dataMin - 5', 'dataMax + 5']}
                                    tick={{ fontSize: 10 }}
                                    tickLine={{ stroke: '#e0e0e0' }}
                                    axisLine={{ stroke: '#e0e0e0' }}
                                    label={{ value: 'уд/мин', angle: -90, position: 'insideLeft' }}
                                />
                                <Tooltip content={<CustomTooltip />} />
                                
                                {/* Референсные линии для нормальных значений */}
                                <ReferenceLine y={60} stroke="#10b981" strokeDasharray="2 2" opacity={0.5} />
                                <ReferenceLine y={100} stroke="#10b981" strokeDasharray="2 2" opacity={0.5} />
                                
                                <Line
                                    type="monotone"
                                    dataKey="bpm"
                                    stroke="#ef4444"
                                    strokeWidth={1.5}
                                    dot={{ 
                                        fill: '#ef4444', 
                                        strokeWidth: 1, 
                                        r: 1.5 
                                    }}
                                    activeDot={{ r: 3, stroke: '#ef4444', strokeWidth: 2 }}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>
                    
                    <div className="mt-4 space-y-2">
                        <div className="text-xs text-muted-foreground text-center">
                            Зеленые линии показывают нормальный диапазон пульса в покое (60-100 уд/мин)
                        </div>
                        <div className="text-xs text-muted-foreground text-center">
                            Детальный просмотр каждого R-R интервала с точностью до секунд
                        </div>
                        <div className="text-xs text-blue-600 text-center">
                            Период: {Math.round((dateRange.to.getTime() - dateRange.from.getTime()) / (1000 * 60))} минут
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
} 