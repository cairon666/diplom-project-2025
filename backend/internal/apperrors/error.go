package apperrors

// AppError - основной интерфейс для всех ошибок приложения
type AppError interface {
	error
	Code() ErrorCode
	HTTPCode() int
	JSON() map[string]interface{}
	WithMessage(message string) AppError
	WithField(key string, value interface{}) AppError
	WithFields(fields map[string]interface{}) AppError
}

// ErrorBuilder - интерфейс для создания ошибок
type ErrorBuilder interface {
	Build() AppError
	WithMessage(message string) ErrorBuilder
	WithField(key string, value interface{}) ErrorBuilder
	WithFields(fields map[string]interface{}) ErrorBuilder
}
