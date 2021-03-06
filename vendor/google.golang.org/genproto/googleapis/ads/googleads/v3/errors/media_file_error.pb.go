// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/errors/media_file_error.proto

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

// Enum describing possible media file errors.
type MediaFileErrorEnum_MediaFileError int32

const (
	// Enum unspecified.
	MediaFileErrorEnum_UNSPECIFIED MediaFileErrorEnum_MediaFileError = 0
	// The received error code is not known in this version.
	MediaFileErrorEnum_UNKNOWN MediaFileErrorEnum_MediaFileError = 1
	// Cannot create a standard icon type.
	MediaFileErrorEnum_CANNOT_CREATE_STANDARD_ICON MediaFileErrorEnum_MediaFileError = 2
	// May only select Standard Icons alone.
	MediaFileErrorEnum_CANNOT_SELECT_STANDARD_ICON_WITH_OTHER_TYPES MediaFileErrorEnum_MediaFileError = 3
	// Image contains both a media file ID and data.
	MediaFileErrorEnum_CANNOT_SPECIFY_MEDIA_FILE_ID_AND_DATA MediaFileErrorEnum_MediaFileError = 4
	// A media file with given type and reference ID already exists.
	MediaFileErrorEnum_DUPLICATE_MEDIA MediaFileErrorEnum_MediaFileError = 5
	// A required field was not specified or is an empty string.
	MediaFileErrorEnum_EMPTY_FIELD MediaFileErrorEnum_MediaFileError = 6
	// A media file may only be modified once per call.
	MediaFileErrorEnum_RESOURCE_REFERENCED_IN_MULTIPLE_OPS MediaFileErrorEnum_MediaFileError = 7
	// Field is not supported for the media sub type.
	MediaFileErrorEnum_FIELD_NOT_SUPPORTED_FOR_MEDIA_SUB_TYPE MediaFileErrorEnum_MediaFileError = 8
	// The media file ID is invalid.
	MediaFileErrorEnum_INVALID_MEDIA_FILE_ID MediaFileErrorEnum_MediaFileError = 9
	// The media subtype is invalid.
	MediaFileErrorEnum_INVALID_MEDIA_SUB_TYPE MediaFileErrorEnum_MediaFileError = 10
	// The media file type is invalid.
	MediaFileErrorEnum_INVALID_MEDIA_FILE_TYPE MediaFileErrorEnum_MediaFileError = 11
	// The mimetype is invalid.
	MediaFileErrorEnum_INVALID_MIME_TYPE MediaFileErrorEnum_MediaFileError = 12
	// The media reference ID is invalid.
	MediaFileErrorEnum_INVALID_REFERENCE_ID MediaFileErrorEnum_MediaFileError = 13
	// The YouTube video ID is invalid.
	MediaFileErrorEnum_INVALID_YOU_TUBE_ID MediaFileErrorEnum_MediaFileError = 14
	// Media file has failed transcoding
	MediaFileErrorEnum_MEDIA_FILE_FAILED_TRANSCODING MediaFileErrorEnum_MediaFileError = 15
	// Media file has not been transcoded.
	MediaFileErrorEnum_MEDIA_NOT_TRANSCODED MediaFileErrorEnum_MediaFileError = 16
	// The media type does not match the actual media file's type.
	MediaFileErrorEnum_MEDIA_TYPE_DOES_NOT_MATCH_MEDIA_FILE_TYPE MediaFileErrorEnum_MediaFileError = 17
	// None of the fields have been specified.
	MediaFileErrorEnum_NO_FIELDS_SPECIFIED MediaFileErrorEnum_MediaFileError = 18
	// One of reference ID or media file ID must be specified.
	MediaFileErrorEnum_NULL_REFERENCE_ID_AND_MEDIA_ID MediaFileErrorEnum_MediaFileError = 19
	// The string has too many characters.
	MediaFileErrorEnum_TOO_LONG MediaFileErrorEnum_MediaFileError = 20
	// The specified type is not supported.
	MediaFileErrorEnum_UNSUPPORTED_TYPE MediaFileErrorEnum_MediaFileError = 21
	// YouTube is unavailable for requesting video data.
	MediaFileErrorEnum_YOU_TUBE_SERVICE_UNAVAILABLE MediaFileErrorEnum_MediaFileError = 22
	// The YouTube video has a non positive duration.
	MediaFileErrorEnum_YOU_TUBE_VIDEO_HAS_NON_POSITIVE_DURATION MediaFileErrorEnum_MediaFileError = 23
	// The YouTube video ID is syntactically valid but the video was not found.
	MediaFileErrorEnum_YOU_TUBE_VIDEO_NOT_FOUND MediaFileErrorEnum_MediaFileError = 24
)

