package service

import (
	"CRUDServer/internal/repository"
	"context"
)

// Save function ...
func Save(ctx context.Context, rps repository.Repository, userform repository.User) error {
	return rps.CreateUser(ctx, userform)
}
// Get function ...
func Get(ctx context.Context, rps repository.Repository, userID string) (repository.User, error) {
	return rps.ReadUser(ctx, userID)
}
// Delete function ...
func Delete(ctx context.Context, rps repository.Repository, userID string) error {
	return rps.DeleteUser(ctx, userID)
}
// Update function...
func Update(ctx context.Context, rps repository.Repository, userform repository.User) error {
	return rps.UpdateUser(ctx, userform)
}