// Package cache represents caching in application
package cache

import (
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
	log "github.com/sirupsen/logrus"
)

type OrderCache struct {
	Orders map[string]model.Order
}

func NewCache() *OrderCache {
	var cache OrderCache
	cache.Orders = make(map[string]model.Order)
	return &cache
}

func (orderCache OrderCache) GetOrder(ctx context.Context, rps repository.Repository, orderID string) (model.Order, error) {
	order := orderCache.Orders[orderID]
	if order.OrderID == "" {
		order, err := rps.Read(ctx, orderID)
		if err != nil {
			return model.Order{}, err
		}
		orderCache.Orders[orderID] = order
		err = AddOrderToStream(order)
		if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"status": "failed",
			}).Info("redis streams info")
		}
		return order, nil
	}
	return order, nil
}

func (orderCache OrderCache) UpdateOrder(ctx context.Context, rps repository.Repository, order model.Order) error {
	orderCache.Orders[order.OrderID] = order
	if err := rps.Update(ctx, order); err != nil {
		return err
	}
	return nil
}
func AddOrderToStream(order model.Order) error {

	return nil
}
