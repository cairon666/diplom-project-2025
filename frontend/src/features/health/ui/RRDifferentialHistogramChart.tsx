import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { LoadingSpinner } from 'ui/loading-spinner';
import { useGetDifferentialHistogramQuery } from '../api';

interface RRDifferentialHistogramChartProps {
    from: string;
    to: string;
    binsCount?: number;
}

export function RRDifferentialHistogramChart({ from, to, binsCount }: RRDifferentialHistogramChartProps) {
    const { data, isLoading, error } = useGetDifferentialHistogramQuery({
        from,
        to,
        bins_count: binsCount,
    });

    if (isLoading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>Дифференциальная гистограмма RR</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center justify-center h-96">
                        <LoadingSpinner />
                        <span className="ml-2">Загрузка данных...</span>
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (error) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>Дифференциальная гистограмма RR</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-red-600 p-4">
                        Ошибка загрузки данных дифференциальной гистограммы
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (!data || data.bins.length === 0) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>Дифференциальная гистограмма RR</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-gray-600 p-4">
                        Нет данных для отображения
                    </div>
                </CardContent>
            </Card>
        );
    }

    const chartData = data.bins.map(bin => ({
        range: `${bin.range_start}-${bin.range_end}`,
        rangeStart: bin.range_start,
        rangeEnd: bin.range_end,
        count: bin.count,
        frequency: Math.round(bin.frequency * 100 * 100) / 100, // округляем до 2 знаков
    }));

    const formatTooltip = (value: any, name: string) => {
        if (name === 'frequency') {
            return [`${value}%`, 'Частота'];
        }
        if (name === 'count') {
            return [value, 'Количество'];
        }
        return [value, name];
    };

    const formatLabel = (label: string) => {
        return `${label} мс`;
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Дифференциальная гистограмма RR интервалов</CardTitle>
                <div className="text-sm text-gray-600">
                    <p>Распределение различий между соседними RR интервалами</p>
                    <div className="grid grid-cols-2 gap-4 mt-2">
                        <div>
                            <span className="font-medium">Среднее: </span>
                            {data.statistics.mean.toFixed(2)} мс
                        </div>
                        <div>
                            <span className="font-medium">Станд. откл.: </span>
                            {data.statistics.std_dev.toFixed(2)} мс
                        </div>
                        <div>
                            <span className="font-medium">RMSSD: </span>
                            {data.statistics.rmssd.toFixed(2)} мс
                        </div>
                        <div>
                            <span className="font-medium">Всего точек: </span>
                            {data.total_count}
                        </div>
                    </div>
                </div>
            </CardHeader>
            <CardContent>
                <div className="h-96 w-full">
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart
                            data={chartData}
                            margin={{
                                top: 20,
                                right: 30,
                                left: 20,
                                bottom: 60,
                            }}
                        >
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis 
                                dataKey="range"
                                angle={-45}
                                textAnchor="end"
                                height={60}
                                fontSize={12}
                            />
                            <YAxis 
                                yAxisId="left"
                                orientation="left"
                                tickFormatter={(value) => `${value}`}
                                label={{ value: 'Количество', angle: -90, position: 'insideLeft' }}
                            />
                            <YAxis 
                                yAxisId="right"
                                orientation="right"
                                tickFormatter={(value) => `${value}%`}
                                label={{ value: 'Частота (%)', angle: 90, position: 'insideRight' }}
                            />
                            <Tooltip 
                                formatter={formatTooltip}
                                labelFormatter={formatLabel}
                                contentStyle={{
                                    backgroundColor: 'white',
                                    border: '1px solid #ccc',
                                    borderRadius: '4px',
                                }}
                            />
                            <Bar 
                                yAxisId="left"
                                dataKey="count" 
                                fill="#3b82f6" 
                                name="count"
                                opacity={0.8}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </CardContent>
        </Card>
    );
} 