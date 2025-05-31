import { useState } from 'react';

export interface ChartVisibilityState {
    heartRate: boolean;
    histogram: boolean;
    differentialHistogram: boolean;
    scatterplot: boolean;
    trends: boolean;
}

export function useChartVisibility(initialState: Partial<ChartVisibilityState> = {}) {
    const [visibleCharts, setVisibleCharts] = useState<ChartVisibilityState>({
        heartRate: false,
        histogram: false,
        differentialHistogram: false,
        scatterplot: false,
        trends: false,
        ...initialState,
    });

    const toggleChart = (chartName: keyof ChartVisibilityState) => {
        setVisibleCharts(prev => ({
            ...prev,
            [chartName]: !prev[chartName]
        }));
    };

    const showAllCharts = () => {
        setVisibleCharts({
            heartRate: true,
            histogram: true,
            differentialHistogram: true,
            scatterplot: true,
            trends: true,
        });
    };

    const hideAllCharts = () => {
        setVisibleCharts({
            heartRate: false,
            histogram: false,
            differentialHistogram: false,
            scatterplot: false,
            trends: false,
        });
    };

    const isAnyChartVisible = Object.values(visibleCharts).some(Boolean);

    return {
        visibleCharts,
        toggleChart,
        showAllCharts,
        hideAllCharts,
        isAnyChartVisible,
    };
} 