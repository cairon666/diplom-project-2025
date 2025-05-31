package apperrors

import "net/http"

// Дополнительные конструкторы ошибок

func ValidationFailed() AppError {
	return New(CodeValidationFailed, "validation failed", http.StatusBadRequest)
}

func ValidationFailedf(format string, args ...interface{}) AppError {
	return Newf(CodeValidationFailed, http.StatusBadRequest, format, args...)
}

func Unauthorized() AppError {
	return New(CodeUnauthorized, "unauthorized", http.StatusUnauthorized)
}

func Unauthorizedf(format string, args ...interface{}) AppError {
	return Newf(CodeUnauthorized, http.StatusUnauthorized, format, args...)
}

func BadRequest() AppError {
	return New(CodeBadRequest, "bad request", http.StatusBadRequest)
}

func BadRequestf(format string, args ...interface{}) AppError {
	return Newf(CodeBadRequest, http.StatusBadRequest, format, args...)
}

func Conflict() AppError {
	return New(CodeConflict, "conflict", http.StatusConflict)
}

func Conflictf(format string, args ...interface{}) AppError {
	return Newf(CodeConflict, http.StatusConflict, format, args...)
}

func TooManyRequests() AppError {
	return New(CodeTooManyRequests, "too many requests", http.StatusTooManyRequests)
}

func ServiceUnavailable() AppError {
	return New(CodeServiceUnavailable, "service unavailable", http.StatusServiceUnavailable)
}

// User specific errors

func UserNotFound() AppError {
	return New(CodeUserNotFound, "user not found", http.StatusNotFound)
}

func UserNotFoundf(format string, args ...interface{}) AppError {
	return Newf(CodeUserNotFound, http.StatusNotFound, format, args...)
}

// RR Intervals specific errors

func DeviceNotFound() AppError {
	return New(CodeDeviceNotFound, "device not found", http.StatusNotFound)
}

func DeviceNotFoundf(format string, args ...interface{}) AppError {
	return Newf(CodeDeviceNotFound, http.StatusNotFound, format, args...)
}

func DeviceAccessDenied() AppError {
	return New(CodeDeviceAccessDenied, "access to device denied", http.StatusForbidden)
}

func DeviceAccessDeniedf(format string, args ...interface{}) AppError {
	return Newf(CodeDeviceAccessDenied, http.StatusForbidden, format, args...)
}

func InvalidRRInterval() AppError {
	return New(CodeInvalidRRInterval, "invalid RR interval value", http.StatusBadRequest)
}

func InvalidRRIntervalf(format string, args ...interface{}) AppError {
	return Newf(CodeInvalidRRInterval, http.StatusBadRequest, format, args...)
}

func InsufficientData() AppError {
	return New(CodeInsufficientData, "insufficient data for analysis", http.StatusBadRequest)
}

func InsufficientDataf(format string, args ...interface{}) AppError {
	return Newf(CodeInsufficientData, http.StatusBadRequest, format, args...)
}

func InvalidTimeRange() AppError {
	return New(CodeInvalidTimeRange, "invalid time range", http.StatusBadRequest)
}

func InvalidTimeRangef(format string, args ...interface{}) AppError {
	return Newf(CodeInvalidTimeRange, http.StatusBadRequest, format, args...)
}

func AnalysisNotPossible() AppError {
	return New(CodeAnalysisNotPossible, "analysis not possible with current data", http.StatusBadRequest)
}

func AnalysisNotPossiblef(format string, args ...interface{}) AppError {
	return Newf(CodeAnalysisNotPossible, http.StatusBadRequest, format, args...)
}

func BatchTooLarge() AppError {
	return New(CodeBatchTooLarge, "batch size exceeds maximum allowed", http.StatusBadRequest)
}

func BatchTooLargef(format string, args ...interface{}) AppError {
	return Newf(CodeBatchTooLarge, http.StatusBadRequest, format, args...)
}

func BatchEmpty() AppError {
	return New(CodeBatchEmpty, "batch cannot be empty", http.StatusBadRequest)
}

func BatchEmptyf(format string, args ...interface{}) AppError {
	return Newf(CodeBatchEmpty, http.StatusBadRequest, format, args...)
}

func InvalidDataFormat() AppError {
	return New(CodeInvalidDataFormat, "invalid data format", http.StatusBadRequest)
}

func InvalidDataFormatf(format string, args ...interface{}) AppError {
	return Newf(CodeInvalidDataFormat, http.StatusBadRequest, format, args...)
}

func DataProcessingError() AppError {
	return New(CodeDataProcessingError, "error processing data", http.StatusInternalServerError)
}

func DataProcessingErrorf(format string, args ...interface{}) AppError {
	return Newf(CodeDataProcessingError, http.StatusInternalServerError, format, args...)
}

func TimeRangeTooSmall() AppError {
	return New(CodeTimeRangeTooSmall, "time range too small for analysis", http.StatusBadRequest)
}

func TimeRangeTooSmallf(format string, args ...interface{}) AppError {
	return Newf(CodeTimeRangeTooSmall, http.StatusBadRequest, format, args...)
}

func TimeRangeTooLarge() AppError {
	return New(CodeTimeRangeTooLarge, "time range too large", http.StatusBadRequest)
}

func TimeRangeTooLargef(format string, args ...interface{}) AppError {
	return Newf(CodeTimeRangeTooLarge, http.StatusBadRequest, format, args...)
}

func ParameterOutOfRange() AppError {
	return New(CodeParameterOutOfRange, "parameter value out of allowed range", http.StatusBadRequest)
}

func ParameterOutOfRangef(format string, args ...interface{}) AppError {
	return Newf(CodeParameterOutOfRange, http.StatusBadRequest, format, args...)
}

func NoValidData() AppError {
	return New(CodeNoValidData, "no valid data found", http.StatusNotFound)
}

func NoValidDataf(format string, args ...interface{}) AppError {
	return Newf(CodeNoValidData, http.StatusNotFound, format, args...)
}

// JWT specific errors

func TokenCreationFailed() AppError {
	return New(CodeTokenCreationFailed, "failed to create token", http.StatusInternalServerError)
}

func TokenCreationFailedf(format string, args ...interface{}) AppError {
	return Newf(CodeTokenCreationFailed, http.StatusInternalServerError, format, args...)
}

// Health data specific errors

func HealthDataQueryFailed() AppError {
	return New(CodeHealthDataQueryFailed, "failed to query health data", http.StatusInternalServerError)
}

func HealthDataQueryFailedf(format string, args ...interface{}) AppError {
	return Newf(CodeHealthDataQueryFailed, http.StatusInternalServerError, format, args...)
}

func HealthDataReadFailed() AppError {
	return New(CodeHealthDataReadFailed, "failed to read health data", http.StatusInternalServerError)
}

func HealthDataReadFailedf(format string, args ...interface{}) AppError {
	return Newf(CodeHealthDataReadFailed, http.StatusInternalServerError, format, args...)
}