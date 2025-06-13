package apperrors

type ErrorCode string

var (
	CodeNotFound                     ErrorCode = "NOT_FOUND"
	CodeAlreadyExists                ErrorCode = "ALREADY_EXISTS"
	CodeEmailAlreadyExists           ErrorCode = "EMAIL_ALREADY_EXISTS"
	CodeLoginNotRegistered           ErrorCode = "LOGIN_NOT_REGISTERED"
	CodeWrongPassword                ErrorCode = "WRONG_PASSWORD"
	CodeInternalError                ErrorCode = "INTERNAL_ERROR"
	CodeInvalidToken                 ErrorCode = "INVALID_TOKEN"
	CodeForbidden                    ErrorCode = "FORBIDDEN"
	CodeInvalidParams                ErrorCode = "INVALID_PARAMS"
	CodeProviderAlreadyConnected     ErrorCode = "PROVIDER_ALREADY_CONNECTED"
	CodeProviderAccountAlreadyLinked ErrorCode = "PROVIDER_ACCOUNT_ALREADY_LINKED"
	CodeNeedEndRegistration          ErrorCode = "NEED_END_REGISTRATION"
	CodeTempIdNotFound               ErrorCode = "TEMP_ID_NOT_FOUND"
	CodeInvalidTelegramHash          ErrorCode = "INVALID_TELEGRAM_HASH"
	CodeTelegramIsNotLinked          ErrorCode = "TELEGRAM_IS_NOT_LINKED"

	// User specific errors.
	CodeUserNotFound ErrorCode = "USER_NOT_FOUND"

	// JWT specific errors.
	CodeTokenCreationFailed ErrorCode = "TOKEN_CREATION_FAILED"

	// Health data specific errors.
	CodeHealthDataQueryFailed ErrorCode = "HEALTH_DATA_QUERY_FAILED"
	CodeHealthDataReadFailed  ErrorCode = "HEALTH_DATA_READ_FAILED"

	// Дополнительные коды ошибок.
	CodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeTooManyRequests    ErrorCode = "TOO_MANY_REQUESTS"
	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// RR Intervals specific errors.
	CodeDeviceNotFound      ErrorCode = "DEVICE_NOT_FOUND"
	CodeDeviceAccessDenied  ErrorCode = "DEVICE_ACCESS_DENIED"
	CodeInvalidRRInterval   ErrorCode = "INVALID_RR_INTERVAL"
	CodeInsufficientData    ErrorCode = "INSUFFICIENT_DATA"
	CodeInvalidTimeRange    ErrorCode = "INVALID_TIME_RANGE"
	CodeAnalysisNotPossible ErrorCode = "ANALYSIS_NOT_POSSIBLE"
	CodeBatchTooLarge       ErrorCode = "BATCH_TOO_LARGE"
	CodeBatchEmpty          ErrorCode = "BATCH_EMPTY"
	CodeInvalidDataFormat   ErrorCode = "INVALID_DATA_FORMAT"
	CodeDataProcessingError ErrorCode = "DATA_PROCESSING_ERROR"
	CodeTimeRangeTooSmall   ErrorCode = "TIME_RANGE_TOO_SMALL"
	CodeTimeRangeTooLarge   ErrorCode = "TIME_RANGE_TOO_LARGE"
	CodeParameterOutOfRange ErrorCode = "PARAMETER_OUT_OF_RANGE"
	CodeNoValidData         ErrorCode = "NO_VALID_DATA"
)
