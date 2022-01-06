package service

import (
	"CRUDServer/internal/cache"
	"CRUDServer/internal/model"
	"context"
)

// Save function ...
func (s Service) Save(ctx context.Context, order model.Order) error {
	return s.rps.Create(ctx, order)
}

// Get function ...
func (s Service) Get(ctx context.Context, cache cache.OrderCache, orderID string) (model.Order, error) {
	return cache.GetOrder(ctx, s.rps, orderID)
}

// Delete function ...
func (s Service) Delete(ctx context.Context, orderID string) error {
	return s.rps.Delete(ctx, orderID)
}

// Update function...
func (s Service) Update(ctx context.Context, cache cache.OrderCache, order model.Order) error {
	return cache.UpdateOrder(ctx, s.rps, order)
}
