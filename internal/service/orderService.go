package service

import (
	"CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"context"
)

// Save function ...
func (s Service) Save(ctx context.Context, order model.Order) error {
	return s.rps.Create(ctx, order)
}

// Get function ...
func (s Service) Get(ctx context.Context, cfg *configs.Config, orderID string) (model.Order, error) {
	return s.orderCache.GetOrder(ctx, cfg, s.rps, orderID)
}

// Delete function ...
func (s Service) Delete(ctx context.Context, orderID string) error {
	return s.rps.Delete(ctx, orderID)
}

// Update function...
func (s Service) Update(ctx context.Context, cfg *configs.Config, order model.Order) error {
	err := s.rps.Update(ctx, order)
	return s.orderCache.UpdateOrder(ctx, cfg, s.rps, order)
}
