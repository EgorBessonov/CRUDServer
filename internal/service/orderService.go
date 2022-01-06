package service

import (
	"CRUDServer/internal/model"
	"context"
)

// Save function ...
func(s Service) Save(ctx context.Context, order model.Order) error {
	return s.rps.Create(ctx, order)
}

// Get function ...
func(s Service) Get(ctx context.Context, userID string) (model.Order, error) {
	return s.rps.Read(ctx, userID)
}

// Delete function ...
func(s Service) Delete(ctx context.Context, userID string) error {
	return s.rps.Delete(ctx, userID)
}

// Update function...
func(s Service) Update(ctx context.Context, order model.Order) error {
	return s.rps.Update(ctx, order)
}
