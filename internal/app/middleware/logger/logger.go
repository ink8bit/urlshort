package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Log(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Infoln(
				"method", r.Method,
				"path", r.URL.Path,
			)

			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			now := time.Now()

			defer func() {
				log.Infow(
					"request completed",
					"status", rw.Status(),
					"bytes", rw.BytesWritten(),
					"duration", time.Since(now),
				)
			}()

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}
