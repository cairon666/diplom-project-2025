import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';
import { LoadingOverlay } from 'ui/loading-overlay';
import { RRStatistics } from './RRStatistics';
import type { RRStatisticalSummary, HRVMetrics } from '../api/rrAnalyticsApi';

interface PeriodStatisticsProps {
    title: string;
    colorClass: string;
    statistics?: RRStatisticalSummary;
    hrvMetrics?: HRVMetrics;
    isLoading: boolean;
}

export function PeriodStatistics({
    title,
    colorClass,
    statistics,
    hrvMetrics,
    isLoading,
}: PeriodStatisticsProps) {
    return (
        <Card className={colorClass}>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <div className={`w-4 h-4 rounded-full ${colorClass.includes('blue') ? 'bg-blue-500' : 'bg-red-500'}`} />
                    {title}
                </CardTitle>
            </CardHeader>
            <CardContent>
                <LoadingOverlay isLoading={isLoading}>
                    <RRStatistics
                        statistics={statistics}
                        hrvMetrics={hrvMetrics}
                        isLoading={isLoading}
                        compact={true}
                    />
                </LoadingOverlay>
            </CardContent>
        </Card>
    );
} 