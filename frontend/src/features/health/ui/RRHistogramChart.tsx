import { useMemo } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from 'ui/card';
import { RRHistogramData } from '../api/rrAnalyticsApi';

interface RRHistogramChartProps {
    data: RRHistogramData | null | undefined;
    className?: string;
}

export function RRHistogramChart({ data, className }: RRHistogramChartProps) {
    // Проверяем данные и предоставляем fallback
    if (!data || !data.bins || !Array.isArray(data.bins) || data.bins.length === 0) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle>Гистограмма R-R интервалов</CardTitle>
                    <CardDescription>
                        Нет данных для построения гистограммы
                    </CardDescription>
                </CardHeader>
                <CardContent className="py-8 text-center">
                    <p className="text-muted-foreground">
                        Недостаточно данных для построения гистограммы. Попробуйте выбрать другой временной диапазон.
                    </p>
                </CardContent>
            </Card>
        );
    }

    const chartData = useMemo(() => {
        if (!data?.bins || data.bins.length === 0) return [];
        
        // Создаем данные только из существующих бинов
        const chartData = data.bins.map(bin => ({
            range: `${bin.range_start || 0}-${bin.range_end || 0}`,
            rangeStart: bin.range_start || 0,
            rangeEnd: bin.range_end || 0,
            count: bin.count || 0,
            frequency: Math.round((bin.frequency || 0) * 100 * 10) / 10,
            isEmpty: false
        }));
        
        // Сортируем по rangeStart для правильного отображения
        return chartData.sort((a, b) => a.rangeStart - b.rangeStart);
    }, [data.bins]);

    const formatTooltip = (value: number, name: string) => {
        if (name === 'count') {
            return [`${value} измерений`, 'Количество'];
        }
        if (name === 'frequency') {
            return [`${value}%`, 'Частота'];
        }
        return [value, name];
    };

    const formatLabel = (label: string) => {
        return `${label} мс`;
    };

    return (
        <Card className={className}>
            <CardHeader>
                <CardTitle>Гистограмма R-R интервалов</CardTitle>
                <CardDescription>
                    Распределение интервалов по значениям. Всего: {(data.total_count || 0).toLocaleString()} измерений
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="h-80">
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart
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
                                dataKey="range"
                                tickFormatter={formatLabel}
                                angle={-45}
                                textAnchor="end"
                                height={60}
                                interval={Math.max(1, Math.ceil(chartData.length / 15))}
                                fontSize={12}
                            />
                            <YAxis />
                            <Tooltip 
                                formatter={formatTooltip}
                                labelFormatter={formatLabel}
                                contentStyle={{
                                    backgroundColor: 'hsl(var(--card))',
                                    border: '1px solid hsl(var(--border))',
                                    borderRadius: '6px',
                                }}
                                content={({ active, payload, label }) => {
                                    if (active && payload && payload.length > 0) {
                                        const data = payload[0].payload;
                                        return (
                                            <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-lg">
                                                <p className="font-medium">{label} мс</p>
                                                <p className="text-blue-600">
                                                    Количество: {data.count || 0} измерений
                                                </p>
                                                <p className="text-gray-600">
                                                    Частота: {data.frequency || 0}%
                                                </p>
                                            </div>
                                        );
                                    }
                                    return null;
                                }}
                            />
                            <Bar 
                                dataKey="count" 
                                fill="hsl(var(--primary))"
                                radius={[2, 2, 0, 0]}
                                opacity={0.8}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
                
                {/* Статистика по гистограмме */}
                <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div className="text-center">
                        <div className="font-semibold">{(data.statistics?.mean || 0).toFixed(1)} мс</div>
                        <div className="text-muted-foreground">Среднее</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold">{(data.statistics?.std_dev || 0).toFixed(1)} мс</div>
                        <div className="text-muted-foreground">Ст. отклонение</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold">{data.bin_width || 0} мс</div>
                        <div className="text-muted-foreground">Ширина бина</div>
                    </div>
                    <div className="text-center">
                        <div className="font-semibold">{chartData.length}</div>
                        <div className="text-muted-foreground">Количество бинов</div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
} 