package middleware

import (
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net"
	"net/http"
	"time"
)

func CombinedLogger() martini.Handler {
	return func(ctx martini.Context, res http.ResponseWriter, req *http.Request, log *log.Logger) {
		ts := time.Now()

		rw := res.(martini.ResponseWriter)

		ctx.Next()

		username := "-"
		if req.URL.User != nil {
			if name := req.URL.User.Username(); name != "" {
				username = name
			}
		}

		host, _, err := net.SplitHostPort(req.RemoteAddr)

		if err != nil {
			host = req.RemoteAddr
		}

		fmt.Printf("%s - %s [%s] \"%s %s %s\" %d %d\n",
			host,
			username,
			ts.Format("02/Jan/2006:15:04:05 -0700"),
			req.Method,
			req.URL.RequestURI(),
			req.Proto,
			rw.Status(),
			rw.Size(),
		)
	}
}
