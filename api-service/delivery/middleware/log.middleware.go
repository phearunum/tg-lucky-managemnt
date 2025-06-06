package middleware

import (
	"api-service/utils"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LoggerRequest(r.Body, r.URL.Path, "Request")

		next.ServeHTTP(w, r) // Call the next handler
	})
}
