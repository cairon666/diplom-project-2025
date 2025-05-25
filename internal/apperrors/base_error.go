package apperrors

import "net/http"

type BaseError struct {
	code     ErrorCode
	message  string
	httpCode int
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Error() string {
	return e.message
}

func (e *BaseError) HTTPCode() int {
	return e.httpCode
}

func NewBaseError(code ErrorCode, message string, httpCode int) *BaseError {
	return &BaseError{
		code:     code,
		message:  message,
		httpCode: httpCode,
	}
}

var (
	ErrNotFound                     = NewBaseError(CodeNotFound, "not found", http.StatusNotFound)
	ErrAlreadyExists                = NewBaseError(CodeAlreadyExists, "already exists", http.StatusBadRequest)
	ErrEmailAlreadyExists           = NewBaseError(CodeEmailAlreadyExists, "email already exists", http.StatusBadRequest)
	ErrLoginNotRegistered           = NewBaseError(CodeLoginNotRegistered, "login not registered", http.StatusBadRequest)
	ErrWrongPassword                = NewBaseError(CodeWrongPassword, "wrong password", http.StatusBadRequest)
	ErrInternalError                = NewBaseError(CodeInternalError, "internal error", http.StatusInternalServerError)
	ErrInvalidToken                 = NewBaseError(CodeInvalidToken, "invalid token", http.StatusUnauthorized)
	ErrForbidden                    = NewBaseError(CodeForbidden, "forbidden", http.StatusForbidden)
	ErrInvalidParams                = NewBaseError(CodeInvalidParams, "invalid params", http.StatusBadRequest)
	ErrProviderAlreadyConnected     = NewBaseError(CodeProviderAlreadyConnected, "provider already connected", http.StatusBadRequest)
	ErrProviderAccountAlreadyLinked = NewBaseError(CodeProviderAccountAlreadyLinked, "provider account already linked", http.StatusBadRequest)
	ErrTempIdNotFound               = NewBaseError(CodeTempIdNotFound, "temp id not found", http.StatusNotFound)
	ErrInvalidTelegramHash          = NewBaseError(CodeInvalidTelegramHash, "invalid telegram hash", http.StatusBadRequest)
	ErrTelegramIsNotLinked          = NewBaseError(CodeTelegramIsNotLinked, "telegram is not linked", http.StatusBadRequest)
)
