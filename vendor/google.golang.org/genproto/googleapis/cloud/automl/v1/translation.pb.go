// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/automl/v1/translation.proto

package automl

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

// Dataset metadata that is specific to translation.
type TranslationDatasetMetadata struct {
	// Required. The BCP-47 language code of the source language.
	SourceLanguageCode string `protobuf:"bytes,1,opt,name=source_language_code,json=sourceLanguageCode,proto3" json:"source_language_code,omitempty"`
	// Required. The BCP-47 language code of the target language.
	TargetLanguageCode   string   `protobuf:"bytes,2,opt,name=target_language_code,json=targetLanguageCode,proto3" json:"target_language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationDatasetMetadata) Reset()         { *m = TranslationDatasetMetadata{} }
func (m *TranslationDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*TranslationDatasetMetadata) ProtoMessage()    {}
func (*TranslationDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3df1fd6bde1101e, []int{0}
}

func (m *TranslationDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationDatasetMetadata.Unmarshal(m, b)
}
func (m *TranslationDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *TranslationDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationDatasetMetadata.Merge(m, src)
}
func (m *TranslationDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_TranslationDatasetMetadata.Size(m)
}
func (m *TranslationDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationDatasetMetadata proto.InternalMessageInfo

func (m *TranslationDatasetMetadata) GetSourceLanguageCode() string {
	if m != nil {
		return m.SourceLanguageCode
	}
	return ""
}

func (m *TranslationDatasetMetadata) GetTargetLanguageCode() string {
	if m != nil {
		return m.TargetLanguageCode
	}
	return ""
}

// Evaluation metrics for the dataset.
type TranslationEvaluationMetrics struct {
	// Output only. BLEU score.
	BleuScore float64 `protobuf:"fixed64,1,opt,name=bleu_score,json=bleuScore,proto3" json:"bleu_score,omitempty"`
	// Output only. BLEU score for base model.
	BaseBleuScore        float64  `protobuf:"fixed64,2,opt,name=base_bleu_score,json=baseBleuScore,proto3" json:"base_bleu_score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationEvaluationMetrics) Reset()         { *m = TranslationEvaluationMetrics{} }
func (m *TranslationEvaluationMetrics) String() string { return proto.CompactTextString(m) }
func (*TranslationEvaluationMetrics) ProtoMessage()    {}
func (*TranslationEvaluationMetrics) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3df1fd6bde1101e, []int{1}
}

func (m *TranslationEvaluationMetrics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationEvaluationMetrics.Unmarshal(m, b)
}
func (m *TranslationEvaluationMetrics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationEvaluationMetrics.Marshal(b, m, deterministic)
}
func (m *TranslationEvaluationMetrics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationEvaluationMetrics.Merge(m, src)
}
func (m *TranslationEvaluationMetrics) XXX_Size() int {
	return xxx_messageInfo_TranslationEvaluationMetrics.Size(m)
}
func (m *TranslationEvaluationMetrics) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationEvaluationMetrics.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationEvaluationMetrics proto.InternalMessageInfo

func (m *TranslationEvaluationMetrics) GetBleuScore() float64 {
	if m != nil {
		return m.BleuScore
	}
	return 0
}

func (m *TranslationEvaluationMetrics) GetBaseBleuScore() float64 {
	if m != nil {
		return m.BaseBleuScore
	}
	return 0
}

