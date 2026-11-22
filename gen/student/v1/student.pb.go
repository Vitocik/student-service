// Code in this file implements the protobuf messages declared in proto/student.proto.
// It uses the protobuf v1 compatibility interface so it can be used by grpc's adapter.
// To replace it with canonical generated code, run scripts/generate-proto.ps1.
package studentv1

import proto "github.com/golang/protobuf/proto"

const _ = proto.ProtoPackageIsVersion4

type Student struct {
	Id               string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FullName         string  `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Group            string  `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	Gpa              float32 `protobuf:"fixed32,4,opt,name=gpa,proto3" json:"gpa,omitempty"`
	EnrolledOn       int64   `protobuf:"varint,5,opt,name=enrolled_on,json=enrolledOn,proto3" json:"enrolled_on,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Student) Reset()         { *m = Student{} }
func (m *Student) String() string { return proto.CompactTextString(m) }
func (*Student) ProtoMessage()    {}
func (m *Student) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}
func (m *Student) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}
func (m *Student) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}
func (m *Student) GetGpa() float32 {
	if m != nil {
		return m.Gpa
	}
	return 0
}
func (m *Student) GetEnrolledOn() int64 {
	if m != nil {
		return m.EnrolledOn
	}
	return 0
}

type CreateStudentRequest struct {
	Student          *Student `protobuf:"bytes,1,opt,name=student,proto3" json:"student,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CreateStudentRequest) Reset()         { *m = CreateStudentRequest{} }
func (m *CreateStudentRequest) String() string { return proto.CompactTextString(m) }
func (*CreateStudentRequest) ProtoMessage()    {}
func (m *CreateStudentRequest) GetStudent() *Student {
	if m != nil {
		return m.Student
	}
	return nil
}

type CreateStudentResponse struct {
	Id               string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CreateStudentResponse) Reset()         { *m = CreateStudentResponse{} }
func (m *CreateStudentResponse) String() string { return proto.CompactTextString(m) }
func (*CreateStudentResponse) ProtoMessage()    {}
func (m *CreateStudentResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetByIDRequest struct {
	Id               string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *GetByIDRequest) Reset()         { *m = GetByIDRequest{} }
func (m *GetByIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetByIDRequest) ProtoMessage()    {}
func (m *GetByIDRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetByIDResponse struct {
	Student          *Student `protobuf:"bytes,1,opt,name=student,proto3" json:"student,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *GetByIDResponse) Reset()         { *m = GetByIDResponse{} }
func (m *GetByIDResponse) String() string { return proto.CompactTextString(m) }
func (*GetByIDResponse) ProtoMessage()    {}
func (m *GetByIDResponse) GetStudent() *Student {
	if m != nil {
		return m.Student
	}
	return nil
}

type ChangeGroupRequest struct {
	Id               string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Group            string `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChangeGroupRequest) Reset()         { *m = ChangeGroupRequest{} }
func (m *ChangeGroupRequest) String() string { return proto.CompactTextString(m) }
func (*ChangeGroupRequest) ProtoMessage()    {}
func (m *ChangeGroupRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}
func (m *ChangeGroupRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

type ChangeGroupResponse struct {
	Changed          bool   `protobuf:"varint,1,opt,name=changed,proto3" json:"changed,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChangeGroupResponse) Reset()         { *m = ChangeGroupResponse{} }
func (m *ChangeGroupResponse) String() string { return proto.CompactTextString(m) }
func (*ChangeGroupResponse) ProtoMessage()    {}
func (m *ChangeGroupResponse) GetChanged() bool {
	if m != nil {
		return m.Changed
	}
	return false
}
