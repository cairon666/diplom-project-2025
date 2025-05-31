import { createFileRoute } from '@tanstack/react-router';
import { ComparativeAnalysisDashboard } from 'src/features/health';
import { z } from 'zod';

const comparativeAnalysisSearchSchema = z.object({
    period1_from: z.string().optional(),
    period1_to: z.string().optional(),
    period2_from: z.string().optional(),
    period2_to: z.string().optional(),
});

function ComparativeAnalysisPage() {
    return (
        <div className="flex-1 p-6">
            <div className="max-w-7xl mx-auto">
                <ComparativeAnalysisDashboard />
            </div>
        </div>
    );
}

export const Route = createFileRoute('/panel/comparative-analysis')({
    component: ComparativeAnalysisPage,
    validateSearch: comparativeAnalysisSearchSchema,
}); 