// Model metadata that is specific to translation.
type TranslationModelMetadata struct {
	// The resource name of the model to use as a baseline to train the custom
	// model. If unset, we use the default base model provided by Google
	// Translate. Format:
	// `projects/{project_id}/locations/{location_id}/models/{model_id}`
	BaseModel string `protobuf:"bytes,1,opt,name=base_model,json=baseModel,proto3" json:"base_model,omitempty"`
	// Output only. Inferred from the dataset.
	// The source language (The BCP-47 language code) that is used for training.
	SourceLanguageCode string `protobuf:"bytes,2,opt,name=source_language_code,json=sourceLanguageCode,proto3" json:"source_language_code,omitempty"`
	// Output only. The target language (The BCP-47 language code) that is used
	// for training.
	TargetLanguageCode   string   `protobuf:"bytes,3,opt,name=target_language_code,json=targetLanguageCode,proto3" json:"target_language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationModelMetadata) Reset()         { *m = TranslationModelMetadata{} }
func (m *TranslationModelMetadata) String() string { return proto.CompactTextString(m) }
func (*TranslationModelMetadata) ProtoMessage()    {}
func (*TranslationModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3df1fd6bde1101e, []int{2}
}

func (m *TranslationModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationModelMetadata.Unmarshal(m, b)
}
func (m *TranslationModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationModelMetadata.Marshal(b, m, deterministic)
}
func (m *TranslationModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationModelMetadata.Merge(m, src)
}
func (m *TranslationModelMetadata) XXX_Size() int {
	return xxx_messageInfo_TranslationModelMetadata.Size(m)
}
func (m *TranslationModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationModelMetadata proto.InternalMessageInfo

func (m *TranslationModelMetadata) GetBaseModel() string {
	if m != nil {
		return m.BaseModel
	}
	return ""
}

func (m *TranslationModelMetadata) GetSourceLanguageCode() string {
	if m != nil {
		return m.SourceLanguageCode
	}
	return ""
}

func (m *TranslationModelMetadata) GetTargetLanguageCode() string {
	if m != nil {
		return m.TargetLanguageCode
	}
	return ""
}

// Annotation details specific to translation.
type TranslationAnnotation struct {
	// Output only . The translated content.
	TranslatedContent    *TextSnippet `protobuf:"bytes,1,opt,name=translated_content,json=translatedContent,proto3" json:"translated_content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TranslationAnnotation) Reset()         { *m = TranslationAnnotation{} }
func (m *TranslationAnnotation) String() string { return proto.CompactTextString(m) }
func (*TranslationAnnotation) ProtoMessage()    {}
func (*TranslationAnnotation) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3df1fd6bde1101e, []int{3}
}

func (m *TranslationAnnotation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationAnnotation.Unmarshal(m, b)
}
func (m *TranslationAnnotation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationAnnotation.Marshal(b, m, deterministic)
}
func (m *TranslationAnnotation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationAnnotation.Merge(m, src)
}
func (m *TranslationAnnotation) XXX_Size() int {
	return xxx_messageInfo_TranslationAnnotation.Size(m)
}
func (m *TranslationAnnotation) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationAnnotation.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationAnnotation proto.InternalMessageInfo

func (m *TranslationAnnotation) GetTranslatedContent() *TextSnippet {
	if m != nil {
		return m.TranslatedContent
	}
	return nil
}

func init() {
	proto.RegisterType((*TranslationDatasetMetadata)(nil), "google.cloud.automl.v1.TranslationDatasetMetadata")
	proto.RegisterType((*TranslationEvaluationMetrics)(nil), "google.cloud.automl.v1.TranslationEvaluationMetrics")
	proto.RegisterType((*TranslationModelMetadata)(nil), "google.cloud.automl.v1.TranslationModelMetadata")
	proto.RegisterType((*TranslationAnnotation)(nil), "google.cloud.automl.v1.TranslationAnnotation")
}

func init() {
	proto.RegisterFile("google/cloud/automl/v1/translation.proto", fileDescriptor_c3df1fd6bde1101e)
}

