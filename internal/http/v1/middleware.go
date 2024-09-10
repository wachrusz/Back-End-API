package handlers

import (
	"net/http"
	"strings"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		contentType = strings.Split(contentType, ";")[0]

		w.Header().Set("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}
