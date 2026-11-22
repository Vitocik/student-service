// Code in this file implements the gRPC service declared in proto/student.proto.
// To replace it with canonical generated code, run scripts/generate-proto.ps1.
package studentv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	StudentService_Create_FullMethodName      = "/student.v1.StudentService/Create"
	StudentService_GetByID_FullMethodName     = "/student.v1.StudentService/GetByID"
	StudentService_ChangeGroup_FullMethodName = "/student.v1.StudentService/ChangeGroup"
)

type StudentServiceClient interface {
	Create(ctx context.Context, in *CreateStudentRequest, opts ...grpc.CallOption) (*CreateStudentResponse, error)
	GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error)
	ChangeGroup(ctx context.Context, in *ChangeGroupRequest, opts ...grpc.CallOption) (*ChangeGroupResponse, error)
}

type studentServiceClient struct{ cc grpc.ClientConnInterface }

func NewStudentServiceClient(cc grpc.ClientConnInterface) StudentServiceClient {
	return &studentServiceClient{cc: cc}
}
func (c *studentServiceClient) Create(ctx context.Context, in *CreateStudentRequest, opts ...grpc.CallOption) (*CreateStudentResponse, error) {
	out := new(CreateStudentResponse)
	err := c.cc.Invoke(ctx, StudentService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *studentServiceClient) GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error) {
	out := new(GetByIDResponse)
	err := c.cc.Invoke(ctx, StudentService_GetByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *studentServiceClient) ChangeGroup(ctx context.Context, in *ChangeGroupRequest, opts ...grpc.CallOption) (*ChangeGroupResponse, error) {
	out := new(ChangeGroupResponse)
	err := c.cc.Invoke(ctx, StudentService_ChangeGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type StudentServiceServer interface {
	Create(context.Context, *CreateStudentRequest) (*CreateStudentResponse, error)
	GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error)
	ChangeGroup(context.Context, *ChangeGroupRequest) (*ChangeGroupResponse, error)
	mustEmbedUnimplementedStudentServiceServer()
}

type UnimplementedStudentServiceServer struct{}

func (UnimplementedStudentServiceServer) Create(context.Context, *CreateStudentRequest) (*CreateStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedStudentServiceServer) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedStudentServiceServer) ChangeGroup(context.Context, *ChangeGroupRequest) (*ChangeGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeGroup not implemented")
}
func (UnimplementedStudentServiceServer) mustEmbedUnimplementedStudentServiceServer() {}

type UnsafeStudentServiceServer interface {
	mustEmbedUnimplementedStudentServiceServer()
}

func RegisterStudentServiceServer(s grpc.ServiceRegistrar, srv StudentServiceServer) {
	s.RegisterService(&StudentService_ServiceDesc, srv)
}

func _StudentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: StudentService_Create_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).Create(ctx, req.(*CreateStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _StudentService_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: StudentService_GetByID_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).GetByID(ctx, req.(*GetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _StudentService_ChangeGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).ChangeGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: StudentService_ChangeGroup_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).ChangeGroup(ctx, req.(*ChangeGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var StudentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "student.v1.StudentService",
	HandlerType: (*StudentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Create", Handler: _StudentService_Create_Handler},
		{MethodName: "GetByID", Handler: _StudentService_GetByID_Handler},
		{MethodName: "ChangeGroup", Handler: _StudentService_ChangeGroup_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/student.proto",
}
