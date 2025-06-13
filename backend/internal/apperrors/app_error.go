package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error - унифицированная структура ошибки.
type Error struct {
	code     ErrorCode
	message  string
	httpCode int
	fields   map[string]interface{}
}

// Реализация интерфейса AppError.
func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func (e *Error) HTTPCode() int {
	return e.httpCode
}

// Is реализует интерфейс для errors.Is
// Ошибки считаются равными, если у них одинаковый код.
func (e *Error) Is(target error) bool {
	var targetErr *Error
	if errors.As(target, &targetErr) {
		return e.code == targetErr.code
	}

	return false
}

func (e *Error) JSON() map[string]interface{} {
	result := map[string]interface{}{
		"error":   e.code,
		"message": e.message,
	}

	if len(e.fields) > 0 {
		result["fields"] = e.fields
	}

	return result
}

func (e *Error) WithMessage(message string) AppError {
	return &Error{
		code:     e.code,
		message:  message,
		httpCode: e.httpCode,
		fields:   copyFields(e.fields),
	}
}

func (e *Error) WithField(key string, value interface{}) AppError {
	fields := copyFields(e.fields)
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields[key] = value

	return &Error{
		code:     e.code,
		message:  e.message,
		httpCode: e.httpCode,
		fields:   fields,
	}
}

func (e *Error) WithFields(newFields map[string]interface{}) AppError {
	fields := copyFields(e.fields)
	if fields == nil {
		fields = make(map[string]interface{})
	}

	for k, v := range newFields {
		fields[k] = v
	}

	return &Error{
		code:     e.code,
		message:  e.message,
		httpCode: e.httpCode,
		fields:   fields,
	}
}

// ErrorBuilder для удобного создания ошибок.
type errorBuilder struct {
	err *Error
}

func (b *errorBuilder) Build() AppError {
	return &Error{
		code:     b.err.code,
		message:  b.err.message,
		httpCode: b.err.httpCode,
		fields:   copyFields(b.err.fields),
	}
}

func (b *errorBuilder) WithMessage(message string) ErrorBuilder {
	b.err.message = message

	return b
}

func (b *errorBuilder) WithField(key string, value interface{}) ErrorBuilder {
	if b.err.fields == nil {
		b.err.fields = make(map[string]interface{})
	}
	b.err.fields[key] = value

	return b
}

func (b *errorBuilder) WithFields(fields map[string]interface{}) ErrorBuilder {
	if b.err.fields == nil {
		b.err.fields = make(map[string]interface{})
	}
	for k, v := range fields {
		b.err.fields[k] = v
	}

	return b
}

// Конструкторы.
func New(code ErrorCode, message string, httpCode int) AppError {
	return &Error{
		code:     code,
		message:  message,
		httpCode: httpCode,
		fields:   nil,
	}
}

func NewBuilder(code ErrorCode, message string, httpCode int) ErrorBuilder {
	return &errorBuilder{
		err: &Error{
			code:     code,
			message:  message,
			httpCode: httpCode,
			fields:   nil,
		},
	}
}

// Функции для создания ошибок с параметрами.
func Newf(code ErrorCode, httpCode int, format string, args ...interface{}) AppError {
	return &Error{
		code:     code,
		message:  fmt.Sprintf(format, args...),
		httpCode: httpCode,
		fields:   nil,
	}
}

// Вспомогательная функция для копирования полей.
func copyFields(fields map[string]interface{}) map[string]interface{} {
	if fields == nil {
		return nil
	}

	result := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		result[k] = v
	}

	return result
}

// Предопределенные ошибки (теперь как функции-фабрики).
func NotFound() AppError {
	return New(CodeNotFound, "not found", http.StatusNotFound)
}

func NotFoundf(format string, args ...interface{}) AppError {
	return Newf(CodeNotFound, http.StatusNotFound, format, args...)
}

func AlreadyExists() AppError {
	return New(CodeAlreadyExists, "already exists", http.StatusBadRequest)
}

func AlreadyExistsf(format string, args ...interface{}) AppError {
	return Newf(CodeAlreadyExists, http.StatusBadRequest, format, args...)
}

func EmailAlreadyExists() AppError {
	return New(CodeEmailAlreadyExists, "email already exists", http.StatusBadRequest)
}

func LoginNotRegistered() AppError {
	return New(CodeLoginNotRegistered, "login not registered", http.StatusBadRequest)
}

func WrongPassword() AppError {
	return New(CodeWrongPassword, "wrong password", http.StatusBadRequest)
}

func InternalError() AppError {
	return New(CodeInternalError, "internal error", http.StatusInternalServerError)
}

func InternalErrorf(format string, args ...interface{}) AppError {
	return Newf(CodeInternalError, http.StatusInternalServerError, format, args...)
}

func InvalidToken() AppError {
	return New(CodeInvalidToken, "invalid token", http.StatusUnauthorized)
}

func Forbidden() AppError {
	return New(CodeForbidden, "forbidden", http.StatusForbidden)
}

func Forbiddenf(format string, args ...interface{}) AppError {
	return Newf(CodeForbidden, http.StatusForbidden, format, args...)
}

func InvalidParams() AppError {
	return New(CodeInvalidParams, "invalid params", http.StatusBadRequest)
}

func InvalidParamsf(format string, args ...interface{}) AppError {
	return Newf(CodeInvalidParams, http.StatusBadRequest, format, args...)
}

func ProviderAlreadyConnected() AppError {
	return New(CodeProviderAlreadyConnected, "provider already connected", http.StatusBadRequest)
}

func ProviderAccountAlreadyLinked() AppError {
	return New(CodeProviderAccountAlreadyLinked, "provider account already linked", http.StatusBadRequest)
}

func TempIdNotFound() AppError {
	return New(CodeTempIdNotFound, "temp id not found", http.StatusNotFound)
}

func InvalidTelegramHash() AppError {
	return New(CodeInvalidTelegramHash, "invalid telegram hash", http.StatusBadRequest)
}

func TelegramIsNotLinked() AppError {
	return New(CodeTelegramIsNotLinked, "telegram is not linked", http.StatusBadRequest)
}

func NeedEndRegistration() AppError {
	return New(CodeNeedEndRegistration, "need end registration", http.StatusBadRequest)
}
