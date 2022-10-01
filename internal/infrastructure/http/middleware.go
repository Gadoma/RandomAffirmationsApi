package http

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.respondError(w, fmt.Sprintf("There was a general error: %q", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}