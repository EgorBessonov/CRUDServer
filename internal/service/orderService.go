package service

import (
	"CRUDServer/internal/model"
	"context"
)

// Save function ...
func (s Service) Save(ctx context.Context, order model.Order) error {
	s.orderCache.Save(order)
	return s.rps.Create(ctx, order)
}

// Get function ...
func (s Service) Get(ctx context.Context, orderID string) (model.Order, error) {
	order := s.orderCache.Get(orderID)
	if order.OrderName == "" {
		order, err := s.rps.Read(ctx, orderID)
		if err != nil {
			return model.Order{}, nil
		}
		s.orderCache.Save(order)
		return order, nil
	}
	return order, nil
}

// Delete function ...
func (s Service) Delete(ctx context.Context, orderID string) error {
	s.orderCache.Delete(orderID)
	return s.rps.Delete(ctx, orderID)
}

// Update function...
func (s Service) Update(ctx context.Context, order model.Order) error {
	s.orderCache.Update(order)
	return s.rps.Update(ctx, order)
}
