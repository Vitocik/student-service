package implementation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"example.com/order-service/order"
)

type service struct {
	repository order.Repository
	logger     log.Logger
}

func NewService(repo order.Repository, logger log.Logger) order.Service {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &service{repository: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, o order.Order) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	if strings.TrimSpace(o.CustomerID) == "" || strings.TrimSpace(o.RestaurantID) == "" {
		return "", order.ErrInvalidOrder
	}
	o.ID = uuid.NewString()
	o.Status = "Pending"
	o.CreatedOn = time.Now().Unix()
	if err := s.repository.CreateOrder(ctx, o); err != nil {
		level.Error(logger).Log("err", err)
		return "", fmt.Errorf("create order: %w", err)
	}
	return o.ID, nil
}

func (s *service) GetByID(ctx context.Context, id string) (order.Order, error) {
	o, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}
	return o, nil
}

func (s *service) ChangeStatus(ctx context.Context, id, status string) error {
	if strings.TrimSpace(status) == "" {
		return order.ErrInvalidOrder
	}
	return s.repository.ChangeOrderStatus(ctx, id, status)
}