var MediaFileErrorEnum_MediaFileError_name = map[int32]string{
	0:  "UNSPECIFIED",
	1:  "UNKNOWN",
	2:  "CANNOT_CREATE_STANDARD_ICON",
	3:  "CANNOT_SELECT_STANDARD_ICON_WITH_OTHER_TYPES",
	4:  "CANNOT_SPECIFY_MEDIA_FILE_ID_AND_DATA",
	5:  "DUPLICATE_MEDIA",
	6:  "EMPTY_FIELD",
	7:  "RESOURCE_REFERENCED_IN_MULTIPLE_OPS",
	8:  "FIELD_NOT_SUPPORTED_FOR_MEDIA_SUB_TYPE",
	9:  "INVALID_MEDIA_FILE_ID",
	10: "INVALID_MEDIA_SUB_TYPE",
	11: "INVALID_MEDIA_FILE_TYPE",
	12: "INVALID_MIME_TYPE",
	13: "INVALID_REFERENCE_ID",
	14: "INVALID_YOU_TUBE_ID",
	15: "MEDIA_FILE_FAILED_TRANSCODING",
	16: "MEDIA_NOT_TRANSCODED",
	17: "MEDIA_TYPE_DOES_NOT_MATCH_MEDIA_FILE_TYPE",
	18: "NO_FIELDS_SPECIFIED",
	19: "NULL_REFERENCE_ID_AND_MEDIA_ID",
	20: "TOO_LONG",
	21: "UNSUPPORTED_TYPE",
	22: "YOU_TUBE_SERVICE_UNAVAILABLE",
	23: "YOU_TUBE_VIDEO_HAS_NON_POSITIVE_DURATION",
	24: "YOU_TUBE_VIDEO_NOT_FOUND",
}

var MediaFileErrorEnum_MediaFileError_value = map[string]int32{
	"UNSPECIFIED":                 0,
	"UNKNOWN":                     1,
	"CANNOT_CREATE_STANDARD_ICON": 2,
	"CANNOT_SELECT_STANDARD_ICON_WITH_OTHER_TYPES": 3,
	"CANNOT_SPECIFY_MEDIA_FILE_ID_AND_DATA":        4,
	"DUPLICATE_MEDIA":                              5,
	"EMPTY_FIELD":                                  6,
	"RESOURCE_REFERENCED_IN_MULTIPLE_OPS":          7,
	"FIELD_NOT_SUPPORTED_FOR_MEDIA_SUB_TYPE":       8,
	"INVALID_MEDIA_FILE_ID":                        9,
	"INVALID_MEDIA_SUB_TYPE":                       10,
	"INVALID_MEDIA_FILE_TYPE":                      11,
	"INVALID_MIME_TYPE":                            12,
	"INVALID_REFERENCE_ID":                         13,
	"INVALID_YOU_TUBE_ID":                          14,
	"MEDIA_FILE_FAILED_TRANSCODING":                15,
	"MEDIA_NOT_TRANSCODED":                         16,
	"MEDIA_TYPE_DOES_NOT_MATCH_MEDIA_FILE_TYPE":    17,
	"NO_FIELDS_SPECIFIED":                          18,
	"NULL_REFERENCE_ID_AND_MEDIA_ID":               19,
	"TOO_LONG":                                     20,
	"UNSUPPORTED_TYPE":                             21,
	"YOU_TUBE_SERVICE_UNAVAILABLE":                 22,
	"YOU_TUBE_VIDEO_HAS_NON_POSITIVE_DURATION":     23,
	"YOU_TUBE_VIDEO_NOT_FOUND":                     24,
}

func (x MediaFileErrorEnum_MediaFileError) String() string {
	return proto.EnumName(MediaFileErrorEnum_MediaFileError_name, int32(x))
}

func (MediaFileErrorEnum_MediaFileError) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a6f239b6b744c564, []int{0, 0}
}

// Container for enum describing possible media file errors.
type MediaFileErrorEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MediaFileErrorEnum) Reset()         { *m = MediaFileErrorEnum{} }
func (m *MediaFileErrorEnum) String() string { return proto.CompactTextString(m) }
func (*MediaFileErrorEnum) ProtoMessage()    {}
func (*MediaFileErrorEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6f239b6b744c564, []int{0}
}

func (m *MediaFileErrorEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MediaFileErrorEnum.Unmarshal(m, b)
}
func (m *MediaFileErrorEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MediaFileErrorEnum.Marshal(b, m, deterministic)
}
func (m *MediaFileErrorEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MediaFileErrorEnum.Merge(m, src)
}
func (m *MediaFileErrorEnum) XXX_Size() int {
	return xxx_messageInfo_MediaFileErrorEnum.Size(m)
}
func (m *MediaFileErrorEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_MediaFileErrorEnum.DiscardUnknown(m)
}

