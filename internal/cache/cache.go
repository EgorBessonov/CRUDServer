// Package cache represents caching in application
package cache

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
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
	streamName  string
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
					cache.streamMessageHandler(msg["method"].(string), order)
				}
			}
		}
	}()
	return &cache
}

// Get method firstly check cache for order and if it isn't there method goes to repository and save this order object in cache
func (orderCache OrderCache) Get(orderID string) model.Order {
	return orderCache.orders[orderID]
}

//Save method
func (orderCache OrderCache) Save(order model.Order) error {
	sendToStreamMessage(orderCache.redisClient, orderCache.streamName, "save", order)
	return nil
}

// Update method update order objects in cache and repository and send message to redis stream
func (orderCache OrderCache) Update(order model.Order) error {
	sendToStreamMessage(orderCache.redisClient, orderCache.streamName, "update", order)
	return nil
}

// DeleteOrder method delete order object from cache and repository and send message to redis stream
func (orderCache OrderCache) Delete(orderID string) error {
	sendToStreamMessage(orderCache.redisClient, orderCache.streamName, "delete", orderID)
	return nil
}
