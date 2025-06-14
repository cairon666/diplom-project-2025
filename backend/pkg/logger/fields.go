package logger

import (
	"log/slog"
	"time"
)

func String(key string, val string) Field {
	return slog.String(key, val)
}

func Int(key string, val int) Field {
	return slog.Int(key, val)
}

func Int32(key string, val int32) Field {
	return slog.Int(key, int(val))
}

func Int64(key string, val int64) Field {
	return slog.Int64(key, val)
}

func Float32(key string, val float32) Field {
	return slog.Float64(key, float64(val))
}

func Float64(key string, val float64) Field {
	return slog.Float64(key, val)
}

func Bool(key string, val bool) Field {
	return slog.Bool(key, val)
}

func Error(err error) Field {
	return slog.Any("error", err)
}

func Stack() Field {
	return slog.String("stack", "stack trace not implemented")
}

func Message(msg string) Field {
	return slog.String("message", msg)
}

func Duration(key string, val time.Duration) Field {
	return slog.Duration(key, val)
}

func Durationp(key string, val *time.Duration) Field {
	if val == nil {
		return slog.Any(key, nil)
	}

	return slog.Duration(key, *val)
}

func Dict(key string, val ...Field) Field {
	attrs := make([]any, len(val))
	for i, v := range val {
		attrs[i] = v
	}

	return slog.Group(key, attrs...)
}

func Time(key string, val time.Time) Field {
	return slog.Time(key, val)
}

func Timep(key string, val *time.Time) Field {
	if val == nil {
		return slog.Any(key, nil)
	}

	return slog.Time(key, *val)
}

func Any(key string, val interface{}) Field {
	return slog.Any(key, val)
}
