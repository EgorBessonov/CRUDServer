// Package cache represents caching in application
package cache

import (
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
)

type OrderCache struct {
	Orders map[string]model.Order
}

func (orderCache OrderCache) GetOrder(orderID string) model.Order {
	order := orderCache.Orders[orderID]

	return order
}

func (orderCache OrderCache) putOrderInCache(ctx context.Context, rps repository.Repository, orderId string) error {

	return nil
}
