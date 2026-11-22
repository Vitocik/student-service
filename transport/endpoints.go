// Package transport определяет DTO и Go kit endpoints.
package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/student-service/student"
)

type CreateRequest struct {
	Student student.Student
}

type CreateResponse struct {
	ID string `json:"id"`
}

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Student student.Student `json:"student"`
}

type ChangeGroupRequest struct {
	ID    string
	Group string
}

type ChangeGroupResponse struct {
	Changed bool `json:"changed"`
}

type Endpoints struct {
	Create      endpoint.Endpoint
	GetByID     endpoint.Endpoint
	ChangeGroup endpoint.Endpoint
}

// MakeEndpoints адаптирует методы сервиса к универсальному типу Go kit Endpoint.
func MakeEndpoints(s student.Service) Endpoints {
	return Endpoints{
		Create: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(CreateRequest)
			if !ok {
				return nil, student.ErrInvalidRequest
			}
			id, err := s.Create(ctx, req.Student)
			if err != nil {
				return nil, err
			}
			return CreateResponse{ID: id}, nil
		},
		GetByID: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(GetByIDRequest)
			if !ok {
				return nil, student.ErrInvalidRequest
			}
			st, err := s.GetByID(ctx, req.ID)
			if err != nil {
				return nil, err
			}
			return GetByIDResponse{Student: st}, nil
		},
		ChangeGroup: func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(ChangeGroupRequest)
			if !ok {
				return nil, student.ErrInvalidRequest
			}
			if err := s.ChangeGroup(ctx, req.ID, req.Group); err != nil {
				return nil, err
			}
			return ChangeGroupResponse{Changed: true}, nil
		},
	}
}
