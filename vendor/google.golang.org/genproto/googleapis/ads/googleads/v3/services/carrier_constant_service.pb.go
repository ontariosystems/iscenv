// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/services/carrier_constant_service.proto

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

// Request message for [CarrierConstantService.GetCarrierConstant][google.ads.googleads.v3.services.CarrierConstantService.GetCarrierConstant].
type GetCarrierConstantRequest struct {
	// Required. Resource name of the carrier constant to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCarrierConstantRequest) Reset()         { *m = GetCarrierConstantRequest{} }
func (m *GetCarrierConstantRequest) String() string { return proto.CompactTextString(m) }
func (*GetCarrierConstantRequest) ProtoMessage()    {}
func (*GetCarrierConstantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3024993599d2a450, []int{0}
}

func (m *GetCarrierConstantRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCarrierConstantRequest.Unmarshal(m, b)
}
func (m *GetCarrierConstantRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCarrierConstantRequest.Marshal(b, m, deterministic)
}
func (m *GetCarrierConstantRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCarrierConstantRequest.Merge(m, src)
}
func (m *GetCarrierConstantRequest) XXX_Size() int {
	return xxx_messageInfo_GetCarrierConstantRequest.Size(m)
}
func (m *GetCarrierConstantRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCarrierConstantRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetCarrierConstantRequest proto.InternalMessageInfo

func (m *GetCarrierConstantRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetCarrierConstantRequest)(nil), "google.ads.googleads.v3.services.GetCarrierConstantRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/services/carrier_constant_service.proto", fileDescriptor_3024993599d2a450)
}

var fileDescriptor_3024993599d2a450 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xcf, 0x8a, 0xd3, 0x40,
	0x18, 0x27, 0x29, 0x08, 0x06, 0xbd, 0xe4, 0xa0, 0x35, 0x15, 0x2c, 0xa5, 0x48, 0xf1, 0x30, 0x03,
	0xcd, 0x45, 0xa6, 0xa8, 0x4c, 0x8b, 0xd4, 0x93, 0x94, 0x0a, 0x3d, 0x48, 0x20, 0x4c, 0x93, 0x31,
	0x0e, 0x24, 0x33, 0x75, 0x26, 0xcd, 0x45, 0xbc, 0xf4, 0x15, 0x7c, 0x03, 0x8f, 0xbe, 0x87, 0x97,
	0x5e, 0xbd, 0xed, 0x69, 0x0f, 0x7b, 0xda, 0x77, 0x58, 0x58, 0xd2, 0xc9, 0xa4, 0x7f, 0xb6, 0xa1,
	0xb7, 0x1f, 0xf9, 0xfd, 0xfb, 0xbe, 0x6f, 0xe2, 0x7c, 0x48, 0x84, 0x48, 0x52, 0x0a, 0x49, 0xac,
	0xa0, 0x86, 0x25, 0x2a, 0x7c, 0xa8, 0xa8, 0x2c, 0x58, 0x44, 0x15, 0x8c, 0x88, 0x94, 0x8c, 0xca,
	0x30, 0x12, 0x5c, 0xe5, 0x84, 0xe7, 0x61, 0xc5, 0x80, 0x95, 0x14, 0xb9, 0x70, 0xbb, 0xda, 0x05,
	0x48, 0xac, 0x40, 0x1d, 0x00, 0x0a, 0x1f, 0x98, 0x00, 0xef, 0x6d, 0x53, 0x85, 0xa4, 0x4a, 0xac,
	0xe5, 0xb9, 0x0e, 0x9d, 0xed, 0xbd, 0x34, 0xce, 0x15, 0x83, 0x84, 0x73, 0x91, 0x93, 0x9c, 0x09,
	0xae, 0x2a, 0xf6, 0xf9, 0x01, 0x1b, 0xa5, 0x8c, 0xd6, 0xb6, 0x57, 0x07, 0xc4, 0x37, 0x46, 0xd3,
	0x38, 0x5c, 0xd2, 0xef, 0xa4, 0x60, 0x42, 0x6a, 0x41, 0xef, 0xa3, 0xf3, 0x62, 0x4a, 0xf3, 0x89,
	0x2e, 0x9d, 0x54, 0x9d, 0x73, 0xfa, 0x63, 0x4d, 0x55, 0xee, 0x0e, 0x9c, 0xa7, 0x66, 0xb0, 0x90,
	0x93, 0x8c, 0xb6, 0xad, 0xae, 0x35, 0x78, 0x3c, 0x6e, 0x5d, 0x63, 0x7b, 0xfe, 0xc4, 0x30, 0x9f,
	0x49, 0x46, 0x87, 0x77, 0x96, 0xf3, 0xec, 0x24, 0xe4, 0x8b, 0x5e, 0xda, 0xfd, 0x67, 0x39, 0xee,
	0xc3, 0x0a, 0x77, 0x04, 0x2e, 0x5d, 0x0b, 0x34, 0x0e, 0xe6, 0x0d, 0x1b, 0xcd, 0xf5, 0x21, 0xc1,
	0x89, 0xb5, 0xf7, 0xfe, 0x0a, 0x1f, 0x6f, 0xb3, 0xf9, 0x7f, 0xf3, 0xdb, 0x1e, 0xb8, 0xaf, 0xcb,
	0xfb, 0xff, 0x3c, 0x62, 0xde, 0x45, 0xc7, 0x5e, 0x05, 0xdf, 0xfc, 0xf2, 0x3a, 0x5b, 0xdc, 0xde,
	0x57, 0x55, 0x68, 0xc5, 0x14, 0x88, 0x44, 0x36, 0xde, 0xd8, 0x4e, 0x3f, 0x12, 0xd9, 0xc5, 0x9d,
	0xc6, 0x9d, 0xf3, 0x57, 0x9a, 0x95, 0x8f, 0x31, 0xb3, 0xbe, 0x7e, 0xaa, 0x02, 0x12, 0x91, 0x12,
	0x9e, 0x00, 0x21, 0x13, 0x98, 0x50, 0xbe, 0x7b, 0x2a, 0xb8, 0xaf, 0x6c, 0xfe, 0x45, 0x47, 0x06,
	0xfc, 0xb1, 0x5b, 0x53, 0x8c, 0xff, 0xda, 0xdd, 0xa9, 0x0e, 0xc4, 0xb1, 0x02, 0x1a, 0x96, 0x68,
	0xe1, 0x83, 0xaa, 0x58, 0x6d, 0x8d, 0x24, 0xc0, 0xb1, 0x0a, 0x6a, 0x49, 0xb0, 0xf0, 0x03, 0x23,
	0xb9, 0xb5, 0xfb, 0xfa, 0x3b, 0x42, 0x38, 0x56, 0x08, 0xd5, 0x22, 0x84, 0x16, 0x3e, 0x42, 0x46,
	0xb6, 0x7c, 0xb4, 0x9b, 0xd3, 0xbf, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x69, 0x87, 0x81, 0x13, 0x49,
	0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CarrierConstantServiceClient is the client API for CarrierConstantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CarrierConstantServiceClient interface {
	// Returns the requested carrier constant in full detail.
	GetCarrierConstant(ctx context.Context, in *GetCarrierConstantRequest, opts ...grpc.CallOption) (*resources.CarrierConstant, error)
}

type carrierConstantServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCarrierConstantServiceClient(cc grpc.ClientConnInterface) CarrierConstantServiceClient {
	return &carrierConstantServiceClient{cc}
}

func (c *carrierConstantServiceClient) GetCarrierConstant(ctx context.Context, in *GetCarrierConstantRequest, opts ...grpc.CallOption) (*resources.CarrierConstant, error) {
	out := new(resources.CarrierConstant)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v3.services.CarrierConstantService/GetCarrierConstant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CarrierConstantServiceServer is the server API for CarrierConstantService service.
type CarrierConstantServiceServer interface {
	// Returns the requested carrier constant in full detail.
	GetCarrierConstant(context.Context, *GetCarrierConstantRequest) (*resources.CarrierConstant, error)
}

// UnimplementedCarrierConstantServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCarrierConstantServiceServer struct {
}

func (*UnimplementedCarrierConstantServiceServer) GetCarrierConstant(ctx context.Context, req *GetCarrierConstantRequest) (*resources.CarrierConstant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCarrierConstant not implemented")
}

func RegisterCarrierConstantServiceServer(s *grpc.Server, srv CarrierConstantServiceServer) {
	s.RegisterService(&_CarrierConstantService_serviceDesc, srv)
}

func _CarrierConstantService_GetCarrierConstant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCarrierConstantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarrierConstantServiceServer).GetCarrierConstant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v3.services.CarrierConstantService/GetCarrierConstant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarrierConstantServiceServer).GetCarrierConstant(ctx, req.(*GetCarrierConstantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CarrierConstantService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v3.services.CarrierConstantService",
	HandlerType: (*CarrierConstantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCarrierConstant",
			Handler:    _CarrierConstantService_GetCarrierConstant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v3/services/carrier_constant_service.proto",
}
