// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/resources/label.proto

package resources

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	common "google.golang.org/genproto/googleapis/ads/googleads/v1/common"
	enums "google.golang.org/genproto/googleapis/ads/googleads/v1/enums"
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

// A label.
type Label struct {
	// Immutable. Name of the resource.
	// Label resource names have the form:
	// `customers/{customer_id}/labels/{label_id}`
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	// Output only. Id of the label. Read only.
	Id *wrappers.Int64Value `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the label.
	//
	// This field is required and should not be empty when creating a new label.
	//
	// The length of this string should be between 1 and 80, inclusive.
	Name *wrappers.StringValue `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// Output only. Status of the label. Read only.
	Status enums.LabelStatusEnum_LabelStatus `protobuf:"varint,4,opt,name=status,proto3,enum=google.ads.googleads.v1.enums.LabelStatusEnum_LabelStatus" json:"status,omitempty"`
	// A type of label displaying text on a colored background.
	TextLabel            *common.TextLabel `protobuf:"bytes,5,opt,name=text_label,json=textLabel,proto3" json:"text_label,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Label) Reset()         { *m = Label{} }
func (m *Label) String() string { return proto.CompactTextString(m) }
func (*Label) ProtoMessage()    {}
func (*Label) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4f50cd20a4c405a, []int{0}
}

func (m *Label) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Label.Unmarshal(m, b)
}
func (m *Label) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Label.Marshal(b, m, deterministic)
}
func (m *Label) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Label.Merge(m, src)
}
func (m *Label) XXX_Size() int {
	return xxx_messageInfo_Label.Size(m)
}
func (m *Label) XXX_DiscardUnknown() {
	xxx_messageInfo_Label.DiscardUnknown(m)
}

var xxx_messageInfo_Label proto.InternalMessageInfo

func (m *Label) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func (m *Label) GetId() *wrappers.Int64Value {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Label) GetName() *wrappers.StringValue {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Label) GetStatus() enums.LabelStatusEnum_LabelStatus {
	if m != nil {
		return m.Status
	}
	return enums.LabelStatusEnum_UNSPECIFIED
}

func (m *Label) GetTextLabel() *common.TextLabel {
	if m != nil {
		return m.TextLabel
	}
	return nil
}

func init() {
	proto.RegisterType((*Label)(nil), "google.ads.googleads.v1.resources.Label")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/resources/label.proto", fileDescriptor_b4f50cd20a4c405a)
}

var fileDescriptor_b4f50cd20a4c405a = []byte{
	// 490 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x4f, 0x6b, 0xdb, 0x30,
	0x1c, 0xc5, 0x76, 0x53, 0xa8, 0xf6, 0xe7, 0xe0, 0x53, 0xd6, 0x95, 0x2e, 0xdd, 0xe8, 0xc8, 0x0e,
	0x93, 0xe2, 0x6e, 0xec, 0xa0, 0x9d, 0x1c, 0x28, 0xed, 0xfe, 0x30, 0x4a, 0x3a, 0x32, 0x18, 0x81,
	0xa0, 0xd8, 0xaa, 0x67, 0xb0, 0x24, 0x23, 0xc9, 0x59, 0xa1, 0xf4, 0xcb, 0x0c, 0x76, 0xd9, 0x47,
	0xd9, 0xa7, 0xe8, 0xb9, 0x5f, 0x60, 0xb0, 0xd3, 0x88, 0xfe, 0x38, 0x81, 0x91, 0xe6, 0xe4, 0x27,
	0x7e, 0xef, 0x3d, 0x3d, 0x9e, 0x7e, 0x06, 0x2f, 0x0b, 0x21, 0x8a, 0x8a, 0x22, 0x92, 0x2b, 0x64,
	0xe1, 0x02, 0xcd, 0x13, 0x24, 0xa9, 0x12, 0x8d, 0xcc, 0xa8, 0x42, 0x15, 0x99, 0xd1, 0x0a, 0xd6,
	0x52, 0x68, 0x11, 0x1f, 0x58, 0x0e, 0x24, 0xb9, 0x82, 0x2d, 0x1d, 0xce, 0x13, 0xd8, 0xd2, 0x77,
	0xd1, 0x3a, 0xc7, 0x4c, 0x30, 0x26, 0x38, 0xd2, 0xf4, 0x52, 0x4f, 0x57, 0x3c, 0x77, 0x07, 0xeb,
	0x04, 0x94, 0x37, 0xcc, 0x5d, 0x3f, 0x55, 0x9a, 0xe8, 0x46, 0x39, 0xc5, 0x13, 0xaf, 0xa8, 0x4b,
	0x74, 0x51, 0xd2, 0x2a, 0x9f, 0xce, 0xe8, 0x37, 0x32, 0x2f, 0x85, 0x74, 0x84, 0x47, 0x2b, 0x04,
	0x9f, 0xcc, 0x8d, 0xf6, 0xdd, 0xc8, 0x9c, 0x66, 0xcd, 0x05, 0xfa, 0x2e, 0x49, 0x5d, 0x53, 0xe9,
	0xbd, 0xf7, 0x56, 0xa4, 0x84, 0x73, 0xa1, 0x89, 0x2e, 0x05, 0x77, 0xd3, 0xa7, 0x3f, 0x23, 0xd0,
	0xf9, 0xb8, 0x08, 0x14, 0x7f, 0x00, 0x0f, 0xbc, 0xf3, 0x94, 0x13, 0x46, 0xbb, 0x41, 0x2f, 0xe8,
	0xef, 0x0c, 0x9f, 0xdf, 0xa4, 0x9d, 0xbf, 0x69, 0x0f, 0xec, 0x2f, 0xdb, 0x71, 0xa8, 0x2e, 0x15,
	0xcc, 0x04, 0x43, 0x46, 0x3e, 0xba, 0xef, 0xc5, 0x9f, 0x08, 0xa3, 0xf1, 0x00, 0x84, 0x65, 0xde,
	0x0d, 0x7b, 0x41, 0xff, 0xde, 0xd1, 0x63, 0x27, 0x80, 0x3e, 0x21, 0x7c, 0xc7, 0xf5, 0x9b, 0xd7,
	0x63, 0x52, 0x35, 0x74, 0x18, 0xdd, 0xa4, 0xd1, 0x28, 0x2c, 0xf3, 0x78, 0x00, 0xb6, 0xcc, 0xad,
	0x91, 0xd1, 0xec, 0xfd, 0xa7, 0x39, 0xd7, 0xb2, 0xe4, 0x85, 0x11, 0x8d, 0x0c, 0x33, 0xfe, 0x02,
	0xb6, 0x6d, 0x89, 0xdd, 0xad, 0x5e, 0xd0, 0x7f, 0x78, 0x84, 0xe1, 0xba, 0xb7, 0x34, 0xbd, 0x43,
	0x93, 0xf3, 0xdc, 0x28, 0x8e, 0x79, 0xc3, 0x56, 0xcf, 0x36, 0x86, 0xb3, 0x8b, 0x4f, 0x01, 0x58,
	0xbe, 0x69, 0xb7, 0x63, 0x02, 0xbd, 0x58, 0x6b, 0x6e, 0xb7, 0x00, 0x7e, 0xa6, 0x97, 0xda, 0x36,
	0xb1, 0xa3, 0x3d, 0xc4, 0xa7, 0xb7, 0xe9, 0xf1, 0xa6, 0xe6, 0xe2, 0x67, 0x59, 0xa3, 0xb4, 0x60,
	0x54, 0x2a, 0x74, 0xe5, 0xe1, 0xb5, 0xdd, 0x12, 0x85, 0xae, 0xcc, 0xf7, 0x7a, 0xf8, 0x27, 0x00,
	0x87, 0x99, 0x60, 0x70, 0xe3, 0xba, 0x0e, 0x81, 0x71, 0x3d, 0x5b, 0xf4, 0x76, 0x16, 0x7c, 0x7d,
	0xef, 0x04, 0x85, 0xa8, 0x08, 0x2f, 0xa0, 0x90, 0x05, 0x2a, 0x28, 0x37, 0xad, 0xa2, 0x65, 0x9a,
	0x3b, 0xfe, 0x96, 0xb7, 0x2d, 0xfa, 0x11, 0x46, 0x27, 0x69, 0xfa, 0x2b, 0x3c, 0x38, 0xb1, 0x96,
	0x69, 0xae, 0xa0, 0x85, 0x0b, 0x34, 0x4e, 0xe0, 0xc8, 0x33, 0x7f, 0x7b, 0xce, 0x24, 0xcd, 0xd5,
	0xa4, 0xe5, 0x4c, 0xc6, 0xc9, 0xa4, 0xe5, 0xdc, 0x86, 0x87, 0x76, 0x80, 0x71, 0x9a, 0x2b, 0x8c,
	0x5b, 0x16, 0xc6, 0xe3, 0x04, 0xe3, 0x96, 0x37, 0xdb, 0x36, 0x61, 0x5f, 0xfd, 0x0b, 0x00, 0x00,
	0xff, 0xff, 0x03, 0x5b, 0xb5, 0x86, 0xd9, 0x03, 0x00, 0x00,
}
