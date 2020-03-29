// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/resources/geographic_view.proto

package resources

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

// A geographic view.
//
// Geographic View includes all metrics aggregated at the country level,
// one row per country. It reports metrics at either actual physical location of
// the user or an area of interest. If other segment fields are used, you may
// get more than one row per country.
type GeographicView struct {
	// Output only. The resource name of the geographic view.
	// Geographic view resource names have the form:
	//
	// `customers/{customer_id}/geographicViews/{country_criterion_id}~{location_type}`
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	// Output only. CriterionId for the geo target for a country.
	CountryGeoTargetConstant *wrappers.StringValue `protobuf:"bytes,2,opt,name=country_geo_target_constant,json=countryGeoTargetConstant,proto3" json:"country_geo_target_constant,omitempty"`
	// Output only. Type of the geo targeting of the campaign.
	LocationType         enums.GeoTargetingTypeEnum_GeoTargetingType `protobuf:"varint,3,opt,name=location_type,json=locationType,proto3,enum=google.ads.googleads.v1.enums.GeoTargetingTypeEnum_GeoTargetingType" json:"location_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *GeographicView) Reset()         { *m = GeographicView{} }
func (m *GeographicView) String() string { return proto.CompactTextString(m) }
func (*GeographicView) ProtoMessage()    {}
func (*GeographicView) Descriptor() ([]byte, []int) {
	return fileDescriptor_45a03e49a5e672c5, []int{0}
}

func (m *GeographicView) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeographicView.Unmarshal(m, b)
}
func (m *GeographicView) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeographicView.Marshal(b, m, deterministic)
}
func (m *GeographicView) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeographicView.Merge(m, src)
}
func (m *GeographicView) XXX_Size() int {
	return xxx_messageInfo_GeographicView.Size(m)
}
func (m *GeographicView) XXX_DiscardUnknown() {
	xxx_messageInfo_GeographicView.DiscardUnknown(m)
}

var xxx_messageInfo_GeographicView proto.InternalMessageInfo

func (m *GeographicView) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func (m *GeographicView) GetCountryGeoTargetConstant() *wrappers.StringValue {
	if m != nil {
		return m.CountryGeoTargetConstant
	}
	return nil
}

func (m *GeographicView) GetLocationType() enums.GeoTargetingTypeEnum_GeoTargetingType {
	if m != nil {
		return m.LocationType
	}
	return enums.GeoTargetingTypeEnum_UNSPECIFIED
}

func init() {
	proto.RegisterType((*GeographicView)(nil), "google.ads.googleads.v1.resources.GeographicView")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/resources/geographic_view.proto", fileDescriptor_45a03e49a5e672c5)
}

var fileDescriptor_45a03e49a5e672c5 = []byte{
	// 496 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcf, 0x6a, 0xd4, 0x40,
	0x1c, 0x26, 0x09, 0x08, 0xc6, 0xb6, 0x87, 0x78, 0x59, 0x6b, 0xd1, 0xad, 0x50, 0x5c, 0x45, 0x66,
	0xc8, 0x0a, 0x15, 0xe2, 0x29, 0x51, 0x59, 0xf0, 0x20, 0x65, 0x5d, 0x72, 0x90, 0x85, 0x30, 0x9b,
	0xfc, 0x3a, 0x1d, 0x48, 0x66, 0xc2, 0xcc, 0x64, 0x97, 0xa5, 0x14, 0x7c, 0x03, 0xdf, 0xc1, 0xa3,
	0x8f, 0xe2, 0xc1, 0x67, 0xe8, 0xb9, 0x8f, 0xe0, 0x49, 0x92, 0x4c, 0xb2, 0x5b, 0x4b, 0xab, 0xb7,
	0x2f, 0xf3, 0x7d, 0xbf, 0xef, 0xf7, 0x37, 0xee, 0x1b, 0x2a, 0x04, 0xcd, 0x01, 0x93, 0x4c, 0xe1,
	0x16, 0xd6, 0x68, 0xe9, 0x63, 0x09, 0x4a, 0x54, 0x32, 0x05, 0x85, 0x29, 0x08, 0x2a, 0x49, 0x79,
	0xc6, 0xd2, 0x64, 0xc9, 0x60, 0x85, 0x4a, 0x29, 0xb4, 0xf0, 0x0e, 0x5b, 0x35, 0x22, 0x99, 0x42,
	0x7d, 0x20, 0x5a, 0xfa, 0xa8, 0x0f, 0xdc, 0x3f, 0xbe, 0xcd, 0x1b, 0x78, 0x55, 0x34, 0xbe, 0x89,
	0x26, 0x92, 0x82, 0x66, 0x9c, 0x26, 0x7a, 0x5d, 0x42, 0x6b, 0xbd, 0xff, 0xb4, 0x8b, 0x2b, 0x19,
	0x3e, 0x65, 0x90, 0x67, 0xc9, 0x02, 0xce, 0xc8, 0x92, 0x09, 0x69, 0x04, 0x8f, 0xb6, 0x04, 0x5d,
	0x3a, 0x43, 0x3d, 0x31, 0x54, 0xf3, 0xb5, 0xa8, 0x4e, 0xf1, 0x4a, 0x92, 0xb2, 0x04, 0xa9, 0x0c,
	0x7f, 0xb0, 0x15, 0x4a, 0x38, 0x17, 0x9a, 0x68, 0x26, 0xb8, 0x61, 0x9f, 0xfd, 0x72, 0xdc, 0xbd,
	0x49, 0xdf, 0x6e, 0xcc, 0x60, 0xe5, 0xcd, 0xdc, 0xdd, 0x2e, 0x45, 0xc2, 0x49, 0x01, 0x03, 0x6b,
	0x68, 0x8d, 0xee, 0x47, 0xf8, 0x32, 0x74, 0x7e, 0x87, 0x2f, 0xdc, 0xe7, 0x9b, 0xde, 0x0d, 0x2a,
	0x99, 0x42, 0xa9, 0x28, 0xf0, 0x75, 0x9f, 0xe9, 0x4e, 0xe7, 0xf2, 0x89, 0x14, 0xe0, 0x7d, 0xb3,
	0xdc, 0xc7, 0xa9, 0xa8, 0xb8, 0x96, 0xeb, 0x64, 0x33, 0x87, 0x24, 0x15, 0x5c, 0x69, 0xc2, 0xf5,
	0xc0, 0x1e, 0x5a, 0xa3, 0x07, 0xe3, 0x03, 0xe3, 0x89, 0xba, 0x6e, 0xd0, 0x67, 0x2d, 0x19, 0xa7,
	0x31, 0xc9, 0x2b, 0x88, 0xc6, 0x4d, 0x09, 0xaf, 0xdc, 0x97, 0x77, 0x95, 0x30, 0x6b, 0x8c, 0xdf,
	0x19, 0xdf, 0xe9, 0xc0, 0x24, 0xbd, 0xc1, 0x78, 0xdc, 0xdd, 0xcd, 0x45, 0xda, 0x4c, 0xa3, 0xd9,
	0xc5, 0xc0, 0x19, 0x5a, 0xa3, 0xbd, 0xf1, 0x7b, 0x74, 0xdb, 0x9e, 0x9b, 0x25, 0xa2, 0xde, 0x88,
	0x71, 0x3a, 0x5b, 0x97, 0xf0, 0x81, 0x57, 0xc5, 0x8d, 0xc7, 0xc8, 0xb9, 0x0c, 0x9d, 0xe9, 0x4e,
	0xe7, 0x5f, 0x3f, 0x05, 0xd9, 0x55, 0x48, 0xfe, 0x7b, 0x7a, 0xde, 0x71, 0x5a, 0x29, 0x2d, 0x0a,
	0x90, 0x0a, 0x9f, 0x77, 0xf0, 0x62, 0xeb, 0x32, 0x6b, 0x91, 0xc2, 0xe7, 0x7f, 0x9d, 0xea, 0x45,
	0xf4, 0xd5, 0x76, 0x8f, 0x52, 0x51, 0xa0, 0x7f, 0x1e, 0x6b, 0xf4, 0xf0, 0x7a, 0xc6, 0x93, 0x7a,
	0xe2, 0x27, 0xd6, 0x97, 0x8f, 0x26, 0x92, 0x8a, 0x9c, 0x70, 0x8a, 0x84, 0xa4, 0x98, 0x02, 0x6f,
	0xf6, 0x81, 0x37, 0x25, 0xdf, 0xf1, 0xfb, 0xbc, 0xed, 0xd1, 0x77, 0xdb, 0x99, 0x84, 0xe1, 0x0f,
	0xfb, 0x70, 0xd2, 0x5a, 0x86, 0x99, 0x42, 0x2d, 0xac, 0x51, 0xec, 0xa3, 0x69, 0xa7, 0xfc, 0xd9,
	0x69, 0xe6, 0x61, 0xa6, 0xe6, 0xbd, 0x66, 0x1e, 0xfb, 0xf3, 0x5e, 0x73, 0x65, 0x1f, 0xb5, 0x44,
	0x10, 0x84, 0x99, 0x0a, 0x82, 0x5e, 0x15, 0x04, 0xb1, 0x1f, 0x04, 0xbd, 0x6e, 0x71, 0xaf, 0x29,
	0xf6, 0xf5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x21, 0xa0, 0xf7, 0xea, 0x03, 0x00, 0x00,
}
