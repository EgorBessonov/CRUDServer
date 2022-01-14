package service

import (
	"CRUDServer/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Save function method generate order uuid and after that save instance in cache and repository
func (s Service) Save(ctx context.Context, order *model.Order) (string, error) {
	order.OrderID = uuid.New().String()
	if err := s.orderCache.Save(order); err != nil {
		return "", fmt.Errorf("service: can't create order - %w", err)
	}
	if err := s.rps.Save(ctx, order); err != nil {
		return "", fmt.Errorf("service: can't create order - %w", err)
	}
	return order.OrderID, nil
}

// Get method look through cache for order and if order wasn't found, method get it from repository and add it in cache
func (s Service) Get(ctx context.Context, orderID string) (*model.Order, error) {
	order := s.orderCache.Get(orderID)
	if order == nil {
		order, err := s.rps.Get(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("service: can't get order - %w", err)
		}
		if err := s.orderCache.Save(order); err != nil {
			return nil, fmt.Errorf("service: can't get order - %w", err)
		}
		return order, nil
	}
	return order, nil
}

// Delete method delete order from repository and cache
func (s Service) Delete(ctx context.Context, orderID string) error {
	if err := s.orderCache.Delete(orderID); err != nil {
		return fmt.Errorf("service: can't delete order - %w", err)
	}
	if err := s.rps.Delete(ctx, orderID); err != nil {
		return fmt.Errorf("service: can't delete order - %w", err)
	}
	return nil
}

// Update method update order instance in repository and cache
func (s Service) Update(ctx context.Context, order *model.Order) error {
	if err := s.orderCache.Update(order); err != nil {
		return fmt.Errorf("service: can't update order - %w", err)
	}
	if err := s.rps.Update(ctx, order); err != nil {
		return fmt.Errorf("service: can't update order - %w", err)
	}
	return nil
}
