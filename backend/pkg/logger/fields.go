package logger

import (
	"time"

	"go.uber.org/zap"
)

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Error(err error) Field {
	return zap.Error(err)
}

func Stack() Field {
	return zap.Stack("stack")
}

func Message(msg string) Field {
	return zap.String("message", msg)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) Field {
	return zap.Durationp(key, val)
}

func Dict(key string, val ...Field) Field {
	return zap.Dict(key, val...)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) Field {
	return zap.Timep(key, val)
}

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}
