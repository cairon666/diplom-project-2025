import { createFileRoute } from '@tanstack/react-router';
import { RRAnalyticsDashboard } from 'src/features/health';
import { z } from 'zod';

const analyticsSearchSchema = z.object({
    from: z.string().optional(),
    to: z.string().optional(),
});

function AnalyticsPage() {
    return (
        <div className="flex-1 p-6">
            <div className="max-w-7xl mx-auto">
                <RRAnalyticsDashboard />
            </div>
        </div>
    );
}

export const Route = createFileRoute('/panel/')({
    component: AnalyticsPage,
    validateSearch: analyticsSearchSchema,
}); 