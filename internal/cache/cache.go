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
	orders      map[string]model.Order
	redisClient *redis.Client
	serviceUUID string
}

// NewCache returns new cache instance with redisdb client
func NewCache(ctx context.Context, cfg configs.Config, rCli *redis.Client) *OrderCache {
	var cache OrderCache
	cache.orders = make(map[string]model.Order)
	cache.redisClient = rCli
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
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
					err := json.Unmarshal([]byte(msgString), &order)
					if err != nil {
						fmt.Print(err)
					}
					cache.streamMsgHandler(msg["method"].(string), order)
				}
			}
		}
	}()
	return &cache
}

// GetOrder method firstly check cache for order and if it isn't there method goes to repository and save this order object in cache
func (orderCache OrderCache) GetOrder(ctx context.Context, cfg *configs.Config, orderID string) (model.Order, error) {
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
func (orderCache OrderCache) UpdateOrder(order model.Order) error {
	sendToStreamUpdateMsg(cfg, orderCache.redisClient, order)
}

// DeleteOrder delete order object from cache and repository and send message to redis stream
func (orderCache OrderCache) DeleteOrder(ctx context.Context, cfg *configs.Config, rps repository.Repository, orderID string) error {
	delete(orderCache.orders, orderID)
	go sendToStreamDeleteMsg(cfg, orderCache.redisClient, orderID)
	return rps.Delete(ctx, orderID)
}

func streamError(err error, method, userID string) {
	if err != nil {
		log.WithFields(log.Fields{
			"status": "operation failed",
			"err":    err,
			"method": method,
			"userid": userID,
		}).Info("redis stream info")
	} else {
		log.WithFields(log.Fields{
			"status": "operation successfully ended",
			"err":    err,
			"method": method,
			"userid": userID,
		}).Info("redis stream info")
	}
}
