// Package cache represents caching in application
package cache

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
	"github.com/go-redis/redis"
)

// OrderCache type represents cache object structure and behavior
type OrderCache struct {
	orders map[string]model.Order
	rCli   *redis.Client
}

// NewCache returns new cache instance with redisdb client
func NewCache(_rCli *redis.Client) *OrderCache {
	var cache OrderCache
	cache.orders = make(map[string]model.Order)
	cache.rCli = _rCli
	return &cache
}

// GetOrder method firstly check cache for order and if it isn't there method goes to repository and save this order object in cache
func (orderCache OrderCache) GetOrder(ctx context.Context, rps repository.Repository, orderID string) (model.Order, error) {
	order := orderCache.orders[orderID]
	if order.OrderName == "" {
		order, err := rps.Read(ctx, orderID)
		if err != nil {
			return model.Order{}, err
		}
		orderCache.orders[orderID] = order
		return order, nil
	}
	return order, nil
}

// UpdateOrder update order objects in cache and repository and send message to redis stream
func (orderCache OrderCache) UpdateOrder(ctx context.Context, rps repository.Repository, order model.Order) error {
	orderCache.orders[order.OrderID] = order
	return rps.Update(ctx, order)
}

// DeleteOrder delete order object from cache and repository and send message to redis stream
func (orderCache OrderCache) DeleteOrder(ctx context.Context, rps repository.Repository, orderID string) error {
	delete(orderCache.orders, orderID)
	return rps.Delete(ctx, orderID)
}

func (orderCache OrderCache) produceMsg(cfg configs.Config, event map[string]interface{}) (string, error) {
	return orderCache.rCli.XAdd(&redis.XAddArgs{
		Stream: cfg.StreamName,
		Values: event,
	}).Result()
}
