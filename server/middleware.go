package server

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) MeasureRequestTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now().Sub(start).Seconds()
		slog.Info("request", slog.Float64("time", end), slog.String("uri", r.Host+r.URL.Path))
	})
}
