package middleware

import (
	log "github.com/MerlinDMC/logrus"
	"github.com/go-martini/martini"
	"net"
	"net/http"
)

func LogrusLogger() martini.Handler {
	return func(ctx martini.Context, user User, res http.ResponseWriter, req *http.Request) {
		rw := res.(martini.ResponseWriter)

		ctx.Next()

		host, _, err := net.SplitHostPort(req.RemoteAddr)

		if err != nil {
			host = req.RemoteAddr
		}

		log.WithFields(log.Fields{
			"remote_addr":         host,
			"user":                user.GetName(),
			"user_agent":          req.UserAgent(),
			"http_method":         req.Method,
			"http_request_uri":    req.URL.RequestURI(),
			"http_status":         rw.Status(),
			"http_content_length": rw.Size(),
		}).Info("request finished")
	}
}
