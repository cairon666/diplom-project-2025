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
)
