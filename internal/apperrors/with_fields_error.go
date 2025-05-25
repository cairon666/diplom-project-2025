package apperrors

import "net/http"

type WithFieldsError struct {
	errorCode ErrorCode
	httpCode  int
	message   string
	fields    map[string]interface{}
}

func NewWithFieldsError(code ErrorCode, message string, httpCode int) *WithFieldsError {
	return &WithFieldsError{
		errorCode: code,
		message:   message,
		httpCode:  httpCode,
		fields:    map[string]interface{}{},
	}
}

func (e *WithFieldsError) AddField(key string, value interface{}) {
	e.fields[key] = value
}

func (e *WithFieldsError) AddFields(fields map[string]interface{}) {
	for key, value := range fields {
		e.fields[key] = value
	}
}

func (e *WithFieldsError) Code() ErrorCode {
	return e.errorCode
}

func (e *WithFieldsError) Error() string {
	return e.message
}

func (e *WithFieldsError) HTTPCode() int {
	return e.httpCode
}

func (e *WithFieldsError) Clone() *WithFieldsError {
	return &WithFieldsError{
		errorCode: e.errorCode,
		message:   e.message,
		httpCode:  e.httpCode,
		fields:    e.fields,
	}
}

func (e *WithFieldsError) JSON() map[string]interface{} {
	return map[string]interface{}{
		"error":   e.errorCode,
		"message": e.message,
		"fields":  e.fields,
	}
}

var (
	ErrNeedEndRegistration = NewWithFieldsError(
		CodeNeedEndRegistration,
		"need end registration",
		http.StatusBadRequest,
	)
)
