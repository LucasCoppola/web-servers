package handlers

import (
	"net/http"
)

func (apiCfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.FileServerHits++
		next.ServeHTTP(w, r)
	})
}
