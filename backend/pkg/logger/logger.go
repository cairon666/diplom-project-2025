package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	l *slog.Logger
}

const logFilePermissions = 0666

func NewProd(path string) (*Logger, error) {
	// Создаем файл для логов
	logFile, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePermissions)
	if err != nil {
		// Если не можем создать файл, используем stdout
		logFile = os.Stdout
	}

	// Создаем JSON handler для продакшена
	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	l := slog.New(handler)

	return &Logger{
		l: l,
	}, nil
}

func NewDev() (*Logger, error) {
	// Создаем текстовый handler для разработки
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	l := slog.New(handler)

	return &Logger{
		l: l,
	}, nil
}

// convertFields конвертирует []Field в []any для совместимости с slog.
func convertFields(fields []Field) []any {
	result := make([]any, len(fields))
	for i, field := range fields {
		result[i] = field
	}

	return result
}

// Основные методы логирования.
func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, convertFields(fields)...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, convertFields(fields)...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, convertFields(fields)...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, convertFields(fields)...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Error(msg, convertFields(fields)...)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Error(msg, convertFields(fields)...)
	panic(msg)
}

// Методы для работы с полями и контекстом.
func (l *Logger) With(fields ...Field) ILogger {
	return &Logger{
		l: l.l.With(convertFields(fields)...),
	}
}

func (l *Logger) WithContext(ctx context.Context) ILogger {
	// Можно добавить извлечение информации из контекста
	// Например, trace_id, user_id и т.д.
	return l // Для базовой реализации возвращаем тот же логгер
}

// InfoContext Методы для работы с контекстом.
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	l.l.InfoContext(ctx, msg, convertFields(fields)...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	l.l.ErrorContext(ctx, msg, convertFields(fields)...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	l.l.DebugContext(ctx, msg, convertFields(fields)...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	l.l.WarnContext(ctx, msg, convertFields(fields)...)
}

// Убеждаемся, что Logger реализует все интерфейсы.
var (
	_ ILogger       = (*Logger)(nil)
	_ FieldLogger   = (*Logger)(nil)
	_ ContextLogger = (*Logger)(nil)
)
