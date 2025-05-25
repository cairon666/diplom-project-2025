package apperrors

type AppError interface {
	JSON() map[string]interface{}
	HTTPCode() int
	error
}
