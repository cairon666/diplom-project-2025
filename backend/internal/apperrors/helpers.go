package apperrors

import "errors"

// Вспомогательные функции для проверки типов ошибок

// IsNotFound проверяет, является ли ошибка типа "not found".
func IsNotFound(err error) bool {
	return IsErrorWithCode(err, CodeNotFound)
}

// IsValidationFailed проверяет, является ли ошибка типа "validation failed".
func IsValidationFailed(err error) bool {
	return IsErrorWithCode(err, CodeValidationFailed)
}

// IsUnauthorized проверяет, является ли ошибка типа "unauthorized".
func IsUnauthorized(err error) bool {
	return IsErrorWithCode(err, CodeUnauthorized)
}

// IsInternalError проверяет, является ли ошибка внутренней ошибкой сервера.
func IsInternalError(err error) bool {
	return IsErrorWithCode(err, CodeInternalError)
}

// IsInvalidParams проверяет, является ли ошибка типа "invalid params".
func IsInvalidParams(err error) bool {
	return IsErrorWithCode(err, CodeInvalidParams)
}

// IsInvalidToken проверяет, является ли ошибка типа "invalid token".
func IsInvalidToken(err error) bool {
	return IsErrorWithCode(err, CodeInvalidToken)
}

// IsForbidden проверяет, является ли ошибка типа "forbidden".
func IsForbidden(err error) bool {
	return IsErrorWithCode(err, CodeForbidden)
}

// IsConflict проверяет, является ли ошибка типа "conflict".
func IsConflict(err error) bool {
	return IsErrorWithCode(err, CodeConflict)
}

// IsErrorWithCode проверяет, имеет ли ошибка определенный код.
func IsErrorWithCode(err error, code ErrorCode) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.code == code
	}

	return false
}

// GetErrorCode возвращает код ошибки, если это AppError.
func GetErrorCode(err error) (ErrorCode, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.code, true
	}

	return "", false
}

// GetErrorFields возвращает поля ошибки, если это AppError.
func GetErrorFields(err error) (map[string]interface{}, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.fields, true
	}

	return nil, false
}
