// Package grpctransport предоставляет тот же сервис через gRPC.
package grpctransport

import (
	"context"
	"errors"
	"fmt"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	studentv1 "example.com/student-service/gen/student/v1"
	"example.com/student-service/student"
	"example.com/student-service/transport"
)

type grpcServer struct {
	studentv1.UnimplementedStudentServiceServer
	create      *kitgrpc.Server
	getByID     *kitgrpc.Server
	changeGroup *kitgrpc.Server
}

// NewService связывает Go kit endpoints с типизированным gRPC-сервисом.
func NewService(endpoints transport.Endpoints) studentv1.StudentServiceServer {
	return &grpcServer{
		create:      kitgrpc.NewServer(endpoints.Create, decodeCreateRequest, encodeCreateResponse),
		getByID:     kitgrpc.NewServer(endpoints.GetByID, decodeGetByIDRequest, encodeGetByIDResponse),
		changeGroup: kitgrpc.NewServer(endpoints.ChangeGroup, decodeChangeGroupRequest, encodeChangeGroupResponse),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *studentv1.CreateStudentRequest) (*studentv1.CreateStudentResponse, error) {
	_, response, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*studentv1.CreateStudentResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *studentv1.GetByIDRequest) (*studentv1.GetByIDResponse, error) {
	_, response, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*studentv1.GetByIDResponse), nil
}

func (s *grpcServer) ChangeGroup(ctx context.Context, req *studentv1.ChangeGroupRequest) (*studentv1.ChangeGroupResponse, error) {
	_, response, err := s.changeGroup.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*studentv1.ChangeGroupResponse), nil
}

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*studentv1.CreateStudentRequest)
	if !ok || req == nil || req.Student == nil {
		return nil, student.ErrInvalidRequest
	}
	return transport.CreateRequest{Student: student.Student{
		FullName:   req.Student.FullName,
		Group:      req.Student.Group,
		GPA:        req.Student.Gpa,
		EnrolledOn: req.Student.EnrolledOn,
	}}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(transport.CreateResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected create response type %T", response)
	}
	return &studentv1.CreateStudentResponse{Id: resp.ID}, nil
}

func decodeGetByIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*studentv1.GetByIDRequest)
	if !ok || req == nil {
		return nil, student.ErrInvalidRequest
	}
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(transport.GetByIDResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected get response type %T", response)
	}
	st := resp.Student
	return &studentv1.GetByIDResponse{Student: &studentv1.Student{
		Id: st.ID, FullName: st.FullName, Group: st.Group, Gpa: st.GPA, EnrolledOn: st.EnrolledOn,
	}}, nil
}

func decodeChangeGroupRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*studentv1.ChangeGroupRequest)
	if !ok || req == nil {
		return nil, student.ErrInvalidRequest
	}
	return transport.ChangeGroupRequest{ID: req.Id, Group: req.Group}, nil
}

func encodeChangeGroupResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(transport.ChangeGroupResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected change-group response type %T", response)
	}
	return &studentv1.ChangeGroupResponse{Changed: resp.Changed}, nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, student.ErrInvalidRequest), errors.Is(err, student.ErrInvalidStudent):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, student.ErrStudentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal service error")
	}
}
