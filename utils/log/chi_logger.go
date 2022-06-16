package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func StatusLevel(status int) zapcore.Level {
	switch {
	case status <= 0:
		return zapcore.WarnLevel
	case status < 400: // for codes in 100s, 200s, 300s
		return zapcore.InfoLevel
	case status >= 400 && status < 500:
		return zapcore.WarnLevel
	case status >= 500:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func StatusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}

func NewChiLogger(servername, env string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			now := time.Now()
			defer func() {
				ZapLogger.Info("Served",
					zap.String("service", servername),
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("user-agent", r.UserAgent()),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
					zap.String("env", env),
					zap.String("timestamp", now.Format(time.RFC3339)),
					zap.Duration("duration", time.Since(now)),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
