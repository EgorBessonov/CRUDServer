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
	"sync"
)

// OrderCache type represents cache object structure and behavior
type OrderCache struct {
	orders      map[string]*model.Order
	redisClient *redis.Client
	streamName  string
	mutex       sync.Mutex
}

// NewCache returns new cache instance with redisdb client
func NewCache(ctx context.Context, cfg configs.Config, rCli *redis.Client) *OrderCache {
	var cache OrderCache
	cache.orders = make(map[string]*model.Order)
	cache.redisClient = rCli
	cache.streamName = cfg.StreamName
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := rCli.XRead(&redis.XReadArgs{
					Streams: []string{cfg.StreamName, "$"},
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
					if err := cache.streamMessageHandler(msg["method"].(string), &order); err != nil {
						log.Errorf("stream error - %e", err)
					}
				}
			}
		}
	}()
	return &cache
}

// Get method firstly check cache for order and if it isn't there method goes to repository and save this order object in cache
func (orderCache *OrderCache) Get(orderID string) *model.Order {
	orderCache.mutex.Lock()
	defer orderCache.mutex.Unlock()
	return orderCache.orders[orderID]
}

//Save method
func (orderCache *OrderCache) Save(order *model.Order) error {
	return orderCache.sendMessageToStream("save", order)
}

// Update method update order objects in cache and repository and send message to redis stream
func (orderCache *OrderCache) Update(order *model.Order) error {
	return orderCache.sendMessageToStream("update", order)
}

// Delete method delete order object from cache and repository and send message to redis stream
func (orderCache *OrderCache) Delete(orderID string) error {
	return orderCache.sendMessageToStream("delete", orderID)
}

func (orderCache *OrderCache) sendMessageToStream(method string, data interface{}) error {
	result := orderCache.redisClient.XAdd(&redis.XAddArgs{
		Stream: orderCache.streamName,
		Values: map[string]interface{}{
			"method": method,
			"data":   data,
		},
	})
	if _, err := result.Result(); err != nil {
		return fmt.Errorf("cache: can't send message  to stream - %w", err)
	}
	return nil
}

func (orderCache *OrderCache) streamMessageHandler(method string, order *model.Order) error {
	orderCache.mutex.Lock()
	defer orderCache.mutex.Unlock()
	switch method {
	case "save", "update":
		orderCache.orders[order.OrderID] = order
		return nil
	case "delete":
		delete(orderCache.orders, order.OrderID)
		return nil
	default:
		return fmt.Errorf("cache handler: invalid method type")
	}
}
