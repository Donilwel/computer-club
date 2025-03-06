package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func LoggerMiddleware(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			log.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"time":   time.Since(start),
			}).Info("HTTP Request")
		})
	}
}