var xxx_messageInfo_MediaFileErrorEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v3.errors.MediaFileErrorEnum_MediaFileError", MediaFileErrorEnum_MediaFileError_name, MediaFileErrorEnum_MediaFileError_value)
	proto.RegisterType((*MediaFileErrorEnum)(nil), "google.ads.googleads.v3.errors.MediaFileErrorEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/errors/media_file_error.proto", fileDescriptor_a6f239b6b744c564)
}

var fileDescriptor_a6f239b6b744c564 = []byte{
	// 672 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0xa5, 0x29, 0xa4, 0x65, 0x5a, 0xda, 0xe9, 0xa4, 0x2f, 0xda, 0x52, 0x20, 0x88, 0x47, 0x51,
	0x71, 0x90, 0x22, 0x36, 0x66, 0x35, 0xf1, 0x5c, 0x27, 0x23, 0x9c, 0x19, 0xcb, 0x1e, 0xbb, 0x0a,
	0x8a, 0x34, 0x0a, 0x24, 0x44, 0x91, 0xd2, 0xb8, 0x8a, 0x4b, 0xd7, 0xfc, 0x05, 0x7b, 0x96, 0x7c,
	0x0a, 0x9f, 0xc2, 0x82, 0x6f, 0x40, 0xf6, 0x24, 0x2e, 0xe1, 0xb5, 0xca, 0xd5, 0x3d, 0xe7, 0xdc,
	0x73, 0xee, 0xc4, 0x17, 0xbd, 0x1a, 0x26, 0xc9, 0x70, 0x3c, 0xa8, 0xf5, 0xfa, 0x69, 0xcd, 0x94,
	0x59, 0x75, 0x55, 0xaf, 0x0d, 0xa6, 0xd3, 0x64, 0x9a, 0xd6, 0xce, 0x07, 0xfd, 0x51, 0x4f, 0x7f,
	0x18, 0x8d, 0x07, 0x3a, 0xef, 0x58, 0x17, 0xd3, 0xe4, 0x32, 0x21, 0xc7, 0x86, 0x6b, 0xf5, 0xfa,
	0xa9, 0x55, 0xc8, 0xac, 0xab, 0xba, 0x65, 0x64, 0x07, 0x47, 0xf3, 0xb1, 0x17, 0xa3, 0x5a, 0x6f,
	0x32, 0x49, 0x2e, 0x7b, 0x97, 0xa3, 0x64, 0x92, 0x1a, 0x75, 0xf5, 0x73, 0x19, 0x91, 0x76, 0x36,
	0xd8, 0x1d, 0x8d, 0x07, 0x90, 0x29, 0x60, 0xf2, 0xf1, 0xbc, 0xfa, 0xa9, 0x8c, 0x36, 0x16, 0xdb,
	0x64, 0x13, 0xad, 0x45, 0x22, 0xf4, 0xc1, 0xe1, 0x2e, 0x07, 0x86, 0x6f, 0x90, 0x35, 0xb4, 0x12,
	0x89, 0x37, 0x42, 0x9e, 0x09, 0xbc, 0x44, 0xee, 0xa3, 0x43, 0x87, 0x0a, 0x21, 0x95, 0x76, 0x02,
	0xa0, 0x0a, 0x74, 0xa8, 0xa8, 0x60, 0x34, 0x60, 0x9a, 0x3b, 0x52, 0xe0, 0x12, 0x79, 0x89, 0x4e,
	0x67, 0x84, 0x10, 0x3c, 0x70, 0xd4, 0x22, 0x41, 0x9f, 0x71, 0xd5, 0xd2, 0x52, 0xb5, 0x20, 0xd0,
	0xaa, 0xe3, 0x43, 0x88, 0x97, 0xc9, 0x09, 0x7a, 0x3c, 0x57, 0xe4, 0xae, 0x1d, 0xdd, 0x06, 0xc6,
	0xa9, 0x76, 0xb9, 0x07, 0x9a, 0x33, 0x4d, 0x05, 0xd3, 0x8c, 0x2a, 0x8a, 0x6f, 0x92, 0x0a, 0xda,
	0x64, 0x91, 0xef, 0x71, 0x27, 0x73, 0xce, 0x59, 0xf8, 0x56, 0x16, 0x18, 0xda, 0xbe, 0xea, 0x68,
	0x97, 0x83, 0xc7, 0x70, 0x99, 0x3c, 0x45, 0x8f, 0x02, 0x08, 0x65, 0x14, 0x38, 0xa0, 0x03, 0x70,
	0x21, 0x00, 0xe1, 0x00, 0xd3, 0x5c, 0xe8, 0x76, 0xe4, 0x29, 0xee, 0x7b, 0xa0, 0xa5, 0x1f, 0xe2,
	0x15, 0xf2, 0x1c, 0x3d, 0xc9, 0x35, 0x3a, 0x37, 0x8f, 0x7c, 0x5f, 0x06, 0x0a, 0x98, 0x76, 0x65,
	0x30, 0x8b, 0x10, 0x46, 0x8d, 0x3c, 0x26, 0x5e, 0x25, 0x77, 0xd1, 0x0e, 0x17, 0x31, 0xf5, 0x38,
	0x5b, 0x8c, 0x87, 0x6f, 0x93, 0x03, 0xb4, 0xbb, 0x08, 0x15, 0x32, 0x44, 0x0e, 0xd1, 0xde, 0x5f,
	0x64, 0x39, 0xb8, 0x46, 0x76, 0xd0, 0x56, 0x01, 0xf2, 0xf6, 0xac, 0xbd, 0x4e, 0xf6, 0xd1, 0xf6,
	0xbc, 0x5d, 0xc4, 0xcf, 0x9c, 0xee, 0x90, 0x3d, 0x54, 0x99, 0x23, 0x1d, 0x19, 0x69, 0x15, 0x35,
	0x72, 0x60, 0x83, 0x3c, 0x44, 0xf7, 0x7e, 0x19, 0xef, 0x52, 0xee, 0x01, 0xd3, 0x2a, 0xa0, 0x22,
	0x74, 0x24, 0xe3, 0xa2, 0x89, 0x37, 0xb3, 0xa9, 0x86, 0x92, 0x2d, 0x3b, 0x87, 0x80, 0x61, 0x4c,
	0x5e, 0xa0, 0x13, 0x83, 0x64, 0xfe, 0x9a, 0x49, 0x08, 0x73, 0x4e, 0x9b, 0x2a, 0xa7, 0xf5, 0x47,
	0xea, 0xad, 0x2c, 0x84, 0x90, 0xe6, 0xb1, 0x43, 0x7d, 0xfd, 0xa1, 0x10, 0x52, 0x45, 0xc7, 0x22,
	0xf2, 0xbc, 0x85, 0xd0, 0xf9, 0xbf, 0x67, 0x66, 0x70, 0x86, 0x2b, 0x64, 0x1d, 0xad, 0x2a, 0x29,
	0xb5, 0x27, 0x45, 0x13, 0x6f, 0x93, 0x6d, 0x84, 0x23, 0x71, 0xfd, 0xf0, 0xb9, 0xc1, 0x0e, 0x79,
	0x80, 0x8e, 0x8a, 0xed, 0x42, 0x08, 0x62, 0xee, 0x80, 0x8e, 0x04, 0x8d, 0x29, 0xf7, 0x68, 0xc3,
	0x03, 0xbc, 0x4b, 0x4e, 0xd1, 0xb3, 0x82, 0x11, 0x73, 0x06, 0x52, 0xb7, 0x68, 0x16, 0x5a, 0x68,
	0x5f, 0x86, 0x5c, 0xf1, 0x18, 0x34, 0x8b, 0x02, 0xaa, 0xb8, 0x14, 0x78, 0x8f, 0x1c, 0xa1, 0xfd,
	0xdf, 0xd8, 0xd9, 0x7a, 0xae, 0x8c, 0x04, 0xc3, 0xfb, 0x8d, 0x1f, 0x4b, 0xa8, 0xfa, 0x3e, 0x39,
	0xb7, 0xfe, 0x7f, 0x5e, 0x8d, 0xca, 0xe2, 0x99, 0xf8, 0xd9, 0x55, 0xf9, 0x4b, 0x6f, 0xd9, 0x4c,
	0x36, 0x4c, 0xc6, 0xbd, 0xc9, 0xd0, 0x4a, 0xa6, 0xc3, 0xda, 0x70, 0x30, 0xc9, 0x6f, 0x6e, 0x7e,
	0xdc, 0x17, 0xa3, 0xf4, 0x5f, 0xb7, 0xfe, 0xda, 0xfc, 0x7c, 0x29, 0x2d, 0x37, 0x29, 0xfd, 0x5a,
	0x3a, 0x6e, 0x9a, 0x61, 0xb4, 0x9f, 0x5a, 0xa6, 0xcc, 0xaa, 0xb8, 0x6e, 0xe5, 0x96, 0xe9, 0xb7,
	0x39, 0xa1, 0x4b, 0xfb, 0x69, 0xb7, 0x20, 0x74, 0xe3, 0x7a, 0xd7, 0x10, 0xbe, 0x97, 0xaa, 0xa6,
	0x6b, 0xdb, 0xb4, 0x9f, 0xda, 0x76, 0x41, 0xb1, 0xed, 0xb8, 0x6e, 0xdb, 0x86, 0xf4, 0xae, 0x9c,
	0xa7, 0xab, 0xff, 0x0c, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x9e, 0xce, 0x44, 0x88, 0x04, 0x00, 0x00,
}
