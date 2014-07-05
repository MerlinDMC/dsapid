package middleware

import (
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/go-martini/martini"
	"net"
	"net/http"
)

func JsonLogger() martini.Handler {
	return func(ctx martini.Context, user User, res http.ResponseWriter, req *http.Request) {
		rw := res.(martini.ResponseWriter)

		ctx.Next()

		host, _, err := net.SplitHostPort(req.RemoteAddr)

		if err != nil {
			host = req.RemoteAddr
		}

		logger.Infof("%s %s \"%s %s\" %d %d",
			host,
			user.GetName(),
			req.Method,
			req.URL.RequestURI(),
			rw.Status(),
			rw.Size())
	}
}
