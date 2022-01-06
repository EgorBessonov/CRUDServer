package cache

import log "github.com/sirupsen/logrus"

func redisOperationError(method string, err error) {
	log.WithFields(log.Fields{
		"method": method,
		"err":    err,
	}).Info("redis info")
}
