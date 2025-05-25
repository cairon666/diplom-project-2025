package logger

import (
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
