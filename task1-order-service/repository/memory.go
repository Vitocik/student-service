package repository

import (
	"context"
	"strings"
	"sync"

	"example.com/order-service/order"
)

type memoryRepository struct {
	mu     sync.RWMutex
	orders map[string]order.Order
}

func NewMemoryRepository() order.Repository {
	return &memoryRepository{orders: make(map[string]order.Order)}
}

func (r *memoryRepository) CreateOrder(ctx context.Context, o order.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID] = o
	return nil
}

func (r *memoryRepository) GetOrderByID(ctx context.Context, id string) (order.Order, error) {
	if err := ctx.Err(); err != nil {
		return order.Order{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[strings.TrimSpace(id)]
	if !ok {
		return order.Order{}, order.ErrOrderNotFound
	}
	return o, nil
}

func (r *memoryRepository) ChangeOrderStatus(ctx context.Context, id, status string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[strings.TrimSpace(id)]
	if !ok {
		return order.ErrOrderNotFound
	}
	o.Status = strings.TrimSpace(status)
	r.orders[o.ID] = o
	return nil
}
