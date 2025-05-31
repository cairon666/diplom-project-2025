import { useState } from 'react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { Button } from 'ui/button';
import { Input } from 'ui/input';
import { Label } from 'ui/label';
import { LuCalendar, LuClock } from 'react-icons/lu';
import { AlertTriangle } from 'lucide-react';

export interface DateTimeRange {
    from: Date;
    to: Date;
}

interface DateTimeRangeSelectorProps {
    value: DateTimeRange;
    onChange: (range: DateTimeRange) => void;
    className?: string;
    hasUnsavedChanges?: boolean;
    onApplyChanges?: () => void;
    manualMode?: boolean;
}


export function DateTimeRangeSelector({ 
    value, 
    onChange, 
    className, 
    hasUnsavedChanges = false,
    onApplyChanges,
    manualMode = false 
}: DateTimeRangeSelectorProps) {
    const [isExpanded, setIsExpanded] = useState(false);

    const handleFromDateChange = (dateStr: string) => {
        const newFrom = new Date(dateStr);
        if (!isNaN(newFrom.getTime())) {
            onChange({
                from: newFrom,
                to: value.to,
            });
        }
    };

    const handleToDateChange = (dateStr: string) => {
        const newTo = new Date(dateStr);
        if (!isNaN(newTo.getTime())) {
            onChange({
                from: value.from,
                to: newTo,
            });
        }
    };

    const formatDateTimeForInput = (date: Date): string => {
        // Формат: YYYY-MM-DDTHH:mm:ss для точности до секунд
        return format(date, "yyyy-MM-dd'T'HH:mm:ss");
    };

    const formatRangeDisplay = (): string => {
        const fromStr = format(value.from, 'dd MMM yyyy, HH:mm:ss', { locale: ru });
        const toStr = format(value.to, 'dd MMM yyyy, HH:mm:ss', { locale: ru });
        return `${fromStr} - ${toStr}`;
    };

    return (
        <Card className={className}>
            <CardHeader>
                <div className="flex items-center justify-between">
                    <CardTitle className="flex items-center gap-2 text-sm">
                        <LuCalendar className="w-4 h-4" />
                        Период времени 
                        {manualMode && (
                            <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                                Ручной режим
                            </span>
                        )}
                        {hasUnsavedChanges && manualMode && (
                            <span className="flex items-center gap-1 text-xs bg-orange-100 text-orange-800 px-2 py-1 rounded">
                                <AlertTriangle className="w-3 h-3" />
                                Не применено
                            </span>
                        )}
                    </CardTitle>
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setIsExpanded(!isExpanded)}
                    >
                        {isExpanded ? 'Свернуть' : 'Развернуть'}
                    </Button>
                </div>
                
                {!isExpanded && (
                    <div className="text-xs text-muted-foreground">
                        {formatRangeDisplay()}
                        {hasUnsavedChanges && manualMode && (
                            <div className="flex items-center gap-2 mt-2">
                                <span className="text-orange-600">Изменения не применены</span>
                                {onApplyChanges && (
                                    <Button size="sm" variant="outline" onClick={onApplyChanges}>
                                        Применить
                                    </Button>
                                )}
                            </div>
                        )}
                    </div>
                )}
            </CardHeader>

            {isExpanded && (
                <CardContent className="space-y-4">
                    {/* Ручной выбор даты и времени */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <Label htmlFor="from-datetime" className="text-xs flex items-center gap-1">
                                <LuClock className="w-3 h-3" />
                                С (дата, время, секунды)
                            </Label>
                            <Input
                                id="from-datetime"
                                type="datetime-local"
                                step="1"
                                value={formatDateTimeForInput(value.from)}
                                onChange={(e) => handleFromDateChange(e.target.value)}
                                className="text-xs"
                            />
                        </div>

                        <div>
                            <Label htmlFor="to-datetime" className="text-xs flex items-center gap-1">
                                <LuClock className="w-3 h-3" />
                                По (дата, время, секунды)
                            </Label>
                            <Input
                                id="to-datetime"
                                type="datetime-local"
                                step="1"
                                value={formatDateTimeForInput(value.to)}
                                onChange={(e) => handleToDateChange(e.target.value)}
                                className="text-xs"
                            />
                        </div>
                    </div>

                    {/* Текущий выбранный диапазон */}
                    <div className="p-2 bg-muted rounded text-xs">
                        <span className="text-muted-foreground">
                            Длительность: {Math.round((value.to.getTime() - value.from.getTime()) / (1000 * 60))} минут
                        </span>
                    </div>
                </CardContent>
            )}
        </Card>
    );
} 