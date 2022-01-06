// Package repository replies for database access
package repository

import (
	"context"
	"CRUDServer/internal/model"
)

// Repository interface represent repository behavior
type Repository interface {
	Create(context.Context, model.Order) error
	Read(context.Context, string) (model.Order, error)
	Update(context.Context, model.Order) error
	Delete(context.Context, string) error
	CreateAuthUser(context.Context, *model.AuthUser) error
	GetAuthUser(context.Context, string) (model.AuthUser, error)
	GetAuthUserByID(context.Context, string) (model.AuthUser, error)
	UpdateAuthUser(ctx context.Context, email, refreshToken string) error
	CloseDBConnection() error
}
