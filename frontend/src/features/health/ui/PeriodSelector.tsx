import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { DateTimeRangeSelector, DateTimeRange } from './DateTimeRangeSelector';

interface PeriodSelectorProps {
    period1: DateTimeRange;
    period2: DateTimeRange;
    onPeriod1Change: (period: DateTimeRange) => void;
    onPeriod2Change: (period: DateTimeRange) => void;
}

export function PeriodSelector({
    period1,
    period2,
    onPeriod1Change,
    onPeriod2Change,
}: PeriodSelectorProps) {
    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Первый период */}
            <Card className="border-blue-200 gap-2">
                <CardHeader>
                    <CardTitle className="text-lg">Период 1</CardTitle>
                </CardHeader>
                <CardContent>
                    <DateTimeRangeSelector
                        value={period1}
                        onChange={onPeriod1Change}
                        className="w-full"
                    />
                </CardContent>
            </Card>

            {/* Второй период */}
            <Card className="border-red-200 gap-2">
                <CardHeader>
                    <CardTitle className="text-lg">Период 2</CardTitle>
                </CardHeader>
                <CardContent>
                    <DateTimeRangeSelector
                        value={period2}
                        onChange={onPeriod2Change}
                        className="w-full"
                    />
                </CardContent>
            </Card>
        </div>
    );
} 