var fileDescriptor_c3df1fd6bde1101e = []byte{
	// 454 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x51, 0x8a, 0xd4, 0x40,
	0x10, 0x25, 0x59, 0x10, 0xa6, 0x45, 0xd4, 0xa0, 0xcb, 0x3a, 0xac, 0x28, 0x23, 0xe8, 0x7e, 0x25,
	0x1b, 0xc5, 0x9f, 0xe8, 0xcf, 0xcc, 0x28, 0xfe, 0xec, 0xc0, 0x92, 0x5d, 0xe6, 0x43, 0x06, 0x42,
	0x4d, 0x52, 0xc6, 0x60, 0x4f, 0x57, 0x48, 0x2a, 0x83, 0x67, 0xf0, 0x0e, 0x5e, 0xc0, 0x33, 0x78,
	0x02, 0x4f, 0xe1, 0xb7, 0xa7, 0x90, 0xee, 0xce, 0xec, 0x36, 0xc3, 0xe4, 0xaf, 0xa9, 0xf7, 0x5e,
	0xbd, 0x57, 0xd5, 0x25, 0xce, 0x4a, 0xa2, 0x52, 0x62, 0x94, 0x4b, 0xea, 0x8a, 0x08, 0x3a, 0xa6,
	0x8d, 0x8c, 0xb6, 0x71, 0xc4, 0x0d, 0xa8, 0x56, 0x02, 0x57, 0xa4, 0xc2, 0xba, 0x21, 0xa6, 0xe0,
	0xd8, 0x32, 0x43, 0xc3, 0x0c, 0x2d, 0x33, 0xdc, 0xc6, 0xe3, 0x67, 0x7d, 0x07, 0xa8, 0xab, 0xe8,
	0x4b, 0x85, 0xb2, 0xc8, 0xd6, 0xf8, 0x15, 0xb6, 0x15, 0x35, 0x56, 0x38, 0x7e, 0x35, 0x60, 0x51,
	0x00, 0x43, 0x56, 0x31, 0x6e, 0xda, 0x9e, 0x78, 0xea, 0x74, 0x02, 0xa5, 0x88, 0x8d, 0x7d, 0x8f,
	0x4e, 0x7e, 0x78, 0x62, 0x7c, 0x7d, 0x9b, 0xea, 0x03, 0x30, 0xb4, 0xc8, 0x0b, 0x64, 0xd0, 0x8d,
	0x82, 0xb7, 0xe2, 0x51, 0x4b, 0x5d, 0x93, 0x63, 0x26, 0x41, 0x95, 0x1d, 0x94, 0x98, 0xe5, 0x54,
	0xe0, 0x89, 0xf7, 0xdc, 0x3b, 0x1b, 0xcd, 0x8e, 0xfe, 0x4e, 0xfd, 0x34, 0xb0, 0x84, 0x8b, 0x1e,
	0x9f, 0x53, 0x81, 0x5a, 0xc6, 0xd0, 0x94, 0xc8, 0x7b, 0x32, 0xdf, 0x91, 0x59, 0x82, 0x2b, 0x9b,
	0xa0, 0x38, 0x75, 0xb2, 0x7c, 0xdc, 0x82, 0xec, 0xcc, 0x6b, 0x81, 0xdc, 0x54, 0x79, 0x1b, 0x3c,
	0x15, 0x62, 0x2d, 0xb1, 0xcb, 0xda, 0x9c, 0x1a, 0x9b, 0xc1, 0x4b, 0x47, 0xba, 0x72, 0xa5, 0x0b,
	0xc1, 0x4b, 0x71, 0x7f, 0x0d, 0x2d, 0x66, 0x0e, 0xc7, 0x37, 0x9c, 0x7b, 0xba, 0x3c, 0xdb, 0xf1,
	0x26, 0x3f, 0x3d, 0x71, 0xe2, 0xf8, 0x2c, 0xa8, 0x40, 0x79, 0x33, 0xb1, 0xf6, 0xd0, 0x4d, 0x36,
	0xba, 0x6a, 0xe7, 0x4c, 0x47, 0xba, 0x62, 0x68, 0xc1, 0xf9, 0xc0, 0x42, 0xcc, 0x64, 0x07, 0x77,
	0x71, 0x3e, 0xb0, 0x8b, 0x23, 0xab, 0x38, 0xb0, 0x86, 0x6f, 0xe2, 0xb1, 0x13, 0x6f, 0x7a, 0xf3,
	0x67, 0x41, 0x2a, 0x82, 0xdd, 0x05, 0x61, 0x91, 0xe5, 0xa4, 0x18, 0x15, 0x9b, 0x8c, 0x77, 0x5f,
	0xbf, 0x08, 0x0f, 0x5f, 0x52, 0x78, 0x8d, 0xdf, 0xf9, 0x4a, 0x55, 0x75, 0x8d, 0x9c, 0x3e, 0xbc,
	0x95, 0xcf, 0xad, 0x7a, 0xf6, 0xdb, 0x13, 0xe3, 0x9c, 0x36, 0x03, 0xea, 0xd9, 0x03, 0x27, 0xc9,
	0xa5, 0xbe, 0x98, 0x4b, 0xef, 0xf3, 0xfb, 0x9e, 0x5b, 0x92, 0x9e, 0x28, 0xa4, 0xa6, 0x8c, 0x4a,
	0x54, 0xe6, 0x9e, 0x22, 0x0b, 0x41, 0x5d, 0xb5, 0xfb, 0x97, 0xf9, 0xce, 0xbe, 0x7e, 0xf9, 0xc7,
	0x9f, 0xac, 0x7c, 0x6e, 0xac, 0xa6, 0x1d, 0xd3, 0xe2, 0x22, 0x5c, 0xc6, 0x7f, 0x76, 0xc0, 0xca,
	0x00, 0x2b, 0x03, 0xc8, 0xd5, 0x32, 0xfe, 0xe7, 0x3f, 0xb1, 0x40, 0x92, 0x18, 0x24, 0x49, 0xac,
	0x26, 0x49, 0x96, 0xf1, 0xfa, 0x8e, 0xb1, 0x7d, 0xf3, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xaf, 0xa9,
	0x82, 0x1a, 0x72, 0x03, 0x00, 0x00,
}
