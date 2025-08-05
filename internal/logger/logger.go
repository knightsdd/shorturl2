package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

var Log *zap.Logger = zap.NewNop()

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(d.Microseconds())
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig = encoderConfig

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

// Logger — middleware-логер для входящих HTTP-запросов.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := &responseData{}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()
		next.ServeHTTP(lw, r)
		duration := time.Since(start)

		Log.Info(
			"HTTP REQUEST",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
		)
		Log.Info(
			"HTTP RESPONSE",
			zap.Int("status", lw.responseData.status),
			zap.Int("size", lw.responseData.size),
		)
	})
}
