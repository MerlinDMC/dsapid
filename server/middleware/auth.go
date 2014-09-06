package middleware

import (
	"encoding/base64"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/go-martini/martini"
	"net"
	"net/http"
	"strings"
)

type User interface {
	GetId() string
	GetName() string
	HasRoles(...dsapid.UserRoleName) bool
	IsGuest() bool
	GetAuthInfo() interface{}
}

func Auth(user_storage storage.UserStorage) martini.Handler {
	return func(ctx martini.Context, res http.ResponseWriter, req *http.Request) {
		var user *dsapid.UserResource

		if h, ok := req.Header["Authorization"]; ok && len(h) > 0 {
			var token string

			parts := strings.SplitN(h[0], " ", 2)

			if parts[0] == "Basic" {
				if b, err := base64.StdEncoding.DecodeString(parts[1]); err == nil {
					parts = strings.Split(string(b), ":")
					if len(parts) == 2 {
						if len(parts[0]) > 0 {
							token = parts[0]
						} else if len(parts[1]) > 0 {
							token = parts[0]
						}
					}
				}
			}

			if token != "" {
				logger.Debugf("got auth token '%s'", token)

				if v, err := user_storage.FindByToken(token); err == nil {
					user = v

					logger.Debugf("found matching user %s (%s)", user.Uuid, user.Name)
				}
			}
		}

		if user == nil {
			user = user_storage.GuestUser()
		}

		ctx.MapTo(user, (*User)(nil))

		ctx.Next()
	}
}

func RequireRoles(roles ...dsapid.UserRoleName) martini.Handler {
	return func(res http.ResponseWriter, user User) {
		logger.Debugf("checking roles on user %s (%s)", user.GetId(), user.GetName())

		if user.IsGuest() || !user.HasRoles(roles...) {
			http.Error(res, "Not allowed", http.StatusUnauthorized)
		}
	}
}

func RequireAdmin() martini.Handler {
	return func(req *http.Request, res http.ResponseWriter, user User) {
		remote_host := req.RemoteAddr

		if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			remote_host = host
		}

		logger.Debugf("checking roles on user %s (%s) from %s", user.GetId(), user.GetName(), remote_host)

		if (user.IsGuest() || !user.HasRoles(dsapid.UserRoleAdmin)) &&
			remote_host != "127.0.0.1" && remote_host != "[::1]" {
			http.Error(res, "Not allowed", http.StatusUnauthorized)
		}
	}
}
