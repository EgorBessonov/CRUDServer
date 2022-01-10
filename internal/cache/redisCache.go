// Package cache represents caching in application
package cache

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// OrderCache type represents cache object structure and behavior
type OrderCache struct {
	orders map[string]model.Order
	redisClient   *redis.Client
	serviceUUID string
}

// NewCache returns new cache instance with redisdb client
func NewCache(_rCli *redis.Client, _serviceUUID string) *OrderCache {
	var cache OrderCache
	cache.orders = make(map[string]model.Order)
	cache.redisClient = _rCli
	cache.serviceUUID = _serviceUUID
	return &cache
}

// GetOrder method firstly check cache for order and if it isn't there method goes to repository and save this order object in cache
func (orderCache OrderCache) GetOrder(ctx context.Context, cfg *configs.Config, rps repository.Repository, orderID string) (model.Order, error) {
	order := orderCache.orders[orderID]
	if order.OrderName == "" {
		order, err := rps.Read(ctx, orderID)
		if err != nil {
			return model.Order{}, err
		}
		orderCache.orders[orderID] = order
		go sendToStreamAddMsg(cfg, orderCache.redisClient, order)
	}
	return order, nil
}

// UpdateOrder update order objects in cache and repository and send message to redis stream
func (orderCache OrderCache) UpdateOrder(ctx context.Context, cfg *configs.Config, rps repository.Repository, order model.Order) error {
	orderCache.orders[order.OrderID] = order
	go sendToStreamUpdateMsg(cfg, orderCache.redisClient, order)
	return rps.Update(ctx, order)
}

// DeleteOrder delete order object from cache and repository and send message to redis stream
func (orderCache OrderCache) DeleteOrder(ctx context.Context, cfg *configs.Config, rps repository.Repository, orderID string) error {
	delete(orderCache.orders, orderID)
	go sendToStreamDeleteMsg(cfg, orderCache.redisClient, orderID)
	return rps.Delete(ctx, orderID)
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
				serviceUUID: cfg.ServiceUUID,
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

func (orderCache OrderCache) ReadStreamMsg(cfg *configs.Config, rCli *redis.Client){
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
}

func streamError(err error, method, userID string) {
	if err != nil {
		log.WithFields(log.Fields{
			"status" : "operation failed",
			"err" : err,
			"method" : method,
			"userid" : userID, 
		}).Info("redis stream info")
	}else {
		log.WithFields(log.Fields{
			"status" : "operation successfully ended",
			"err" : err,
			"method" : method,
			"userid" : userID, 
		}).Info("redis stream info")
	}
}

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