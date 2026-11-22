package order

import (
	"context"
	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order data")
)

type Order struct {
	ID           string      `json:"id,omitempty"`
	CustomerID   string      `json:"customer_id"`
	Status       string      `json:"status,omitempty"`
	CreatedOn    int64       `json:"created_on,omitempty"`
	RestaurantID string      `json:"restaurant_id"`
	OrderItems   []OrderItem `json:"order_items,omitempty"`
}

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

type Service interface {
	Create(context.Context, Order) (string, error)
	GetByID(context.Context, string) (Order, error)
	ChangeStatus(context.Context, string, string) error
}

type Repository interface {
	CreateOrder(context.Context, Order) error
	GetOrderByID(context.Context, string) (Order, error)
	ChangeOrderStatus(context.Context, string, string) error
}
