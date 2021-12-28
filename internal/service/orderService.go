package service

import (
	"CRUDServer/internal/models"
	"CRUDServer/internal/repository"
	"context"
)

// Save function ...
func Save(ctx context.Context, rps repository.Repository, order models.Order) error {
	return rps.Create(ctx, order)
}

// Get function ...
func Get(ctx context.Context, rps repository.Repository, userID string) (models.Order, error) {
	return rps.Read(ctx, userID)
}

// Delete function ...
func Delete(ctx context.Context, rps repository.Repository, userID string) error {
	return rps.Delete(ctx, userID)
}

// Update function...
func Update(ctx context.Context, rps repository.Repository, order models.Order) error {
	return rps.Update(ctx, order)
}
