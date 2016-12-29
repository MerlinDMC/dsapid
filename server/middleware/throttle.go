package middleware

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/throttle"
)

func Throttle(quota *throttle.Quota) martini.Handler {
	log.WithFields(log.Fields{
		"limit":  quota.Limit,
		"within": quota.Within,
	}).Info("adding throttling")

	return throttle.Policy(quota)
}
