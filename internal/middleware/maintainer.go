package middleware

import "net/http"

func VerifyMaintainer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tp, err := ParseCredentials(r.Context())
		if err != nil {
			http.Error(w, "Authorization failed: invalid token", http.StatusUnauthorized)
			return
		}
		if !tp.IsMaintainer {
			http.Error(w, "This action is available for maintainer only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
