package logger

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}

type Field = zap.Field

func NewProd() (*Logger, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	conf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{filepath.Join(cwd, "/logs/all.log")},
		ErrorOutputPaths: []string{filepath.Join(cwd, "/logs/err.log")},
	}
	l, err := conf.Build(zap.AddCallerSkip(1))

	return &Logger{
		l: l,
	}, err
}

func NewDev() (*Logger, error) {
	conf := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	l, err := conf.Build(zap.AddCallerSkip(1))

	return &Logger{
		l: l,
	}, err
}

// Основные методы логирования
func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

// Методы для работы с полями и контекстом
func (l *Logger) With(fields ...Field) ILogger {
	return &Logger{
		l: l.l.With(fields...),
	}
}

func (l *Logger) WithContext(ctx context.Context) ILogger {
	// Можно добавить извлечение информации из контекста
	// Например, trace_id, user_id и т.д.
	return l // Для базовой реализации возвращаем тот же логгер
}

// Методы для работы с контекстом
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

// Убеждаемся, что Logger реализует все интерфейсы
var (
	_ ILogger       = (*Logger)(nil)
	_ FieldLogger   = (*Logger)(nil)
	_ ContextLogger = (*Logger)(nil)
)
