package logger

import (
	"context"
)

// ILogger определяет интерфейс для логирования
type ILogger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	With(fields ...Field) ILogger
	WithContext(ctx context.Context) ILogger
}

// FieldLogger расширенный интерфейс с поддержкой контекста и дополнительных методов
type FieldLogger interface {
	ILogger
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
}

// ContextLogger интерфейс для логгера с контекстом
type ContextLogger interface {
	InfoContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)
	DebugContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
}

// NoOpLogger - заглушка для тестирования
type NoOpLogger struct{}

func (n NoOpLogger) Info(msg string, fields ...Field)            {}
func (n NoOpLogger) Error(msg string, fields ...Field)           {}
func (n NoOpLogger) Debug(msg string, fields ...Field)           {}
func (n NoOpLogger) Warn(msg string, fields ...Field)            {}
func (n NoOpLogger) With(fields ...Field) ILogger                { return n }
func (n NoOpLogger) WithContext(ctx context.Context) ILogger     { return n }
func (n NoOpLogger) Fatal(msg string, fields ...Field)           {}
func (n NoOpLogger) Panic(msg string, fields ...Field)           {}
func (n NoOpLogger) InfoContext(ctx context.Context, msg string, fields ...Field)  {}
func (n NoOpLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {}
func (n NoOpLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {}
func (n NoOpLogger) WarnContext(ctx context.Context, msg string, fields ...Field)  {}

// Ensure NoOpLogger implements all interfaces
var (
	_ ILogger       = (*NoOpLogger)(nil)
	_ FieldLogger   = (*NoOpLogger)(nil)
	_ ContextLogger = (*NoOpLogger)(nil)
) 