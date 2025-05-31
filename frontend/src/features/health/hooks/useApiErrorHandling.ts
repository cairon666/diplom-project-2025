import { 
    ApiErrorCode, 
    isBaseError, 
} from '../../../store/apiErrors';

export function useApiErrorHandling() {
    const getErrorMessage = (error: unknown): string => {
        if (!error) return 'Неизвестная ошибка';

        // Проверяем, является ли это базовой ошибкой API
        if (isBaseError(error)) {
            const { data } = error;
            switch (data.error as ApiErrorCode) {
                case ApiErrorCode.CodeNotFound:
                    return 'Данные не найдены для выбранного периода';
                case ApiErrorCode.CodeInsufficientData:
                    return 'Недостаточно данных для анализа в выбранном периоде';
                case ApiErrorCode.CodeInvalidTimeRange:
                    return 'Некорректный временной диапазон';
                case ApiErrorCode.CodeTimeRangeTooSmall:
                    return 'Слишком маленький временной диапазон для анализа';
                case ApiErrorCode.CodeTimeRangeTooLarge:
                    return 'Слишком большой временной диапазон. Попробуйте уменьшить период';
                case ApiErrorCode.CodeDeviceNotFound:
                    return 'Устройство не найдено';
                case ApiErrorCode.CodeDeviceAccessDenied:
                    return 'Доступ к устройству запрещен';
                case ApiErrorCode.CodeInvalidRRInterval:
                    return 'Обнаружены некорректные R-R интервалы';
                case ApiErrorCode.CodeAnalysisNotPossible:
                    return 'Анализ невозможен с текущими данными';
                case ApiErrorCode.CodeBatchTooLarge:
                    return 'Слишком большой объем данных для обработки';
                case ApiErrorCode.CodeBatchEmpty:
                    return 'Отсутствуют данные для анализа';
                case ApiErrorCode.CodeInvalidDataFormat:
                    return 'Некорректный формат данных';
                case ApiErrorCode.CodeDataProcessingError:
                    return 'Ошибка обработки данных на сервере';
                case ApiErrorCode.CodeParameterOutOfRange:
                    return 'Параметры анализа выходят за допустимые пределы';
                case ApiErrorCode.CodeNoValidData:
                    return 'Нет валидных данных в выбранном периоде';
                case ApiErrorCode.CodeValidationFailed:
                    return 'Ошибка валидации запроса';
                case ApiErrorCode.CodeUnauthorized:
                    return 'Требуется авторизация';
                case ApiErrorCode.CodeForbidden:
                    return 'Доступ запрещен';
                case ApiErrorCode.CodeBadRequest:
                    return 'Некорректный запрос';
                case ApiErrorCode.CodeConflict:
                    return 'Конфликт данных';
                case ApiErrorCode.CodeTooManyRequests:
                    return 'Слишком много запросов. Попробуйте позже';
                case ApiErrorCode.CodeServiceUnavailable:
                    return 'Сервис временно недоступен';
                case ApiErrorCode.CodeInternalError:
                    return 'Внутренняя ошибка сервера';
                case ApiErrorCode.CodeHealthDataQueryFailed:
                    return 'Ошибка запроса данных о здоровье';
                case ApiErrorCode.CodeHealthDataReadFailed:
                    return 'Ошибка чтения данных о здоровье';
                default:
                    return data.message || 'Ошибка API';
            }
        }

        // Обработка обычных ошибок
        if (error && typeof error === 'object') {
            const err = error as { data?: { message?: string }; message?: string; status?: number };
            
            // Обработка HTTP статусов
            if (err.status) {
                switch (err.status) {
                    case 400:
                        return 'Некорректный запрос';
                    case 401:
                        return 'Требуется авторизация';
                    case 403:
                        return 'Доступ запрещен';
                    case 404:
                        return 'Данные не найдены';
                    case 409:
                        return 'Конфликт данных';
                    case 429:
                        return 'Слишком много запросов';
                    case 500:
                        return 'Внутренняя ошибка сервера';
                    case 503:
                        return 'Сервис временно недоступен';
                    default:
                        return `Ошибка сервера (${err.status})`;
                }
            }
            
            if (err.data?.message) return err.data.message;
            if (err.message) return err.message;
        }
        
        return 'Неизвестная ошибка сервера';
    };

    const combineErrors = (...errors: unknown[]): unknown | null => {
        return errors.find(error => !!error) || null;
    };

    return {
        getErrorMessage,
        combineErrors,
    };
} 