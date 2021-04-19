package middlewares

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	logger, err := zap.NewProduction()
	// panic setting up - shouldn't happen
	if err != nil {
		panic(fmt.Errorf("new logger: %w", err))
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		sugar.Infow("request",
			"method", r.Method,
			// TODO log other fields
			// status code
			// duration
		)
	})
}
