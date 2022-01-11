package cache

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

func sendToStreamMessage(rCli *redis.Client, streamName, method string, data interface{}) {
	result := rCli.XAdd(&redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"method": method,
			"data":   data,
		},
	})
	_, err := result.Result()
	streamErrorCheck(err, "delete")
}

func (orderCache OrderCache) readStreamMessage(cfg *configs.Config, rCli *redis.Client) {
	result, err := rCli.XRead(&redis.XReadArgs{
		Streams: []string{cfg.StreamName, "0"},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"status": "failed",
			"err":    err,
		}).Info("redis stream info")
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString, ok := msg["data"].(string)
	if ok {
		order := model.Order{}
		err := order.UnmarshalBinary([]byte(msgString))
		if err != nil {
			fmt.Print(err)
		}
		orderCache.streamMessageHandler(msg["method"].(string), order)
	}
}

func (orderCache OrderCache) streamMessageHandler(method string, order model.Order) {
	switch method {
	case "save", "update":
		orderCache.orders[order.OrderID] = order
	case "delete":
		delete(orderCache.orders, order.OrderID)
	default:
		streamErrorCheck(errors.New("invalid type of method"), "streamMessageHandler")
	}
}

func streamErrorCheck(err error, method string) {
	if err != nil {
		log.WithFields(log.Fields{
			"status": "operation failed",
			"method": method,
			"error":  err,
		}).Info("redis stream info")
		return
	}
	log.WithFields(log.Fields{
		"status": "successfully ended",
		"method": method,
	}).Info("redis stream info")
}
