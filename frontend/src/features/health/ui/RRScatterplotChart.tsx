import React from 'react';
import { ScatterChart, Scatter, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { LoadingSpinner } from 'ui/loading-spinner';
import { useGetScatterplotQuery } from '../api';

interface RRScatterplotChartProps {
    from: string;
    to: string;
}

export function RRScatterplotChart({ from, to }: RRScatterplotChartProps) {
    const { data, isLoading, error } = useGetScatterplotQuery({
        from,
        to,
    });

    if (isLoading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>Диаграмма Пуанкаре</CardTitle>
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
                    <CardTitle>Диаграмма Пуанкаре</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-red-600 p-4">
                        Ошибка загрузки данных скаттерограммы
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (!data || data.points.length === 0) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>Диаграмма Пуанкаре</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-gray-600 p-4">
                        Нет данных для отображения
                    </div>
                </CardContent>
            </Card>
        );
    }

    // Ограничиваем количество точек для производительности
    const maxPoints = 2000;
    const chartData = data.points.length > maxPoints 
        ? data.points.slice(0, maxPoints)
        : data.points;

    const formatTooltip = (value: any, name: string, props: any) => {
        if (name === 'rr_n1') {
            return [`${props.payload.rr_n1} мс`, 'RR(n+1)'];
        }
        return [value, name];
    };

    const formatLabel = (label: string, payload: any[]) => {
        if (payload && payload.length > 0) {
            return `RR(n): ${payload[0].payload.rr_n} мс`;
        }
        return label;
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Диаграмма Пуанкаре RR интервалов</CardTitle>
                <div className="text-sm text-gray-600">
                    <p>Скаттерограмма RR(n) vs RR(n+1) - анализ вариабельности сердечного ритма</p>
                    <div className="grid grid-cols-2 gap-4 mt-2">
                        <div>
                            <span className="font-medium">SD1: </span>
                            {data.statistics.sd1.toFixed(2)} мс
                        </div>
                        <div>
                            <span className="font-medium">SD2: </span>
                            {data.statistics.sd2.toFixed(2)} мс
                        </div>
                        <div>
                            <span className="font-medium">SD1/SD2: </span>
                            {data.statistics.sd1_sd2_ratio.toFixed(3)}
                        </div>
                        <div>
                            <span className="font-medium">CSI: </span>
                            {data.statistics.csi.toFixed(3)}
                        </div>
                        <div>
                            <span className="font-medium">CVI: </span>
                            {data.statistics.cvi.toFixed(3)}
                        </div>
                        <div>
                            <span className="font-medium">Всего точек: </span>
                            {data.total_count} {data.points.length < data.total_count && `(показано ${data.points.length})`}
                        </div>
                    </div>
                </div>
            </CardHeader>
            <CardContent>
                <div className="h-96 w-full">
                    <ResponsiveContainer width="100%" height="100%">
                        <ScatterChart
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
                                type="number"
                                dataKey="rr_n"
                                name="RR(n)"
                                unit="мс"
                                label={{ value: 'RR(n), мс', position: 'bottom' }}
                            />
                            <YAxis 
                                type="number"
                                dataKey="rr_n1"
                                name="RR(n+1)"
                                unit="мс"
                                label={{ value: 'RR(n+1), мс', angle: -90, position: 'insideLeft' }}
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
                            <Scatter name="RR интервалы" fill="#3b82f6" opacity={0.6}>
                                {chartData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill="#3b82f6" />
                                ))}
                            </Scatter>
                        </ScatterChart>
                    </ResponsiveContainer>
                </div>
                
                {/* Дополнительная информация об эллипсе */}
                {data.ellipse && (
                    <div className="mt-4 p-3 bg-gray-50 rounded-lg">
                        <h4 className="font-medium text-sm mb-2">Параметры эллипса:</h4>
                        <div className="grid grid-cols-3 gap-2 text-xs">
                            <div>
                                <span className="font-medium">Центр X: </span>
                                {data.ellipse.center_x.toFixed(2)} мс
                            </div>
                            <div>
                                <span className="font-medium">Центр Y: </span>
                                {data.ellipse.center_y.toFixed(2)} мс
                            </div>
                            <div>
                                <span className="font-medium">Площадь: </span>
                                {data.ellipse.area.toFixed(2)} мс²
                            </div>
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
} 