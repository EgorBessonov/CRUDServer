// Package repository replies for database access
package repository

import (
	"context"
	"CRUDServer/internal/models"
)

// Repository interface represent repository behavior
type Repository interface {
	Create(context.Context, models.Order) error
	Read(context.Context, string) (models.Order, error)
	Update(context.Context, models.Order) error
	Delete(context.Context, string) error
	CreateAuthUser(context.Context, *models.AuthUser) error
	GetAuthUser(context.Context, string) (models.AuthUser, error)
	GetAuthUserByID(context.Context, string) (models.AuthUser, error)
	UpdateAuthUser(ctx context.Context, email, refreshToken string) error
	CloseDBConnection() error
}
