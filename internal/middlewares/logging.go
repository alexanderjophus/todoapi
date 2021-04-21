package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

func Logging(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			log.Infow("request",
				"method", r.Method,
				// TODO log other fields
				// status code
				// duration
			)
		})
	}
}
