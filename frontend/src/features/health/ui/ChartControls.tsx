import React from 'react';
import { Button } from 'ui/button';
import { 
    Eye, 
    EyeOff, 
    ChevronDown, 
    ChevronRight, 
    LucideIcon 
} from 'lucide-react';
import { ChartVisibilityState } from '../hooks/useChartVisibility';

export interface ChartVisibilityControlsProps {
    showAllCharts: () => void;
    hideAllCharts: () => void;
}

export const ChartVisibilityControls: React.FC<ChartVisibilityControlsProps> = ({
    showAllCharts,
    hideAllCharts,
}) => (
    <div className="flex items-center gap-2 px-3 py-1 bg-muted/50 rounded-lg">
        <span className="text-sm font-medium">Графики:</span>
        <Button
            variant="outline"
            size="sm"
            onClick={showAllCharts}
            className="h-7 px-2 text-xs"
        >
            <Eye className="h-3 w-3 mr-1" />
            Показать все
        </Button>
        <Button
            variant="outline"
            size="sm"
            onClick={hideAllCharts}
            className="h-7 px-2 text-xs"
        >
            <EyeOff className="h-3 w-3 mr-1" />
            Скрыть все
        </Button>
    </div>
);

export interface ChartToggleButtonProps {
    chartName: keyof ChartVisibilityState;
    title: string;
    icon: LucideIcon;
    isVisible: boolean;
    isValid: boolean;
    invalidMessage?: string;
    onToggle: (chartName: keyof ChartVisibilityState) => void;
}

export const ChartToggleButton: React.FC<ChartToggleButtonProps> = ({
    chartName,
    title,
    icon: Icon,
    isVisible,
    isValid,
    invalidMessage,
    onToggle,
}) => (
    <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
            <Icon className="h-5 w-5" />
            <h3 className="text-xl font-semibold">{title}</h3>
            {!isValid && (
                <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                    Недоступен
                </span>
            )}
        </div>
        <Button
            variant={isVisible ? "default" : "outline"}
            size="sm"
            onClick={() => onToggle(chartName)}
            disabled={!isValid}
            className="flex items-center gap-2"
            title={!isValid ? invalidMessage : undefined}
        >
            {isVisible ? (
                <>
                    <ChevronDown className="h-4 w-4" />
                    Скрыть
                </>
            ) : (
                <>
                    <ChevronRight className="h-4 w-4" />
                    Показать
                </>
            )}
        </Button>
    </div>
); 