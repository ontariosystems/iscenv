// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/services/location_view_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v3/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Request message for [LocationViewService.GetLocationView][google.ads.googleads.v3.services.LocationViewService.GetLocationView].
type GetLocationViewRequest struct {
	// Required. The resource name of the location view to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetLocationViewRequest) Reset()         { *m = GetLocationViewRequest{} }
func (m *GetLocationViewRequest) String() string { return proto.CompactTextString(m) }
func (*GetLocationViewRequest) ProtoMessage()    {}
func (*GetLocationViewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd56f84abb1064dc, []int{0}
}

func (m *GetLocationViewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLocationViewRequest.Unmarshal(m, b)
}
func (m *GetLocationViewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLocationViewRequest.Marshal(b, m, deterministic)
}
func (m *GetLocationViewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLocationViewRequest.Merge(m, src)
}
func (m *GetLocationViewRequest) XXX_Size() int {
	return xxx_messageInfo_GetLocationViewRequest.Size(m)
}
func (m *GetLocationViewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLocationViewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetLocationViewRequest proto.InternalMessageInfo

func (m *GetLocationViewRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetLocationViewRequest)(nil), "google.ads.googleads.v3.services.GetLocationViewRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/services/location_view_service.proto", fileDescriptor_dd56f84abb1064dc)
}

var fileDescriptor_dd56f84abb1064dc = []byte{
	// 411 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x3f, 0xcb, 0xd3, 0x40,
	0x18, 0x27, 0x29, 0x08, 0x06, 0x45, 0x88, 0xa0, 0x25, 0x0a, 0x96, 0xd2, 0xa1, 0x74, 0xb8, 0x43,
	0x83, 0x20, 0xa7, 0x0e, 0x97, 0x25, 0x0e, 0x22, 0xa5, 0x42, 0x06, 0x09, 0x84, 0x6b, 0xf2, 0x18,
	0x0f, 0x92, 0x5c, 0xcd, 0xa5, 0xe9, 0x20, 0x2e, 0x7e, 0x05, 0xbf, 0x81, 0xa3, 0xdf, 0xc3, 0xa5,
	0xab, 0x9b, 0x93, 0x83, 0x93, 0x1f, 0xe1, 0x5d, 0xde, 0x97, 0xf4, 0x72, 0x69, 0xfa, 0xbe, 0x2d,
	0xdd, 0x7e, 0xe4, 0xf7, 0xe7, 0x79, 0x9e, 0xdf, 0xc5, 0x7a, 0x95, 0x0a, 0x91, 0x66, 0x80, 0x59,
	0x22, 0xb1, 0x82, 0x0d, 0xaa, 0x5d, 0x2c, 0xa1, 0xac, 0x79, 0x0c, 0x12, 0x67, 0x22, 0x66, 0x15,
	0x17, 0x45, 0x54, 0x73, 0xd8, 0x44, 0xed, 0x67, 0xb4, 0x2a, 0x45, 0x25, 0xec, 0x91, 0xb2, 0x20,
	0x96, 0x48, 0xd4, 0xb9, 0x51, 0xed, 0x22, 0xed, 0x76, 0x9e, 0x9f, 0xca, 0x2f, 0x41, 0x8a, 0x75,
	0x79, 0x63, 0x80, 0x0a, 0x76, 0x1e, 0x6b, 0xdb, 0x8a, 0x63, 0x56, 0x14, 0xa2, 0xda, 0x29, 0x64,
	0xcb, 0x3e, 0xec, 0xb1, 0x71, 0xc6, 0xa1, 0xa8, 0x5a, 0xe2, 0x49, 0x8f, 0xf8, 0xc8, 0x21, 0x4b,
	0xa2, 0x25, 0x7c, 0x62, 0x35, 0x17, 0xa5, 0x12, 0x8c, 0x3d, 0xeb, 0x81, 0x0f, 0xd5, 0xdb, 0x76,
	0x62, 0xc0, 0x61, 0xb3, 0x80, 0xcf, 0x6b, 0x90, 0x95, 0x3d, 0xb5, 0xee, 0xea, 0x95, 0xa2, 0x82,
	0xe5, 0x30, 0x34, 0x46, 0xc6, 0xf4, 0xb6, 0x37, 0xf8, 0x4b, 0xcd, 0xc5, 0x1d, 0xcd, 0xbc, 0x63,
	0x39, 0x3c, 0xbb, 0x30, 0xac, 0xfb, 0xfd, 0x84, 0xf7, 0xea, 0x56, 0xfb, 0x97, 0x61, 0xdd, 0xbb,
	0x16, 0x6e, 0xbf, 0x40, 0xe7, 0x1a, 0x42, 0xc7, 0xf7, 0x71, 0xf0, 0x49, 0x67, 0xd7, 0x1c, 0xea,
	0xfb, 0xc6, 0xfe, 0x1f, 0x7a, 0x78, 0xc1, 0xb7, 0xdf, 0xff, 0xbe, 0x9b, 0x4f, 0x6d, 0xdc, 0xb4,
	0xfd, 0xe5, 0x80, 0x79, 0x1d, 0xaf, 0x65, 0x25, 0x72, 0x28, 0x25, 0x9e, 0x75, 0xf5, 0x37, 0x21,
	0x12, 0xcf, 0xbe, 0x3a, 0x8f, 0xb6, 0x74, 0xb8, 0x1f, 0xd8, 0xa2, 0x15, 0x97, 0x28, 0x16, 0xb9,
	0x77, 0x69, 0x58, 0x93, 0x58, 0xe4, 0x67, 0xcf, 0xf2, 0x86, 0x47, 0x2a, 0x9a, 0x37, 0x6f, 0x30,
	0x37, 0x3e, 0xbc, 0x69, 0xdd, 0xa9, 0xc8, 0x58, 0x91, 0x22, 0x51, 0xa6, 0x38, 0x85, 0x62, 0xf7,
	0x42, 0x78, 0x3f, 0xef, 0xf4, 0x3f, 0xf9, 0x52, 0x83, 0x1f, 0xe6, 0xc0, 0xa7, 0xf4, 0xa7, 0x39,
	0xf2, 0x55, 0x20, 0x4d, 0x24, 0x52, 0xb0, 0x41, 0x81, 0x8b, 0xda, 0xc1, 0x72, 0xab, 0x25, 0x21,
	0x4d, 0x64, 0xd8, 0x49, 0xc2, 0xc0, 0x0d, 0xb5, 0xe4, 0xbf, 0x39, 0x51, 0xdf, 0x09, 0xa1, 0x89,
	0x24, 0xa4, 0x13, 0x11, 0x12, 0xb8, 0x84, 0x68, 0xd9, 0xf2, 0xd6, 0x6e, 0x4f, 0xf7, 0x2a, 0x00,
	0x00, 0xff, 0xff, 0xfa, 0xae, 0xea, 0x6c, 0x3a, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LocationViewServiceClient is the client API for LocationViewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LocationViewServiceClient interface {
	// Returns the requested location view in full detail.
	GetLocationView(ctx context.Context, in *GetLocationViewRequest, opts ...grpc.CallOption) (*resources.LocationView, error)
}

type locationViewServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationViewServiceClient(cc grpc.ClientConnInterface) LocationViewServiceClient {
	return &locationViewServiceClient{cc}
}

func (c *locationViewServiceClient) GetLocationView(ctx context.Context, in *GetLocationViewRequest, opts ...grpc.CallOption) (*resources.LocationView, error) {
	out := new(resources.LocationView)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v3.services.LocationViewService/GetLocationView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationViewServiceServer is the server API for LocationViewService service.
type LocationViewServiceServer interface {
	// Returns the requested location view in full detail.
	GetLocationView(context.Context, *GetLocationViewRequest) (*resources.LocationView, error)
}

// UnimplementedLocationViewServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLocationViewServiceServer struct {
}

func (*UnimplementedLocationViewServiceServer) GetLocationView(ctx context.Context, req *GetLocationViewRequest) (*resources.LocationView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocationView not implemented")
}

func RegisterLocationViewServiceServer(s *grpc.Server, srv LocationViewServiceServer) {
	s.RegisterService(&_LocationViewService_serviceDesc, srv)
}

func _LocationViewService_GetLocationView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLocationViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationViewServiceServer).GetLocationView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v3.services.LocationViewService/GetLocationView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationViewServiceServer).GetLocationView(ctx, req.(*GetLocationViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LocationViewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v3.services.LocationViewService",
	HandlerType: (*LocationViewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLocationView",
			Handler:    _LocationViewService_GetLocationView_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v3/services/location_view_service.proto",
}
