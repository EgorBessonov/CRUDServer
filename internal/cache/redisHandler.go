package cache

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"os"
)

type redisMsg struct{
	*model.Order
	serviceUUID string
}

func (msg redisMsg) MarshalBinary()([]byte, error){
	return json.Marshal(msg)
}

func (msg redisMsg) UnmarshalBinary(data []byte) error{
	if err := json.Unmarshal(data, &msg); err != nil{
		return err
	}
	return nil
}

func sendToStreamAddMsg(cfg *configs.Config, rCli *redis.Client, order model.Order) {
	result := rCli.XAdd(&redis.XAddArgs{
		Stream: cfg.StreamName,
		Values: map[string]interface{}{
			"method": "save",
			"data": &redisMsg{
				Order: &model.Order{
					OrderID: order.OrderID,
					OrderName: order.OrderName,
					OrderCost: order.OrderCost,
					IsDelivered: order.IsDelivered,
				},
				serviceUUID: cfg.ServiceUUID,
			},
		},
	})
	_, err := result.Result()
	streamError(err, "delete", cfg.ServiceUUID)
}

func sendToStreamUpdateMsg(cfg *configs.Config, rCli *redis.Client, order model.Order) {
	result := rCli.XAdd(&redis.XAddArgs{
		Stream: cfg.StreamName,
		Values: map[string]interface{}{
			"method":"update",
			"data": &redisMsg{
				Order: &model.Order{
					OrderID: order.OrderID,
					OrderName: order.OrderName,
					OrderCost: order.OrderCost,
					IsDelivered: order.IsDelivered,
				},
			},
		},
	})
	_, err := result.Result()
	streamError(err, "delete", cfg.ServiceUUID)
}

func sendToStreamDeleteMsg(cfg *configs.Config, rCli *redis.Client, orderID string) {
	result := rCli.XAdd(&redis.XAddArgs{
		Stream: cfg.StreamName,
		Values: map[string]interface{}{
			"method":"delete",
			"data": &redisMsg{
				Order: &model.Order{
					OrderID: orderID,
				},
				serviceUUID: cfg.ServiceUUID,

			},
		},
	})
	fmt.Print(cfg.ServiceUUID)
	_, err := result.Result()
	streamError(err, "delete", cfg.ServiceUUID)
}

func (orderCache OrderCache) readStreamMsg(cfg *configs.Config, rCli *redis.Client){
	result, err := rCli.XRead(&redis.XReadArgs{
		Streams: []string{cfg.StreamName, "0"},
		Count: 1,
		Block: 0,
	}).Result()
	if err != nil{
		log.WithFields(log.Fields{
			"status" : "failed",
			"err": err,
		}).Info("redis stream info")
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString, ok := msg["data"].(string)
	if ok{
		order := model.Order{}
		err := json.Unmarshal([]byte(msgString), &order)
		if err != nil{
			fmt.Print(err)
		}
		orderCache.streamMsgHandler(msg["method"].(string), order)
	}
}

func (orderCache OrderCache)streamMsgHandler(method string, order model.Order){
	switch method{
	case "save", "update":
		orderCache.orders[order.OrderID] = order
		streamError(nil, method, "")
	case "delete":
		delete(orderCache.orders, order.OrderID)
		streamError(nil, method, "")
	}

func redisOperationError(method string, err error) {
	log.WithFields(log.Fields{
		"method": method,
		"userid": os.Getenv("SERVICEUUID"),
		"err":    err,
	}).Info("redis info")
}
