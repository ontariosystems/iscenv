// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v2/services/ad_group_extension_setting_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v2/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	status "google.golang.org/genproto/googleapis/rpc/status"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status1 "google.golang.org/grpc/status"
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

// Request message for
// [AdGroupExtensionSettingService.GetAdGroupExtensionSetting][google.ads.googleads.v2.services.AdGroupExtensionSettingService.GetAdGroupExtensionSetting].
type GetAdGroupExtensionSettingRequest struct {
	// Required. The resource name of the ad group extension setting to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdGroupExtensionSettingRequest) Reset()         { *m = GetAdGroupExtensionSettingRequest{} }
func (m *GetAdGroupExtensionSettingRequest) String() string { return proto.CompactTextString(m) }
func (*GetAdGroupExtensionSettingRequest) ProtoMessage()    {}
func (*GetAdGroupExtensionSettingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf10cd4825f3524e, []int{0}
}

func (m *GetAdGroupExtensionSettingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdGroupExtensionSettingRequest.Unmarshal(m, b)
}
func (m *GetAdGroupExtensionSettingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdGroupExtensionSettingRequest.Marshal(b, m, deterministic)
}
func (m *GetAdGroupExtensionSettingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdGroupExtensionSettingRequest.Merge(m, src)
}
func (m *GetAdGroupExtensionSettingRequest) XXX_Size() int {
	return xxx_messageInfo_GetAdGroupExtensionSettingRequest.Size(m)
}
func (m *GetAdGroupExtensionSettingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdGroupExtensionSettingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdGroupExtensionSettingRequest proto.InternalMessageInfo

func (m *GetAdGroupExtensionSettingRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

// Request message for
// [AdGroupExtensionSettingService.MutateAdGroupExtensionSettings][google.ads.googleads.v2.services.AdGroupExtensionSettingService.MutateAdGroupExtensionSettings].
type MutateAdGroupExtensionSettingsRequest struct {
	// Required. The ID of the customer whose ad group extension settings are being
	// modified.
	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	// Required. The list of operations to perform on individual ad group extension
	// settings.
	Operations []*AdGroupExtensionSettingOperation `protobuf:"bytes,2,rep,name=operations,proto3" json:"operations,omitempty"`
	// If true, successful operations will be carried out and invalid
	// operations will return errors. If false, all operations will be carried
	// out in one transaction if and only if they are all valid.
	// Default is false.
	PartialFailure bool `protobuf:"varint,3,opt,name=partial_failure,json=partialFailure,proto3" json:"partial_failure,omitempty"`
	// If true, the request is validated but not executed. Only errors are
	// returned, not results.
	ValidateOnly         bool     `protobuf:"varint,4,opt,name=validate_only,json=validateOnly,proto3" json:"validate_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateAdGroupExtensionSettingsRequest) Reset()         { *m = MutateAdGroupExtensionSettingsRequest{} }
func (m *MutateAdGroupExtensionSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*MutateAdGroupExtensionSettingsRequest) ProtoMessage()    {}
func (*MutateAdGroupExtensionSettingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf10cd4825f3524e, []int{1}
}

func (m *MutateAdGroupExtensionSettingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsRequest.Unmarshal(m, b)
}
func (m *MutateAdGroupExtensionSettingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsRequest.Marshal(b, m, deterministic)
}
func (m *MutateAdGroupExtensionSettingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateAdGroupExtensionSettingsRequest.Merge(m, src)
}
func (m *MutateAdGroupExtensionSettingsRequest) XXX_Size() int {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsRequest.Size(m)
}
func (m *MutateAdGroupExtensionSettingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateAdGroupExtensionSettingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MutateAdGroupExtensionSettingsRequest proto.InternalMessageInfo

func (m *MutateAdGroupExtensionSettingsRequest) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *MutateAdGroupExtensionSettingsRequest) GetOperations() []*AdGroupExtensionSettingOperation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *MutateAdGroupExtensionSettingsRequest) GetPartialFailure() bool {
	if m != nil {
		return m.PartialFailure
	}
	return false
}

func (m *MutateAdGroupExtensionSettingsRequest) GetValidateOnly() bool {
	if m != nil {
		return m.ValidateOnly
	}
	return false
}

// A single operation (create, update, remove) on an ad group extension setting.
type AdGroupExtensionSettingOperation struct {
	// FieldMask that determines which resource fields are modified in an update.
	UpdateMask *field_mask.FieldMask `protobuf:"bytes,4,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	// The mutate operation.
	//
	// Types that are valid to be assigned to Operation:
	//	*AdGroupExtensionSettingOperation_Create
	//	*AdGroupExtensionSettingOperation_Update
	//	*AdGroupExtensionSettingOperation_Remove
	Operation            isAdGroupExtensionSettingOperation_Operation `protobuf_oneof:"operation"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *AdGroupExtensionSettingOperation) Reset()         { *m = AdGroupExtensionSettingOperation{} }
func (m *AdGroupExtensionSettingOperation) String() string { return proto.CompactTextString(m) }
func (*AdGroupExtensionSettingOperation) ProtoMessage()    {}
func (*AdGroupExtensionSettingOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf10cd4825f3524e, []int{2}
}

func (m *AdGroupExtensionSettingOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdGroupExtensionSettingOperation.Unmarshal(m, b)
}
func (m *AdGroupExtensionSettingOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdGroupExtensionSettingOperation.Marshal(b, m, deterministic)
}
func (m *AdGroupExtensionSettingOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdGroupExtensionSettingOperation.Merge(m, src)
}
func (m *AdGroupExtensionSettingOperation) XXX_Size() int {
	return xxx_messageInfo_AdGroupExtensionSettingOperation.Size(m)
}
func (m *AdGroupExtensionSettingOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_AdGroupExtensionSettingOperation.DiscardUnknown(m)
}

var xxx_messageInfo_AdGroupExtensionSettingOperation proto.InternalMessageInfo

func (m *AdGroupExtensionSettingOperation) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type isAdGroupExtensionSettingOperation_Operation interface {
	isAdGroupExtensionSettingOperation_Operation()
}

type AdGroupExtensionSettingOperation_Create struct {
	Create *resources.AdGroupExtensionSetting `protobuf:"bytes,1,opt,name=create,proto3,oneof"`
}

type AdGroupExtensionSettingOperation_Update struct {
	Update *resources.AdGroupExtensionSetting `protobuf:"bytes,2,opt,name=update,proto3,oneof"`
}

type AdGroupExtensionSettingOperation_Remove struct {
	Remove string `protobuf:"bytes,3,opt,name=remove,proto3,oneof"`
}

func (*AdGroupExtensionSettingOperation_Create) isAdGroupExtensionSettingOperation_Operation() {}

func (*AdGroupExtensionSettingOperation_Update) isAdGroupExtensionSettingOperation_Operation() {}

func (*AdGroupExtensionSettingOperation_Remove) isAdGroupExtensionSettingOperation_Operation() {}

func (m *AdGroupExtensionSettingOperation) GetOperation() isAdGroupExtensionSettingOperation_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *AdGroupExtensionSettingOperation) GetCreate() *resources.AdGroupExtensionSetting {
	if x, ok := m.GetOperation().(*AdGroupExtensionSettingOperation_Create); ok {
		return x.Create
	}
	return nil
}

func (m *AdGroupExtensionSettingOperation) GetUpdate() *resources.AdGroupExtensionSetting {
	if x, ok := m.GetOperation().(*AdGroupExtensionSettingOperation_Update); ok {
		return x.Update
	}
	return nil
}

func (m *AdGroupExtensionSettingOperation) GetRemove() string {
	if x, ok := m.GetOperation().(*AdGroupExtensionSettingOperation_Remove); ok {
		return x.Remove
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AdGroupExtensionSettingOperation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AdGroupExtensionSettingOperation_Create)(nil),
		(*AdGroupExtensionSettingOperation_Update)(nil),
		(*AdGroupExtensionSettingOperation_Remove)(nil),
	}
}

// Response message for an ad group extension setting mutate.
type MutateAdGroupExtensionSettingsResponse struct {
	// Errors that pertain to operation failures in the partial failure mode.
	// Returned only when partial_failure = true and all errors occur inside the
	// operations. If any errors occur outside the operations (e.g. auth errors),
	// we return an RPC level error.
	PartialFailureError *status.Status `protobuf:"bytes,3,opt,name=partial_failure_error,json=partialFailureError,proto3" json:"partial_failure_error,omitempty"`
	// All results for the mutate.
	Results              []*MutateAdGroupExtensionSettingResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *MutateAdGroupExtensionSettingsResponse) Reset() {
	*m = MutateAdGroupExtensionSettingsResponse{}
}
func (m *MutateAdGroupExtensionSettingsResponse) String() string { return proto.CompactTextString(m) }
func (*MutateAdGroupExtensionSettingsResponse) ProtoMessage()    {}
func (*MutateAdGroupExtensionSettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf10cd4825f3524e, []int{3}
}

func (m *MutateAdGroupExtensionSettingsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsResponse.Unmarshal(m, b)
}
func (m *MutateAdGroupExtensionSettingsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsResponse.Marshal(b, m, deterministic)
}
func (m *MutateAdGroupExtensionSettingsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateAdGroupExtensionSettingsResponse.Merge(m, src)
}
func (m *MutateAdGroupExtensionSettingsResponse) XXX_Size() int {
	return xxx_messageInfo_MutateAdGroupExtensionSettingsResponse.Size(m)
}
func (m *MutateAdGroupExtensionSettingsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateAdGroupExtensionSettingsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MutateAdGroupExtensionSettingsResponse proto.InternalMessageInfo

func (m *MutateAdGroupExtensionSettingsResponse) GetPartialFailureError() *status.Status {
	if m != nil {
		return m.PartialFailureError
	}
	return nil
}

func (m *MutateAdGroupExtensionSettingsResponse) GetResults() []*MutateAdGroupExtensionSettingResult {
	if m != nil {
		return m.Results
	}
	return nil
}

// The result for the ad group extension setting mutate.
type MutateAdGroupExtensionSettingResult struct {
	// Returned for successful operations.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateAdGroupExtensionSettingResult) Reset()         { *m = MutateAdGroupExtensionSettingResult{} }
func (m *MutateAdGroupExtensionSettingResult) String() string { return proto.CompactTextString(m) }
func (*MutateAdGroupExtensionSettingResult) ProtoMessage()    {}
func (*MutateAdGroupExtensionSettingResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf10cd4825f3524e, []int{4}
}

func (m *MutateAdGroupExtensionSettingResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateAdGroupExtensionSettingResult.Unmarshal(m, b)
}
func (m *MutateAdGroupExtensionSettingResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateAdGroupExtensionSettingResult.Marshal(b, m, deterministic)
}
func (m *MutateAdGroupExtensionSettingResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateAdGroupExtensionSettingResult.Merge(m, src)
}
func (m *MutateAdGroupExtensionSettingResult) XXX_Size() int {
	return xxx_messageInfo_MutateAdGroupExtensionSettingResult.Size(m)
}
func (m *MutateAdGroupExtensionSettingResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateAdGroupExtensionSettingResult.DiscardUnknown(m)
}

var xxx_messageInfo_MutateAdGroupExtensionSettingResult proto.InternalMessageInfo

func (m *MutateAdGroupExtensionSettingResult) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetAdGroupExtensionSettingRequest)(nil), "google.ads.googleads.v2.services.GetAdGroupExtensionSettingRequest")
	proto.RegisterType((*MutateAdGroupExtensionSettingsRequest)(nil), "google.ads.googleads.v2.services.MutateAdGroupExtensionSettingsRequest")
	proto.RegisterType((*AdGroupExtensionSettingOperation)(nil), "google.ads.googleads.v2.services.AdGroupExtensionSettingOperation")
	proto.RegisterType((*MutateAdGroupExtensionSettingsResponse)(nil), "google.ads.googleads.v2.services.MutateAdGroupExtensionSettingsResponse")
	proto.RegisterType((*MutateAdGroupExtensionSettingResult)(nil), "google.ads.googleads.v2.services.MutateAdGroupExtensionSettingResult")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v2/services/ad_group_extension_setting_service.proto", fileDescriptor_bf10cd4825f3524e)
}

var fileDescriptor_bf10cd4825f3524e = []byte{
	// 776 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x4d, 0x6f, 0xd3, 0x4a,
	0x14, 0x7d, 0x71, 0xaa, 0xbe, 0xd7, 0x49, 0xfb, 0x9e, 0x34, 0x4f, 0x40, 0x14, 0x50, 0x09, 0x6e,
	0x81, 0x28, 0x42, 0xb6, 0x64, 0x76, 0x2e, 0x5d, 0x38, 0xa8, 0x4d, 0x8b, 0x54, 0x5a, 0xb9, 0xa8,
	0x0b, 0x14, 0xc9, 0x9a, 0xda, 0x53, 0x63, 0xd5, 0xf6, 0x98, 0x99, 0x71, 0x44, 0x55, 0x75, 0x83,
	0x58, 0xb1, 0x45, 0xe2, 0x07, 0xb0, 0x64, 0xc9, 0xcf, 0xe8, 0x0e, 0xb1, 0xeb, 0x8a, 0x05, 0x2b,
	0x76, 0xfc, 0x03, 0x64, 0x8f, 0x27, 0x1f, 0x15, 0x8e, 0x91, 0xca, 0xee, 0x7a, 0xee, 0xc9, 0x39,
	0x77, 0xee, 0x3d, 0x73, 0x03, 0xb6, 0x7d, 0x42, 0xfc, 0x10, 0xeb, 0xc8, 0x63, 0xba, 0x08, 0xb3,
	0x68, 0x68, 0xe8, 0x0c, 0xd3, 0x61, 0xe0, 0x62, 0xa6, 0x23, 0xcf, 0xf1, 0x29, 0x49, 0x13, 0x07,
	0xbf, 0xe2, 0x38, 0x66, 0x01, 0x89, 0x1d, 0x86, 0x39, 0x0f, 0x62, 0xdf, 0x29, 0x30, 0x5a, 0x42,
	0x09, 0x27, 0xb0, 0x2d, 0x7e, 0xaf, 0x21, 0x8f, 0x69, 0x23, 0x2a, 0x6d, 0x68, 0x68, 0x92, 0xaa,
	0xd5, 0x2b, 0x13, 0xa3, 0x98, 0x91, 0x94, 0xce, 0x56, 0x13, 0x2a, 0xad, 0x5b, 0x92, 0x23, 0x09,
	0x74, 0x14, 0xc7, 0x84, 0x23, 0x1e, 0x90, 0x98, 0x15, 0xd9, 0x1b, 0x13, 0x59, 0x37, 0x0c, 0x70,
	0xcc, 0x8b, 0xc4, 0xed, 0x89, 0xc4, 0x51, 0x80, 0x43, 0xcf, 0x39, 0xc4, 0x2f, 0xd0, 0x30, 0x20,
	0xb4, 0x00, 0x14, 0xd5, 0xeb, 0xf9, 0xd7, 0x61, 0x7a, 0x54, 0xa0, 0x22, 0xc4, 0x8e, 0x2f, 0x71,
	0xd3, 0xc4, 0xd5, 0x19, 0x47, 0x3c, 0x2d, 0x44, 0xd5, 0x1d, 0x70, 0xa7, 0x8f, 0xb9, 0xe5, 0xf5,
	0xb3, 0xc2, 0x37, 0x64, 0xdd, 0xfb, 0xa2, 0x6c, 0x1b, 0xbf, 0x4c, 0x31, 0xe3, 0xb0, 0x03, 0x96,
	0xe4, 0x2d, 0x9d, 0x18, 0x45, 0xb8, 0x59, 0x6b, 0xd7, 0x3a, 0x0b, 0xbd, 0xfa, 0x57, 0x4b, 0xb1,
	0x17, 0x65, 0xe6, 0x29, 0x8a, 0xb0, 0xfa, 0x46, 0x01, 0x77, 0x77, 0x52, 0x8e, 0x38, 0x2e, 0xa1,
	0x64, 0x92, 0x73, 0x15, 0x34, 0xdc, 0x94, 0x71, 0x12, 0x61, 0xea, 0x04, 0xde, 0x24, 0x23, 0x90,
	0xe7, 0xdb, 0x1e, 0xf4, 0x01, 0x20, 0x09, 0xa6, 0xa2, 0x4f, 0x4d, 0xa5, 0x5d, 0xef, 0x34, 0x8c,
	0x9e, 0x56, 0x35, 0x2c, 0xad, 0x44, 0x7c, 0x57, 0x52, 0x15, 0x42, 0x63, 0x6a, 0x78, 0x1f, 0xfc,
	0x97, 0x20, 0xca, 0x03, 0x14, 0x3a, 0x47, 0x28, 0x08, 0x53, 0x8a, 0x9b, 0xf5, 0x76, 0xad, 0xf3,
	0x8f, 0xfd, 0x6f, 0x71, 0xbc, 0x29, 0x4e, 0xe1, 0x0a, 0x58, 0x1a, 0xa2, 0x30, 0xf0, 0x10, 0xc7,
	0x0e, 0x89, 0xc3, 0x93, 0xe6, 0x5c, 0x0e, 0x5b, 0x94, 0x87, 0xbb, 0x71, 0x78, 0xa2, 0x7e, 0x52,
	0x40, 0xbb, 0xaa, 0x06, 0xb8, 0x06, 0x1a, 0x69, 0x92, 0xf3, 0x64, 0x83, 0xca, 0x79, 0x1a, 0x46,
	0x4b, 0x5e, 0x4e, 0xce, 0x52, 0xdb, 0xcc, 0x66, 0xb9, 0x83, 0xd8, 0xb1, 0x0d, 0x04, 0x3c, 0x8b,
	0xe1, 0x33, 0x30, 0xef, 0x52, 0x8c, 0xb8, 0x98, 0x45, 0xc3, 0x30, 0x4b, 0x9b, 0x32, 0xf2, 0x67,
	0x59, 0x57, 0xb6, 0xfe, 0xb2, 0x0b, 0xae, 0x8c, 0x55, 0x68, 0x34, 0x95, 0x3f, 0xc1, 0x2a, 0xb8,
	0x60, 0x13, 0xcc, 0x53, 0x1c, 0x91, 0xa1, 0x68, 0xe9, 0x42, 0x96, 0x11, 0xdf, 0xbd, 0x06, 0x58,
	0x18, 0xcd, 0x40, 0xfd, 0x5c, 0x03, 0xf7, 0xaa, 0xbc, 0xc3, 0x12, 0x12, 0x33, 0x0c, 0x37, 0xc1,
	0xb5, 0x4b, 0xd3, 0x72, 0x30, 0xa5, 0x84, 0xe6, 0x02, 0x0d, 0x03, 0xca, 0xb2, 0x69, 0xe2, 0x6a,
	0xfb, 0xb9, 0xdd, 0xed, 0xff, 0xa7, 0xe7, 0xb8, 0x91, 0xc1, 0xa1, 0x03, 0xfe, 0xa6, 0x98, 0xa5,
	0x21, 0x97, 0xde, 0xda, 0xa8, 0xf6, 0xd6, 0xcc, 0x12, 0xed, 0x9c, 0xcd, 0x96, 0xac, 0xea, 0x13,
	0xb0, 0xf2, 0x1b, 0xf8, 0xcc, 0x54, 0xbf, 0x78, 0x60, 0xd3, 0x6f, 0xcb, 0x38, 0x9f, 0x03, 0xcb,
	0x25, 0x34, 0xfb, 0xa2, 0x38, 0xf8, 0xa3, 0x06, 0x5a, 0xe5, 0xcf, 0x19, 0x3e, 0xae, 0xbe, 0x5d,
	0xe5, 0x32, 0x68, 0x5d, 0xc1, 0x13, 0xaa, 0x7d, 0x61, 0x4d, 0x5f, 0xf4, 0xf5, 0x97, 0x6f, 0xef,
	0x94, 0x47, 0xd0, 0xcc, 0x16, 0xe9, 0xe9, 0x54, 0x66, 0x5d, 0x6e, 0x01, 0xa6, 0x77, 0x75, 0x54,
	0x62, 0x08, 0xbd, 0x7b, 0x06, 0xdf, 0x2b, 0x60, 0x79, 0xb6, 0x6d, 0x60, 0xff, 0x8a, 0x53, 0x95,
	0x4b, 0xab, 0xb5, 0x75, 0x75, 0x22, 0xe1, 0x60, 0x15, 0x5d, 0x58, 0xd7, 0x27, 0xf6, 0xdf, 0x83,
	0xf1, 0x2a, 0xca, 0x5b, 0xd2, 0x53, 0xd7, 0xb3, 0x96, 0x8c, 0x7b, 0x70, 0x3a, 0x01, 0x5e, 0xef,
	0x9e, 0x95, 0x76, 0xc4, 0x8c, 0x72, 0x7d, 0xb3, 0xd6, 0x6d, 0xdd, 0x3c, 0xb7, 0x9a, 0xe3, 0x1a,
	0x8b, 0x28, 0x09, 0x98, 0xe6, 0x92, 0xa8, 0xf7, 0x56, 0x01, 0xab, 0x2e, 0x89, 0x2a, 0xef, 0xd3,
	0x5b, 0x99, 0x6d, 0xb9, 0xbd, 0x6c, 0x4d, 0xed, 0xd5, 0x9e, 0x6f, 0x15, 0x44, 0x3e, 0x09, 0x51,
	0xec, 0x6b, 0x84, 0xfa, 0xba, 0x8f, 0xe3, 0x7c, 0x89, 0xe9, 0x63, 0xe9, 0xf2, 0xbf, 0xea, 0x35,
	0x19, 0x7c, 0x50, 0xea, 0x7d, 0xcb, 0xfa, 0xa8, 0xb4, 0xfb, 0x82, 0xd0, 0xf2, 0x98, 0x26, 0xc2,
	0x2c, 0x3a, 0x30, 0xb4, 0x42, 0x98, 0x9d, 0x4b, 0xc8, 0xc0, 0xf2, 0xd8, 0x60, 0x04, 0x19, 0x1c,
	0x18, 0x03, 0x09, 0xf9, 0xae, 0xac, 0x8a, 0x73, 0xd3, 0xb4, 0x3c, 0x66, 0x9a, 0x23, 0x90, 0x69,
	0x1e, 0x18, 0xa6, 0x29, 0x61, 0x87, 0xf3, 0x79, 0x9d, 0x0f, 0x7f, 0x06, 0x00, 0x00, 0xff, 0xff,
	0xbf, 0x21, 0x55, 0xcd, 0x51, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdGroupExtensionSettingServiceClient is the client API for AdGroupExtensionSettingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdGroupExtensionSettingServiceClient interface {
	// Returns the requested ad group extension setting in full detail.
	GetAdGroupExtensionSetting(ctx context.Context, in *GetAdGroupExtensionSettingRequest, opts ...grpc.CallOption) (*resources.AdGroupExtensionSetting, error)
	// Creates, updates, or removes ad group extension settings. Operation
	// statuses are returned.
	MutateAdGroupExtensionSettings(ctx context.Context, in *MutateAdGroupExtensionSettingsRequest, opts ...grpc.CallOption) (*MutateAdGroupExtensionSettingsResponse, error)
}

type adGroupExtensionSettingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdGroupExtensionSettingServiceClient(cc grpc.ClientConnInterface) AdGroupExtensionSettingServiceClient {
	return &adGroupExtensionSettingServiceClient{cc}
}

func (c *adGroupExtensionSettingServiceClient) GetAdGroupExtensionSetting(ctx context.Context, in *GetAdGroupExtensionSettingRequest, opts ...grpc.CallOption) (*resources.AdGroupExtensionSetting, error) {
	out := new(resources.AdGroupExtensionSetting)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v2.services.AdGroupExtensionSettingService/GetAdGroupExtensionSetting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adGroupExtensionSettingServiceClient) MutateAdGroupExtensionSettings(ctx context.Context, in *MutateAdGroupExtensionSettingsRequest, opts ...grpc.CallOption) (*MutateAdGroupExtensionSettingsResponse, error) {
	out := new(MutateAdGroupExtensionSettingsResponse)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v2.services.AdGroupExtensionSettingService/MutateAdGroupExtensionSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdGroupExtensionSettingServiceServer is the server API for AdGroupExtensionSettingService service.
type AdGroupExtensionSettingServiceServer interface {
	// Returns the requested ad group extension setting in full detail.
	GetAdGroupExtensionSetting(context.Context, *GetAdGroupExtensionSettingRequest) (*resources.AdGroupExtensionSetting, error)
	// Creates, updates, or removes ad group extension settings. Operation
	// statuses are returned.
	MutateAdGroupExtensionSettings(context.Context, *MutateAdGroupExtensionSettingsRequest) (*MutateAdGroupExtensionSettingsResponse, error)
}

// UnimplementedAdGroupExtensionSettingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAdGroupExtensionSettingServiceServer struct {
}

func (*UnimplementedAdGroupExtensionSettingServiceServer) GetAdGroupExtensionSetting(ctx context.Context, req *GetAdGroupExtensionSettingRequest) (*resources.AdGroupExtensionSetting, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method GetAdGroupExtensionSetting not implemented")
}
func (*UnimplementedAdGroupExtensionSettingServiceServer) MutateAdGroupExtensionSettings(ctx context.Context, req *MutateAdGroupExtensionSettingsRequest) (*MutateAdGroupExtensionSettingsResponse, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method MutateAdGroupExtensionSettings not implemented")
}

func RegisterAdGroupExtensionSettingServiceServer(s *grpc.Server, srv AdGroupExtensionSettingServiceServer) {
	s.RegisterService(&_AdGroupExtensionSettingService_serviceDesc, srv)
}

func _AdGroupExtensionSettingService_GetAdGroupExtensionSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdGroupExtensionSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdGroupExtensionSettingServiceServer).GetAdGroupExtensionSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v2.services.AdGroupExtensionSettingService/GetAdGroupExtensionSetting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdGroupExtensionSettingServiceServer).GetAdGroupExtensionSetting(ctx, req.(*GetAdGroupExtensionSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdGroupExtensionSettingService_MutateAdGroupExtensionSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateAdGroupExtensionSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdGroupExtensionSettingServiceServer).MutateAdGroupExtensionSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v2.services.AdGroupExtensionSettingService/MutateAdGroupExtensionSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdGroupExtensionSettingServiceServer).MutateAdGroupExtensionSettings(ctx, req.(*MutateAdGroupExtensionSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdGroupExtensionSettingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v2.services.AdGroupExtensionSettingService",
	HandlerType: (*AdGroupExtensionSettingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAdGroupExtensionSetting",
			Handler:    _AdGroupExtensionSettingService_GetAdGroupExtensionSetting_Handler,
		},
		{
			MethodName: "MutateAdGroupExtensionSettings",
			Handler:    _AdGroupExtensionSettingService_MutateAdGroupExtensionSettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v2/services/ad_group_extension_setting_service.proto",
}
