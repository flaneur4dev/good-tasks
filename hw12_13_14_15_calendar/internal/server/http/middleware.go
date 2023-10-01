package httpserver

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)

		log := fmt.Sprintf("%s [%s] %s %s %s %d %s",
			ip, t.String(),
			r.Method, r.URL.RequestURI(), r.Proto,
			time.Since(t).Milliseconds(),
			r.UserAgent(),
		)

		_, err = s.fd.Write([]byte(log + "\n"))
		if err != nil {
			s.log.Error("failed to write to logfile: " + err.Error())
		}
	})
}
