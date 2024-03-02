package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// LoggerMiddleware — middleware-логер для входящих HTTP-запросов.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			statusCode: 0,
			sizeBytes:  0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Debug("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.Int("statusCode", responseData.statusCode),
			zap.Int64("responseSizeBytes", responseData.sizeBytes),
		)
	})
}

type responseData struct {
	statusCode int
	sizeBytes  int64
}

type loggingResponseWriter struct {
	ResponseWriter http.ResponseWriter
	responseData   *responseData
}

func (w loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w loggingResponseWriter) Write(data []byte) (int, error) {
	sizeBytes, err := w.ResponseWriter.Write(data)
	w.responseData.sizeBytes += int64(sizeBytes)
	return sizeBytes, err
}

func (w loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.statusCode = statusCode
}
