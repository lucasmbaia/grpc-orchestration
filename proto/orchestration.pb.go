// Code generated by protoc-gen-go.
// source: orchestration.proto
// DO NOT EDIT!

/*
Package orchestration is a generated protocol buffer package.

It is generated from these files:
	orchestration.proto

It has these top-level messages:
	Task
	Result
*/
package orchestration

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Task struct {
	Tracking   string `protobuf:"bytes,1,opt,name=tracking" json:"tracking,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Parameters string `protobuf:"bytes,3,opt,name=parameters" json:"parameters,omitempty"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Task) GetTracking() string {
	if m != nil {
		return m.Tracking
	}
	return ""
}

func (m *Task) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Task) GetParameters() string {
	if m != nil {
		return m.Parameters
	}
	return ""
}

type Result struct {
	Response map[string]string `protobuf:"bytes,1,rep,name=response" json:"response,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetResponse() map[string]string {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "orchestration.Task")
	proto.RegisterType((*Result)(nil), "orchestration.Result")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for OrchestrationService service

type OrchestrationServiceClient interface {
	CallTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*Result, error)
	Health(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type orchestrationServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrchestrationServiceClient(cc *grpc.ClientConn) OrchestrationServiceClient {
	return &orchestrationServiceClient{cc}
}

func (c *orchestrationServiceClient) CallTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/orchestration.OrchestrationService/CallTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestrationServiceClient) Health(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/orchestration.OrchestrationService/Health", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrchestrationService service

type OrchestrationServiceServer interface {
	CallTask(context.Context, *Task) (*Result, error)
	Health(context.Context, *google_protobuf1.Empty) (*google_protobuf1.Empty, error)
}

func RegisterOrchestrationServiceServer(s *grpc.Server, srv OrchestrationServiceServer) {
	s.RegisterService(&_OrchestrationService_serviceDesc, srv)
}

func _OrchestrationService_CallTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestrationServiceServer).CallTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orchestration.OrchestrationService/CallTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestrationServiceServer).CallTask(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrchestrationService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestrationServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orchestration.OrchestrationService/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestrationServiceServer).Health(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrchestrationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orchestration.OrchestrationService",
	HandlerType: (*OrchestrationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallTask",
			Handler:    _OrchestrationService_CallTask_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _OrchestrationService_Health_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "orchestration.proto",
}

func init() { proto.RegisterFile("orchestration.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x4d, 0x5b, 0x4b, 0x1d, 0x2d, 0xca, 0xb4, 0xca, 0x12, 0x45, 0xca, 0x7a, 0xe9, 0x69,
	0x17, 0xeb, 0xa5, 0xe8, 0x41, 0x50, 0x0a, 0x1e, 0x04, 0xa1, 0x8a, 0x07, 0x6f, 0x69, 0x19, 0xdb,
	0xa5, 0x69, 0xb2, 0x24, 0x69, 0xa1, 0x57, 0x0f, 0xbe, 0x80, 0x77, 0x5f, 0xca, 0x57, 0xf0, 0x41,
	0x64, 0xb3, 0x6d, 0xe9, 0x8a, 0x9e, 0x32, 0xf3, 0xfd, 0x33, 0xc9, 0x3f, 0x13, 0x68, 0x68, 0x33,
	0x1c, 0x93, 0x75, 0x46, 0xb8, 0x44, 0xab, 0x28, 0x35, 0xda, 0x69, 0xac, 0x17, 0x20, 0x3f, 0x19,
	0x69, 0x3d, 0x92, 0x14, 0x8b, 0x34, 0x89, 0x85, 0x52, 0xda, 0x79, 0x6c, 0xf3, 0x62, 0x7e, 0xbc,
	0x54, 0x7d, 0x36, 0x98, 0xbd, 0xc6, 0x34, 0x4d, 0xdd, 0x22, 0x17, 0xc3, 0x67, 0xa8, 0x3c, 0x09,
	0x3b, 0x41, 0x0e, 0x35, 0x67, 0xc4, 0x70, 0x92, 0xa8, 0x51, 0xc0, 0x5a, 0xac, 0xbd, 0xd3, 0x5f,
	0xe7, 0x88, 0x50, 0x51, 0x62, 0x4a, 0x41, 0xc9, 0x73, 0x1f, 0xe3, 0x29, 0x40, 0x2a, 0x8c, 0x98,
	0x92, 0x23, 0x63, 0x83, 0xb2, 0x57, 0x36, 0x48, 0xf8, 0xce, 0xa0, 0xda, 0x27, 0x3b, 0x93, 0x0e,
	0xaf, 0xa1, 0x66, 0xc8, 0xa6, 0x5a, 0x59, 0x0a, 0x58, 0xab, 0xdc, 0xde, 0xed, 0x9c, 0x45, 0xc5,
	0xa1, 0xf2, 0xc2, 0xec, 0xf0, 0x55, 0x3d, 0xe5, 0xcc, 0xa2, 0xbf, 0x6e, 0xe2, 0x57, 0x50, 0x2f,
	0x48, 0x78, 0x00, 0xe5, 0x09, 0x2d, 0x96, 0x3e, 0xb3, 0x10, 0x9b, 0xb0, 0x3d, 0x17, 0x72, 0xb6,
	0xf2, 0x98, 0x27, 0x97, 0xa5, 0x2e, 0xeb, 0x7c, 0x32, 0x68, 0x3e, 0x6c, 0xbe, 0xf6, 0x48, 0x66,
	0x9e, 0x0c, 0x09, 0xbb, 0x50, 0xbb, 0x15, 0x52, 0xfa, 0xe9, 0x1b, 0xbf, 0x0c, 0x65, 0x90, 0x1f,
	0xfe, 0xe9, 0x32, 0xdc, 0xc2, 0x7b, 0xa8, 0xde, 0x91, 0x90, 0x6e, 0x8c, 0x47, 0x51, 0xbe, 0xdb,
	0x68, 0xb5, 0xdb, 0xa8, 0x97, 0xed, 0x96, 0xff, 0xc3, 0x43, 0x7c, 0xfb, 0xfa, 0xfe, 0x28, 0xed,
	0x21, 0xc4, 0xf3, 0xf3, 0x78, 0xec, 0xef, 0xb8, 0xd9, 0x7f, 0x29, 0xfe, 0xe6, 0xa0, 0xea, 0x9b,
	0x2e, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x59, 0xd0, 0x3f, 0xfa, 0x01, 0x00, 0x00,
}
