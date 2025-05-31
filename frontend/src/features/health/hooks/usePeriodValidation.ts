import { differenceInMinutes } from 'date-fns';
import { DateTimeRange } from '../ui/DateTimeRangeSelector';

export interface ValidationRules {
    minMinutes: number;
    maxMinutes: number;
}

export interface PeriodValidationConfig {
    basicAnalysis: ValidationRules;
    detailedHeartRate: ValidationRules;
    trends: ValidationRules;
    histogram: ValidationRules;
    scatterplot: ValidationRules;
}

const DEFAULT_VALIDATION_CONFIG: PeriodValidationConfig = {
    basicAnalysis: { minMinutes: 0.5, maxMinutes: 1440 }, // 30 сек - 24 часа
    detailedHeartRate: { minMinutes: 0.5, maxMinutes: 15 }, // 30 сек - 15 мин
    trends: { minMinutes: 10, maxMinutes: 1440 }, // 10 мин - 24 часа
    histogram: { minMinutes: 1, maxMinutes: 1440 }, // 1 мин - 24 часа
    scatterplot: { minMinutes: 2, maxMinutes: 1440 }, // 2 мин - 24 часа
};

export function usePeriodValidation(
    dateRange: DateTimeRange,
    config: Partial<PeriodValidationConfig> = {}
) {
    const validationConfig = { ...DEFAULT_VALIDATION_CONFIG, ...config };

    const getDurationInMinutes = (): number => {
        return differenceInMinutes(dateRange.to, dateRange.from);
    };

    const isValidForType = (type: keyof PeriodValidationConfig): boolean => {
        const duration = getDurationInMinutes();
        const rules = validationConfig[type];
        return duration >= rules.minMinutes && duration <= rules.maxMinutes;
    };

    const getInvalidMessage = (type: keyof PeriodValidationConfig): string => {
        const duration = getDurationInMinutes();
        const rules = validationConfig[type];
        const durationText = duration < 1 
            ? `${Math.round(duration * 60)} сек` 
            : `${Math.round(duration)} мин`;
        
        if (duration < rules.minMinutes) {
            const minText = rules.minMinutes < 1 
                ? `${Math.round(rules.minMinutes * 60)} секунд`
                : `${rules.minMinutes} минут`;
            return `Период слишком короткий (${durationText}). Минимум: ${minText}.`;
        }
        
        if (duration > rules.maxMinutes) {
            const maxText = rules.maxMinutes >= 60
                ? `${Math.round(rules.maxMinutes / 60)} часов`
                : `${rules.maxMinutes} минут`;
            return `Период слишком длинный (${durationText}). Максимум: ${maxText}.`;
        }
        
        return 'Недопустимый период времени.';
    };

    const shouldShowWarning = (type: 'heartRate' | 'trends'): boolean => {
        const duration = getDurationInMinutes();
        if (type === 'heartRate') {
            return duration > 15; // Показываем предупреждение если >= 15 минут
        }
        if (type === 'trends') {
            return duration < 10; // Показываем предупреждение если <= 10 минут
        }
        return false;
    };

    return {
        getDurationInMinutes,
        isValidForBasicAnalysis: () => isValidForType('basicAnalysis'),
        isValidForDetailedHeartRate: () => isValidForType('detailedHeartRate'),
        isValidForTrends: () => isValidForType('trends'),
        isValidForHistogram: () => isValidForType('histogram'),
        isValidForScatterplot: () => isValidForType('scatterplot'),
        getInvalidMessage,
        shouldShowHeartRateWarning: () => shouldShowWarning('heartRate'),
        shouldShowTrendsWarning: () => shouldShowWarning('trends'),
    };
} 