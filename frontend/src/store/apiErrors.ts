import { FetchBaseQueryError } from '@reduxjs/toolkit/query/react';

export enum ApiErrorCode {
    CodeNotFound = 'NOT_FOUND',
    CodeAlreadyExists = 'ALREADY_EXISTS',
    CodeEmailAlreadyExists = 'EMAIL_ALREADY_EXISTS',
    CodeLoginNotRegistered = 'LOGIN_NOT_REGISTERED',
    CodeWrongPassword = 'WRONG_PASSWORD',
    CodeInternalError = 'INTERNAL_ERROR',
    CodeInvalidToken = 'INVALID_TOKEN',
    CodeForbidden = 'FORBIDDEN',
    CodeInvalidParams = 'INVALID_PARAMS',
    CodeProviderAlreadyConnected = 'PROVIDER_ALREADY_CONNECTED',
    CodeProviderAccountAlreadyLinked = 'PROVIDER_ACCOUNT_ALREADY_LINKED',
    CodeNeedEndRegistration = 'NEED_END_REGISTRATION',
    CodeTelegramAlreadyRegistered = 'TELEGRAM_ALREADY_REGISTERED',
    CodeTempIdNotFound = 'TEMP_ID_NOT_FOUND',
    CodeInvalidTelegramHash = 'INVALID_TELEGRAM_HASH',
    CodeTelegramIsNotLinked = 'TELEGRAM_IS_NOT_LINKED',
    
    // User specific errors
    CodeUserNotFound = 'USER_NOT_FOUND',
    
    // JWT specific errors
    CodeTokenCreationFailed = 'TOKEN_CREATION_FAILED',
    
    // Health data specific errors
    CodeHealthDataQueryFailed = 'HEALTH_DATA_QUERY_FAILED',
    CodeHealthDataReadFailed = 'HEALTH_DATA_READ_FAILED',
    
    // Additional error codes
    CodeValidationFailed = 'VALIDATION_FAILED',
    CodeUnauthorized = 'UNAUTHORIZED',
    CodeBadRequest = 'BAD_REQUEST',
    CodeConflict = 'CONFLICT',
    CodeTooManyRequests = 'TOO_MANY_REQUESTS',
    CodeServiceUnavailable = 'SERVICE_UNAVAILABLE',
    
    // RR Intervals specific errors
    CodeDeviceNotFound = 'DEVICE_NOT_FOUND',
    CodeDeviceAccessDenied = 'DEVICE_ACCESS_DENIED',
    CodeInvalidRRInterval = 'INVALID_RR_INTERVAL',
    CodeInsufficientData = 'INSUFFICIENT_DATA',
    CodeInvalidTimeRange = 'INVALID_TIME_RANGE',
    CodeAnalysisNotPossible = 'ANALYSIS_NOT_POSSIBLE',
    CodeBatchTooLarge = 'BATCH_TOO_LARGE',
    CodeBatchEmpty = 'BATCH_EMPTY',
    CodeInvalidDataFormat = 'INVALID_DATA_FORMAT',
    CodeDataProcessingError = 'DATA_PROCESSING_ERROR',
    CodeTimeRangeTooSmall = 'TIME_RANGE_TOO_SMALL',
    CodeTimeRangeTooLarge = 'TIME_RANGE_TOO_LARGE',
    CodeParameterOutOfRange = 'PARAMETER_OUT_OF_RANGE',
    CodeNoValidData = 'NO_VALID_DATA',
}

export interface AppBaseError {
    message: string;
    error: ApiErrorCode;
}

export interface AppErrorWithFields<Fields> extends AppBaseError {
    fields: Fields;
}

export function isFetchBaseQueryError(
    error: unknown,
): error is FetchBaseQueryError {
    return typeof error === 'object' && error !== null && 'status' in error;
}

export function isApiError(error: unknown): error is AppBaseError {
    return (
        typeof error === 'object' &&
        error != null &&
        'message' in error &&
        typeof (error as any).message === 'string' &&
        'error' in (error as any) &&
        typeof (error as any).error === 'string'
    );
}

export function isBaseError(
    error: unknown,
): error is FetchBaseQueryError & { data: AppBaseError } {
    return isFetchBaseQueryError(error) && isApiError(error.data);
}

export function isApiErrorWithFields<Fields>(
    error: unknown,
): error is FetchBaseQueryError & { data: AppErrorWithFields<Fields> } {
    return (
        isBaseError(error) &&
        'fields' in error.data &&
        typeof error.data.fields === 'object'
    );
}
