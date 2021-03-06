// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/enums/spending_limit_type.proto

package enums

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The possible spending limit types used by certain resources as an
// alternative to absolute money values in micros.
type SpendingLimitTypeEnum_SpendingLimitType int32

const (
	// Not specified.
	SpendingLimitTypeEnum_UNSPECIFIED SpendingLimitTypeEnum_SpendingLimitType = 0
	// Used for return value only. Represents value unknown in this version.
	SpendingLimitTypeEnum_UNKNOWN SpendingLimitTypeEnum_SpendingLimitType = 1
	// Infinite, indicates unlimited spending power.
	SpendingLimitTypeEnum_INFINITE SpendingLimitTypeEnum_SpendingLimitType = 2
)

var SpendingLimitTypeEnum_SpendingLimitType_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "INFINITE",
}

var SpendingLimitTypeEnum_SpendingLimitType_value = map[string]int32{
	"UNSPECIFIED": 0,
	"UNKNOWN":     1,
	"INFINITE":    2,
}

func (x SpendingLimitTypeEnum_SpendingLimitType) String() string {
	return proto.EnumName(SpendingLimitTypeEnum_SpendingLimitType_name, int32(x))
}

func (SpendingLimitTypeEnum_SpendingLimitType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_05de225243b8687a, []int{0, 0}
}

// Message describing spending limit types.
type SpendingLimitTypeEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpendingLimitTypeEnum) Reset()         { *m = SpendingLimitTypeEnum{} }
func (m *SpendingLimitTypeEnum) String() string { return proto.CompactTextString(m) }
func (*SpendingLimitTypeEnum) ProtoMessage()    {}
func (*SpendingLimitTypeEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_05de225243b8687a, []int{0}
}

func (m *SpendingLimitTypeEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpendingLimitTypeEnum.Unmarshal(m, b)
}
func (m *SpendingLimitTypeEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpendingLimitTypeEnum.Marshal(b, m, deterministic)
}
func (m *SpendingLimitTypeEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpendingLimitTypeEnum.Merge(m, src)
}
func (m *SpendingLimitTypeEnum) XXX_Size() int {
	return xxx_messageInfo_SpendingLimitTypeEnum.Size(m)
}
func (m *SpendingLimitTypeEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_SpendingLimitTypeEnum.DiscardUnknown(m)
}

var xxx_messageInfo_SpendingLimitTypeEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v3.enums.SpendingLimitTypeEnum_SpendingLimitType", SpendingLimitTypeEnum_SpendingLimitType_name, SpendingLimitTypeEnum_SpendingLimitType_value)
	proto.RegisterType((*SpendingLimitTypeEnum)(nil), "google.ads.googleads.v3.enums.SpendingLimitTypeEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/enums/spending_limit_type.proto", fileDescriptor_05de225243b8687a)
}

var fileDescriptor_05de225243b8687a = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x50, 0xcf, 0x4a, 0xfb, 0x30,
	0x00, 0xfe, 0xad, 0x3f, 0x50, 0xc9, 0x04, 0x67, 0x41, 0x0f, 0xe2, 0x0e, 0xdb, 0x03, 0x24, 0x87,
	0x1c, 0x84, 0x78, 0x90, 0x4e, 0xbb, 0x51, 0x94, 0x38, 0xd8, 0x56, 0x45, 0x0a, 0xa3, 0x9a, 0x10,
	0x02, 0x6d, 0x12, 0x96, 0x6e, 0xb0, 0xd7, 0xf1, 0xe8, 0xa3, 0xf8, 0x28, 0xfa, 0x12, 0xd2, 0x64,
	0xed, 0x65, 0xe8, 0x25, 0x7c, 0xe4, 0xfb, 0x93, 0x2f, 0x1f, 0xb8, 0x12, 0x5a, 0x8b, 0x82, 0xa3,
	0x9c, 0x59, 0xe4, 0x61, 0x8d, 0x36, 0x18, 0x71, 0xb5, 0x2e, 0x2d, 0xb2, 0x86, 0x2b, 0x26, 0x95,
	0x58, 0x16, 0xb2, 0x94, 0xd5, 0xb2, 0xda, 0x1a, 0x0e, 0xcd, 0x4a, 0x57, 0x3a, 0xec, 0x7b, 0x35,
	0xcc, 0x99, 0x85, 0xad, 0x11, 0x6e, 0x30, 0x74, 0xc6, 0x8b, 0xcb, 0x26, 0xd7, 0x48, 0x94, 0x2b,
	0xa5, 0xab, 0xbc, 0x92, 0x5a, 0x59, 0x6f, 0x1e, 0x3e, 0x83, 0xb3, 0xd9, 0x2e, 0xf9, 0xa1, 0x0e,
	0x9e, 0x6f, 0x0d, 0x8f, 0xd5, 0xba, 0x1c, 0xde, 0x80, 0xd3, 0x3d, 0x22, 0x3c, 0x01, 0xdd, 0x05,
	0x9d, 0x4d, 0xe3, 0xdb, 0x64, 0x9c, 0xc4, 0x77, 0xbd, 0x7f, 0x61, 0x17, 0x1c, 0x2e, 0xe8, 0x3d,
	0x7d, 0x7c, 0xa2, 0xbd, 0x4e, 0x78, 0x0c, 0x8e, 0x12, 0x3a, 0x4e, 0x68, 0x32, 0x8f, 0x7b, 0xc1,
	0xe8, 0xbb, 0x03, 0x06, 0x6f, 0xba, 0x84, 0x7f, 0xb6, 0x1b, 0x9d, 0xef, 0x3d, 0x32, 0xad, 0x7b,
	0x4d, 0x3b, 0x2f, 0xa3, 0x9d, 0x51, 0xe8, 0x22, 0x57, 0x02, 0xea, 0x95, 0x40, 0x82, 0x2b, 0xd7,
	0xba, 0xd9, 0xc7, 0x48, 0xfb, 0xcb, 0x5c, 0xd7, 0xee, 0x7c, 0x0f, 0xfe, 0x4f, 0xa2, 0xe8, 0x23,
	0xe8, 0x4f, 0x7c, 0x54, 0xc4, 0x2c, 0xf4, 0xb0, 0x46, 0x29, 0x86, 0xf5, 0x4f, 0xed, 0x67, 0xc3,
	0x67, 0x11, 0xb3, 0x59, 0xcb, 0x67, 0x29, 0xce, 0x1c, 0xff, 0x15, 0x0c, 0xfc, 0x25, 0x21, 0x11,
	0xb3, 0x84, 0xb4, 0x0a, 0x42, 0x52, 0x4c, 0x88, 0xd3, 0xbc, 0x1e, 0xb8, 0x62, 0xf8, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x78, 0x15, 0xb3, 0xa4, 0xc6, 0x01, 0x00, 0x00,
}
