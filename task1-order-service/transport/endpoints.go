package transport

import (
	"context"
	"example.com/order-service/order"
	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct{ Order order.Order }
type CreateResponse struct {
	ID string `json:"id"`
}
type GetByIDRequest struct{ ID string }
type GetByIDResponse struct {
	Order order.Order `json:"order"`
}
type ChangeStatusRequest struct {
	ID     string
	Status string
}
type ChangeStatusResponse struct {
	Changed bool `json:"changed"`
}
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeStatus endpoint.Endpoint
}

func MakeEndpoints(s order.Service) Endpoints {
	return Endpoints{
		Create: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(CreateRequest)
			if !ok {
				return nil, order.ErrInvalidOrder
			}
			id, err := s.Create(ctx, req.Order)
			if err != nil {
				return nil, err
			}
			return CreateResponse{ID: id}, nil
		},
		GetByID: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(GetByIDRequest)
			if !ok {
				return nil, order.ErrInvalidOrder
			}
			o, err := s.GetByID(ctx, req.ID)
			if err != nil {
				return nil, err
			}
			return GetByIDResponse{Order: o}, nil
		},
		ChangeStatus: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(ChangeStatusRequest)
			if !ok {
				return nil, order.ErrInvalidOrder
			}
			if err := s.ChangeStatus(ctx, req.ID, req.Status); err != nil {
				return nil, err
			}
			return ChangeStatusResponse{Changed: true}, nil
		},
	}
}
