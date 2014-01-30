package middleware

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func AllowCORS() martini.Handler {
	return func(res http.ResponseWriter) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, X-Requested-With, X-CSRF-Token")
	}
}
