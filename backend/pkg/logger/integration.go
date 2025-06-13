package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WithError создает поле для логирования ошибки.
func WithError(err error) Field {
	return Error(err)
}

// WithErrorCode создает поле для кода ошибки.
func WithErrorCode(code string) Field {
	return String("error_code", code)
}

// WithHTTPStatus создает поле для HTTP статуса.
func WithHTTPStatus(status int) Field {
	return Int("http_status", status)
}

// WithUserID создает поле для ID пользователя.
func WithUserID(userID int64) Field {
	return Int64("user_id", userID)
}

// WithRequestID создает поле для ID запроса.
func WithRequestID(requestID string) Field {
	return String("request_id", requestID)
}

// WithTraceID создает поле для trace ID.
func WithTraceID(traceID string) Field {
	return String("trace_id", traceID)
}

// WithOperation создает поле для названия операции.
func WithOperation(operation string) Field {
	return String("operation", operation)
}

// WithDuration создает поле для продолжительности операции.
func WithDuration(duration time.Duration) Field {
	return Duration("duration", duration)
}

// WithRequestInfo создает поля для информации о HTTP запросе.
func WithRequestInfo(method, path, userAgent string) []Field {
	return []Field{
		String("http_method", method),
		String("http_path", path),
		String("user_agent", userAgent),
	}
}

// ContextKeys для извлечения информации из контекста.
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	TraceIDKey   ContextKey = "trace_id"
)

// FromContext извлекает значения из контекста и создает поля для логгера.
func FromContext(ctx context.Context) []Field {
	var fields []Field

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		fields = append(fields, WithRequestID(requestID))
	}

	if userID, ok := ctx.Value(UserIDKey).(int64); ok && userID > 0 {
		fields = append(fields, WithUserID(userID))
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		fields = append(fields, WithTraceID(traceID))
	}

	return fields
}

// WithContextFields создает логгер с полями из контекста.
func WithContextFields(logger ILogger, ctx context.Context) ILogger {
	fields := FromContext(ctx)
	if len(fields) > 0 {
		return logger.With(fields...)
	}

	return logger
}

// GinMiddleware создает middleware для Gin с логированием запросов.
func GinMiddleware(logger ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Добавляем request ID в контекст
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Логируем начало запроса
		logger.Info("Request started",
			WithRequestID(requestID),
			String("method", c.Request.Method),
			String("path", c.Request.URL.Path),
			String("remote_addr", c.ClientIP()),
			String("user_agent", c.Request.UserAgent()),
		)

		// Выполняем запрос
		c.Next()

		// Логируем окончание запроса
		duration := time.Since(start)
		status := c.Writer.Status()

		logFields := []Field{
			WithRequestID(requestID),
			String("method", c.Request.Method),
			String("path", c.Request.URL.Path),
			Int("status", status),
			WithDuration(duration),
			Int("response_size", c.Writer.Size()),
		}

		if status >= 400 {
			logger.Error("Request completed with error", logFields...)
		} else {
			logger.Info("Request completed", logFields...)
		}
	}
}

// ContextExtractor извлекает дополнительную информацию из HTTP запроса.
func ContextExtractor(r *http.Request) context.Context {
	ctx := r.Context()

	// Извлекаем Request ID
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
	}

	// Извлекаем Trace ID
	if traceID := r.Header.Get("X-Trace-Id"); traceID != "" {
		ctx = context.WithValue(ctx, TraceIDKey, traceID)
	}

	return ctx
}

// generateRequestID генерирует уникальный ID для запроса.
func generateRequestID() string {
	// Простая реализация для демонстрации
	return time.Now().Format("20060102-150405.000000")
}

// LoggedService базовая структура для сервисов с логгером.
type LoggedService struct {
	logger    ILogger
	operation string
}

// NewLoggedService создает новый сервис с логгером.
func NewLoggedService(logger ILogger, operation string) *LoggedService {
	return &LoggedService{
		logger:    logger.With(WithOperation(operation)),
		operation: operation,
	}
}

// Logger возвращает логгер сервиса.
func (s *LoggedService) Logger() ILogger {
	return s.logger
}

// WithContext возвращает логгер с контекстом.
func (s *LoggedService) WithContext(ctx context.Context) ILogger {
	return WithContextFields(s.logger, ctx)
}

// LogError логирует ошибку с дополнительной информацией.
func (s *LoggedService) LogError(ctx context.Context, err error, msg string, fields ...Field) {
	allFields := []Field{WithError(err)}
	allFields = append(allFields, FromContext(ctx)...)
	allFields = append(allFields, fields...)

	s.logger.Error(msg, allFields...)
}

// LogInfo логирует информационное сообщение.
func (s *LoggedService) LogInfo(ctx context.Context, msg string, fields ...Field) {
	allFields := FromContext(ctx)
	allFields = append(allFields, fields...)

	s.logger.Info(msg, allFields...)
}
