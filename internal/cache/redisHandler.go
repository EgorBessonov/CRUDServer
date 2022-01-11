package cache

import (
	"CRUDServer/internal/model"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

func (orderCache *OrderCache) sendMessageToStream(method string, data interface{}) {
	result := orderCache.redisClient.XAdd(&redis.XAddArgs{
		Stream: orderCache.streamName,
		Values: map[string]interface{}{
			"method": method,
			"data":   data,
		},
	})
	if _, err := result.Result(); err != nil {
		log.WithFields(log.Fields{
			"method": "sendMessageTOStream",
			"error":  err,
		}).Error("redis stream error")
	}
}

func (orderCache *OrderCache) streamMessageHandler(method string, order model.Order) {
	orderCache.mutex.Lock()
	defer orderCache.mutex.Unlock()
	switch method {
	case "save", "update":
		orderCache.orders[order.OrderID] = order
		log.WithFields(log.Fields{
			"status": "cache updated by stream",
		}).Info("cache info")
	case "delete":
		delete(orderCache.orders, order.OrderID)
		log.WithFields(log.Fields{
			"status": "element was deleted by stream",
		}).Info("cache info")
	default:
		log.WithFields(log.Fields{
			"error":  "invalid method type",
			"method": "streamMessageHandler()",
		}).Error("Redis operation failed")
	}
}
