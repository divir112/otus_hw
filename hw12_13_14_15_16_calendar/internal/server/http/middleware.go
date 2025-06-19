package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Middleware struct {
	logger Logger
}

func NewMiddlware(logger Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		var ip string
		if remoteAddr != "" {
			ip = strings.Split(remoteAddr, ":")[0]
		}

		date := time.Now()
		formatTime := date.Format("25/Feb/2020:19:11:24 +0600")
		method := r.Method
		path := r.URL.Path
		version := r.Proto

		next.ServeHTTP(w, r)
		duration := time.Since(date).Microseconds()
		userAgent := r.UserAgent()
		m.logger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %s", ip, formatTime, method, path, version, duration, userAgent))
	})
}
