package logging

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

type middlewareLogging func(next http.Handler) http.Handler

func New(log *slog.Logger) middlewareLogging {
	return func(next http.Handler) http.Handler {
		const op = "middleware.logging"

		logger := log.With(
			slog.String("op", op),
		)

		logger.Info("middleware logger enabled")

		handler := func(w http.ResponseWriter, r *http.Request) {
			requestMessage := logger.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_address", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("requestID", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				requestMessage.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("time", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(handler)
	}
}
