package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

func LogRequest(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Msgf("Request:  method  %s, URL  %s",
				r.Method, r.URL)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
