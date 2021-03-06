// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/services/product_group_view_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v1/resources"
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

// Request message for [ProductGroupViewService.GetProductGroupView][google.ads.googleads.v1.services.ProductGroupViewService.GetProductGroupView].
type GetProductGroupViewRequest struct {
	// Required. The resource name of the product group view to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProductGroupViewRequest) Reset()         { *m = GetProductGroupViewRequest{} }
func (m *GetProductGroupViewRequest) String() string { return proto.CompactTextString(m) }
func (*GetProductGroupViewRequest) ProtoMessage()    {}
func (*GetProductGroupViewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b1170300877f7364, []int{0}
}

func (m *GetProductGroupViewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetProductGroupViewRequest.Unmarshal(m, b)
}
func (m *GetProductGroupViewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetProductGroupViewRequest.Marshal(b, m, deterministic)
}
func (m *GetProductGroupViewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProductGroupViewRequest.Merge(m, src)
}
func (m *GetProductGroupViewRequest) XXX_Size() int {
	return xxx_messageInfo_GetProductGroupViewRequest.Size(m)
}
func (m *GetProductGroupViewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProductGroupViewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetProductGroupViewRequest proto.InternalMessageInfo

func (m *GetProductGroupViewRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetProductGroupViewRequest)(nil), "google.ads.googleads.v1.services.GetProductGroupViewRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/services/product_group_view_service.proto", fileDescriptor_b1170300877f7364)
}

var fileDescriptor_b1170300877f7364 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x3f, 0x8b, 0xd4, 0x40,
	0x1c, 0x25, 0x39, 0x10, 0x0c, 0xda, 0xc4, 0xe2, 0x8e, 0x78, 0xe0, 0x72, 0x5c, 0x71, 0x5c, 0x31,
	0x43, 0x5c, 0x6c, 0x46, 0x2d, 0x66, 0x0b, 0x23, 0x08, 0xb2, 0x9c, 0x90, 0x42, 0x02, 0x61, 0x2e,
	0xf9, 0x19, 0x07, 0x92, 0x4c, 0x9c, 0x99, 0xe4, 0x0a, 0xb1, 0xd1, 0xde, 0xc6, 0x6f, 0x60, 0xe9,
	0x47, 0xb9, 0x56, 0x6c, 0xac, 0x2c, 0xac, 0xfc, 0x14, 0x92, 0x9d, 0x4c, 0xf6, 0x0f, 0x1b, 0xb6,
	0x7b, 0xcc, 0x7b, 0xbf, 0xf7, 0x7e, 0x7f, 0xc6, 0xa3, 0x85, 0x10, 0x45, 0x09, 0x98, 0xe5, 0x0a,
	0x1b, 0xd8, 0xa3, 0x2e, 0xc4, 0x0a, 0x64, 0xc7, 0x33, 0x50, 0xb8, 0x91, 0x22, 0x6f, 0x33, 0x9d,
	0x16, 0x52, 0xb4, 0x4d, 0xda, 0x71, 0xb8, 0x49, 0x07, 0x0e, 0x35, 0x52, 0x68, 0xe1, 0xcf, 0x4c,
	0x1d, 0x62, 0xb9, 0x42, 0xa3, 0x05, 0xea, 0x42, 0x64, 0x2d, 0x02, 0x32, 0x15, 0x22, 0x41, 0x89,
	0x56, 0xee, 0x4f, 0x31, 0xee, 0xc1, 0xa9, 0xad, 0x6d, 0x38, 0x66, 0x75, 0x2d, 0x34, 0xd3, 0x5c,
	0xd4, 0x6a, 0x60, 0x8f, 0x37, 0xd8, 0xac, 0xe4, 0x50, 0xeb, 0x81, 0x78, 0xb4, 0x41, 0xbc, 0xe3,
	0x50, 0xe6, 0xe9, 0x35, 0xbc, 0x67, 0x1d, 0x17, 0xd2, 0x08, 0xce, 0x5e, 0x78, 0x41, 0x04, 0x7a,
	0x69, 0x62, 0xa3, 0x3e, 0x35, 0xe6, 0x70, 0x73, 0x05, 0x1f, 0x5a, 0x50, 0xda, 0xbf, 0xf0, 0xee,
	0xdb, 0xde, 0xd2, 0x9a, 0x55, 0x70, 0xe2, 0xcc, 0x9c, 0x8b, 0xbb, 0x8b, 0xa3, 0x3f, 0xd4, 0xbd,
	0xba, 0x67, 0x99, 0xd7, 0xac, 0x82, 0xc7, 0x5f, 0x5d, 0xef, 0x78, 0xd7, 0xe5, 0x8d, 0x19, 0xdc,
	0xff, 0xe5, 0x78, 0x0f, 0xf6, 0x84, 0xf8, 0xcf, 0xd0, 0xa1, 0x95, 0xa1, 0xe9, 0xde, 0x82, 0xf9,
	0x64, 0xf5, 0xb8, 0x4e, 0xb4, 0x5b, 0x7b, 0xf6, 0xea, 0x37, 0xdd, 0x9e, 0xe8, 0xf3, 0xcf, 0xbf,
	0xdf, 0xdc, 0x27, 0xfe, 0xbc, 0x3f, 0xc3, 0xc7, 0x2d, 0xe6, 0x79, 0xd6, 0x2a, 0x2d, 0x2a, 0x90,
	0x0a, 0x5f, 0xda, 0xbb, 0x8c, 0x46, 0x0a, 0x5f, 0x7e, 0x0a, 0x1e, 0xde, 0xd2, 0x93, 0x75, 0xf0,
	0x80, 0x1a, 0xae, 0x50, 0x26, 0xaa, 0xc5, 0x17, 0xd7, 0x3b, 0xcf, 0x44, 0x75, 0x70, 0xc4, 0xc5,
	0xe9, 0xc4, 0xda, 0x96, 0xfd, 0x7d, 0x96, 0xce, 0xdb, 0x97, 0x83, 0x43, 0x21, 0x4a, 0x56, 0x17,
	0x48, 0xc8, 0x02, 0x17, 0x50, 0xaf, 0xae, 0x87, 0xd7, 0x99, 0xd3, 0x3f, 0xf7, 0xa9, 0x05, 0xdf,
	0xdd, 0xa3, 0x88, 0xd2, 0x1f, 0xee, 0x2c, 0x32, 0x86, 0x34, 0x57, 0xc8, 0xc0, 0x1e, 0xc5, 0x21,
	0x1a, 0x82, 0xd5, 0xad, 0x95, 0x24, 0x34, 0x57, 0xc9, 0x28, 0x49, 0xe2, 0x30, 0xb1, 0x92, 0x7f,
	0xee, 0xb9, 0x79, 0x27, 0x84, 0xe6, 0x8a, 0x90, 0x51, 0x44, 0x48, 0x1c, 0x12, 0x62, 0x65, 0xd7,
	0x77, 0x56, 0x7d, 0xce, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0xd8, 0xc6, 0xbe, 0x10, 0x60, 0x03,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProductGroupViewServiceClient is the client API for ProductGroupViewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProductGroupViewServiceClient interface {
	// Returns the requested product group view in full detail.
	GetProductGroupView(ctx context.Context, in *GetProductGroupViewRequest, opts ...grpc.CallOption) (*resources.ProductGroupView, error)
}

type productGroupViewServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductGroupViewServiceClient(cc grpc.ClientConnInterface) ProductGroupViewServiceClient {
	return &productGroupViewServiceClient{cc}
}

func (c *productGroupViewServiceClient) GetProductGroupView(ctx context.Context, in *GetProductGroupViewRequest, opts ...grpc.CallOption) (*resources.ProductGroupView, error) {
	out := new(resources.ProductGroupView)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.ProductGroupViewService/GetProductGroupView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductGroupViewServiceServer is the server API for ProductGroupViewService service.
type ProductGroupViewServiceServer interface {
	// Returns the requested product group view in full detail.
	GetProductGroupView(context.Context, *GetProductGroupViewRequest) (*resources.ProductGroupView, error)
}

// UnimplementedProductGroupViewServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProductGroupViewServiceServer struct {
}

func (*UnimplementedProductGroupViewServiceServer) GetProductGroupView(ctx context.Context, req *GetProductGroupViewRequest) (*resources.ProductGroupView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductGroupView not implemented")
}

func RegisterProductGroupViewServiceServer(s *grpc.Server, srv ProductGroupViewServiceServer) {
	s.RegisterService(&_ProductGroupViewService_serviceDesc, srv)
}

func _ProductGroupViewService_GetProductGroupView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductGroupViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductGroupViewServiceServer).GetProductGroupView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.ProductGroupViewService/GetProductGroupView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductGroupViewServiceServer).GetProductGroupView(ctx, req.(*GetProductGroupViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProductGroupViewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v1.services.ProductGroupViewService",
	HandlerType: (*ProductGroupViewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProductGroupView",
			Handler:    _ProductGroupViewService_GetProductGroupView_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v1/services/product_group_view_service.proto",
}
