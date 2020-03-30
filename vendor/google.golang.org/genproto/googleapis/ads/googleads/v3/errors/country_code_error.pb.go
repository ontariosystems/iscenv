// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/errors/country_code_error.proto

package errors

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

// Enum describing country code errors.
type CountryCodeErrorEnum_CountryCodeError int32

const (
	// Enum unspecified.
	CountryCodeErrorEnum_UNSPECIFIED CountryCodeErrorEnum_CountryCodeError = 0
	// The received error code is not known in this version.
	CountryCodeErrorEnum_UNKNOWN CountryCodeErrorEnum_CountryCodeError = 1
	// The country code is invalid.
	CountryCodeErrorEnum_INVALID_COUNTRY_CODE CountryCodeErrorEnum_CountryCodeError = 2
)

var CountryCodeErrorEnum_CountryCodeError_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "INVALID_COUNTRY_CODE",
}

var CountryCodeErrorEnum_CountryCodeError_value = map[string]int32{
	"UNSPECIFIED":          0,
	"UNKNOWN":              1,
	"INVALID_COUNTRY_CODE": 2,
}

func (x CountryCodeErrorEnum_CountryCodeError) String() string {
	return proto.EnumName(CountryCodeErrorEnum_CountryCodeError_name, int32(x))
}

func (CountryCodeErrorEnum_CountryCodeError) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c64e23a252258a7c, []int{0, 0}
}

// Container for enum describing country code errors.
type CountryCodeErrorEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountryCodeErrorEnum) Reset()         { *m = CountryCodeErrorEnum{} }
func (m *CountryCodeErrorEnum) String() string { return proto.CompactTextString(m) }
func (*CountryCodeErrorEnum) ProtoMessage()    {}
func (*CountryCodeErrorEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_c64e23a252258a7c, []int{0}
}

func (m *CountryCodeErrorEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CountryCodeErrorEnum.Unmarshal(m, b)
}
func (m *CountryCodeErrorEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CountryCodeErrorEnum.Marshal(b, m, deterministic)
}
func (m *CountryCodeErrorEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountryCodeErrorEnum.Merge(m, src)
}
func (m *CountryCodeErrorEnum) XXX_Size() int {
	return xxx_messageInfo_CountryCodeErrorEnum.Size(m)
}
func (m *CountryCodeErrorEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_CountryCodeErrorEnum.DiscardUnknown(m)
}

var xxx_messageInfo_CountryCodeErrorEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v3.errors.CountryCodeErrorEnum_CountryCodeError", CountryCodeErrorEnum_CountryCodeError_name, CountryCodeErrorEnum_CountryCodeError_value)
	proto.RegisterType((*CountryCodeErrorEnum)(nil), "google.ads.googleads.v3.errors.CountryCodeErrorEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/errors/country_code_error.proto", fileDescriptor_c64e23a252258a7c)
}

var fileDescriptor_c64e23a252258a7c = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xd1, 0x4a, 0xf3, 0x30,
	0x1c, 0xc5, 0xbf, 0xf5, 0x03, 0x85, 0xec, 0xc2, 0x52, 0x26, 0x88, 0xc8, 0x2e, 0xfa, 0x00, 0xc9,
	0x45, 0x2e, 0x84, 0x78, 0x95, 0xb5, 0x75, 0x54, 0x25, 0x1b, 0xea, 0x2a, 0x4a, 0xa1, 0x74, 0x4b,
	0x09, 0x83, 0x2d, 0xff, 0x91, 0x74, 0x03, 0x5f, 0xc7, 0x4b, 0x1f, 0xc5, 0x47, 0xf1, 0xca, 0x47,
	0x90, 0x36, 0xae, 0x17, 0x03, 0xbd, 0xca, 0xe1, 0xf0, 0x3b, 0x27, 0x87, 0x3f, 0xba, 0x54, 0x00,
	0x6a, 0x55, 0x91, 0x52, 0x5a, 0xe2, 0x64, 0xa3, 0x76, 0x94, 0x54, 0xc6, 0x80, 0xb1, 0x64, 0x01,
	0x5b, 0x5d, 0x9b, 0xd7, 0x62, 0x01, 0xb2, 0x2a, 0x5a, 0x0f, 0x6f, 0x0c, 0xd4, 0x10, 0x0c, 0x1d,
	0x8d, 0x4b, 0x69, 0x71, 0x17, 0xc4, 0x3b, 0x8a, 0x5d, 0xf0, 0xfc, 0x62, 0x5f, 0xbc, 0x59, 0x92,
	0x52, 0x6b, 0xa8, 0xcb, 0x7a, 0x09, 0xda, 0xba, 0x74, 0x38, 0x47, 0x83, 0xc8, 0x35, 0x47, 0x20,
	0xab, 0xa4, 0x89, 0x24, 0x7a, 0xbb, 0x0e, 0x6f, 0x90, 0x7f, 0xe8, 0x07, 0x27, 0xa8, 0x3f, 0x13,
	0x0f, 0xd3, 0x24, 0x4a, 0xaf, 0xd3, 0x24, 0xf6, 0xff, 0x05, 0x7d, 0x74, 0x3c, 0x13, 0xb7, 0x62,
	0xf2, 0x24, 0xfc, 0x5e, 0x70, 0x86, 0x06, 0xa9, 0xc8, 0xf8, 0x5d, 0x1a, 0x17, 0xd1, 0x64, 0x26,
	0x1e, 0xef, 0x9f, 0x8b, 0x68, 0x12, 0x27, 0xbe, 0x37, 0xfa, 0xea, 0xa1, 0x70, 0x01, 0x6b, 0xfc,
	0xf7, 0xd0, 0xd1, 0xe9, 0xe1, 0x87, 0xd3, 0x66, 0xe1, 0xb4, 0xf7, 0x12, 0xff, 0x04, 0x15, 0xac,
	0x4a, 0xad, 0x30, 0x18, 0x45, 0x54, 0xa5, 0xdb, 0xfd, 0xfb, 0x53, 0x6d, 0x96, 0xf6, 0xb7, 0xcb,
	0x5d, 0xb9, 0xe7, 0xcd, 0xfb, 0x3f, 0xe6, 0xfc, 0xdd, 0x1b, 0x8e, 0x5d, 0x19, 0x97, 0x16, 0x3b,
	0xd9, 0xa8, 0x8c, 0xe2, 0xf6, 0x4b, 0xfb, 0xb1, 0x07, 0x72, 0x2e, 0x6d, 0xde, 0x01, 0x79, 0x46,
	0x73, 0x07, 0x7c, 0x7a, 0xa1, 0x73, 0x19, 0xe3, 0xd2, 0x32, 0xd6, 0x21, 0x8c, 0x65, 0x94, 0x31,
	0x07, 0xcd, 0x8f, 0xda, 0x75, 0xf4, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x0e, 0x92, 0xbb, 0xd6,
	0x01, 0x00, 0x00,
}
