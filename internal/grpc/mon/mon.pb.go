// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: proto/mon.proto

package mon

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricKey_MetricType int32

const (
	MetricKey_GAUGE   MetricKey_MetricType = 0
	MetricKey_COUNTER MetricKey_MetricType = 1
)

// Enum value maps for MetricKey_MetricType.
var (
	MetricKey_MetricType_name = map[int32]string{
		0: "GAUGE",
		1: "COUNTER",
	}
	MetricKey_MetricType_value = map[string]int32{
		"GAUGE":   0,
		"COUNTER": 1,
	}
)

func (x MetricKey_MetricType) Enum() *MetricKey_MetricType {
	p := new(MetricKey_MetricType)
	*p = x
	return p
}

func (x MetricKey_MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricKey_MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_mon_proto_enumTypes[0].Descriptor()
}

func (MetricKey_MetricType) Type() protoreflect.EnumType {
	return &file_proto_mon_proto_enumTypes[0]
}

func (x MetricKey_MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricKey_MetricType.Descriptor instead.
func (MetricKey_MetricType) EnumDescriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{0, 0}
}

type MetricKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string               `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Type MetricKey_MetricType `protobuf:"varint,2,opt,name=type,proto3,enum=mon.MetricKey_MetricType" json:"type,omitempty"`
}

func (x *MetricKey) Reset() {
	*x = MetricKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricKey) ProtoMessage() {}

func (x *MetricKey) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricKey.ProtoReflect.Descriptor instead.
func (*MetricKey) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{0}
}

func (x *MetricKey) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MetricKey) GetType() MetricKey_MetricType {
	if x != nil {
		return x.Type
	}
	return MetricKey_GAUGE
}

type MetricValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetricId string `protobuf:"bytes,1,opt,name=metricId,proto3" json:"metricId,omitempty"`
	// Types that are assignable to Value:
	//
	//	*MetricValue_Gauge
	//	*MetricValue_Counter
	Value isMetricValue_Value `protobuf_oneof:"value"`
}

func (x *MetricValue) Reset() {
	*x = MetricValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricValue) ProtoMessage() {}

func (x *MetricValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricValue.ProtoReflect.Descriptor instead.
func (*MetricValue) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{1}
}

func (x *MetricValue) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (m *MetricValue) GetValue() isMetricValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *MetricValue) GetGauge() float64 {
	if x, ok := x.GetValue().(*MetricValue_Gauge); ok {
		return x.Gauge
	}
	return 0
}

func (x *MetricValue) GetCounter() int64 {
	if x, ok := x.GetValue().(*MetricValue_Counter); ok {
		return x.Counter
	}
	return 0
}

type isMetricValue_Value interface {
	isMetricValue_Value()
}

type MetricValue_Gauge struct {
	Gauge float64 `protobuf:"fixed64,2,opt,name=gauge,proto3,oneof"`
}

type MetricValue_Counter struct {
	Counter int64 `protobuf:"varint,3,opt,name=counter,proto3,oneof"`
}

func (*MetricValue_Gauge) isMetricValue_Value() {}

func (*MetricValue_Counter) isMetricValue_Value() {}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key *MetricKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequest.ProtoReflect.Descriptor instead.
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{2}
}

func (x *GetMetricRequest) GetKey() *MetricKey {
	if x != nil {
		return x.Key
	}
	return nil
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *MetricValue `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponse.ProtoReflect.Descriptor instead.
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricResponse) GetMetric() *MetricValue {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *MetricValue `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricRequest) GetMetric() *MetricValue {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *MetricValue `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricResponse) Reset() {
	*x = UpdateMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricResponse) ProtoMessage() {}

func (x *UpdateMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMetricResponse) GetMetric() *MetricValue {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateMetricsBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*MetricValue `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricsBatchRequest) Reset() {
	*x = UpdateMetricsBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mon_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsBatchRequest) ProtoMessage() {}

func (x *UpdateMetricsBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mon_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsBatchRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsBatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_mon_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateMetricsBatchRequest) GetMetrics() []*MetricValue {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_proto_mon_proto protoreflect.FileDescriptor

var file_proto_mon_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x6d, 0x6f, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64,
	0x12, 0x2d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19,
	0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a,
	0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x45, 0x52, 0x10, 0x01, 0x22, 0x66, 0x0a, 0x0b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x48,
	0x00, 0x52, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x34, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x20, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x3d, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x22, 0x3f, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x06, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x22, 0x40, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f,
	0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x47, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x32, 0x99,
	0x02, 0x0a, 0x0e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3a, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x15,
	0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a,
	0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x18, 0x2e,
	0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4c, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1e, 0x2e, 0x6d, 0x6f, 0x6e, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x38, 0x0a, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x44, 0x42, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x13, 0x5a, 0x11, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x6f, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_mon_proto_rawDescOnce sync.Once
	file_proto_mon_proto_rawDescData = file_proto_mon_proto_rawDesc
)

func file_proto_mon_proto_rawDescGZIP() []byte {
	file_proto_mon_proto_rawDescOnce.Do(func() {
		file_proto_mon_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_mon_proto_rawDescData)
	})
	return file_proto_mon_proto_rawDescData
}

var file_proto_mon_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_mon_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_mon_proto_goTypes = []any{
	(MetricKey_MetricType)(0),         // 0: mon.MetricKey.MetricType
	(*MetricKey)(nil),                 // 1: mon.MetricKey
	(*MetricValue)(nil),               // 2: mon.MetricValue
	(*GetMetricRequest)(nil),          // 3: mon.GetMetricRequest
	(*GetMetricResponse)(nil),         // 4: mon.GetMetricResponse
	(*UpdateMetricRequest)(nil),       // 5: mon.UpdateMetricRequest
	(*UpdateMetricResponse)(nil),      // 6: mon.UpdateMetricResponse
	(*UpdateMetricsBatchRequest)(nil), // 7: mon.UpdateMetricsBatchRequest
	(*emptypb.Empty)(nil),             // 8: google.protobuf.Empty
}
var file_proto_mon_proto_depIdxs = []int32{
	0,  // 0: mon.MetricKey.type:type_name -> mon.MetricKey.MetricType
	1,  // 1: mon.GetMetricRequest.key:type_name -> mon.MetricKey
	2,  // 2: mon.GetMetricResponse.metric:type_name -> mon.MetricValue
	2,  // 3: mon.UpdateMetricRequest.metric:type_name -> mon.MetricValue
	2,  // 4: mon.UpdateMetricResponse.metric:type_name -> mon.MetricValue
	2,  // 5: mon.UpdateMetricsBatchRequest.metrics:type_name -> mon.MetricValue
	3,  // 6: mon.MonitorService.GetMetric:input_type -> mon.GetMetricRequest
	5,  // 7: mon.MonitorService.UpdateMetric:input_type -> mon.UpdateMetricRequest
	7,  // 8: mon.MonitorService.UpdateMetricsBatch:input_type -> mon.UpdateMetricsBatchRequest
	8,  // 9: mon.MonitorService.PingDB:input_type -> google.protobuf.Empty
	4,  // 10: mon.MonitorService.GetMetric:output_type -> mon.GetMetricResponse
	6,  // 11: mon.MonitorService.UpdateMetric:output_type -> mon.UpdateMetricResponse
	8,  // 12: mon.MonitorService.UpdateMetricsBatch:output_type -> google.protobuf.Empty
	8,  // 13: mon.MonitorService.PingDB:output_type -> google.protobuf.Empty
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_proto_mon_proto_init() }
func file_proto_mon_proto_init() {
	if File_proto_mon_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_mon_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MetricKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MetricValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_mon_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricsBatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_mon_proto_msgTypes[1].OneofWrappers = []any{
		(*MetricValue_Gauge)(nil),
		(*MetricValue_Counter)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_mon_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_mon_proto_goTypes,
		DependencyIndexes: file_proto_mon_proto_depIdxs,
		EnumInfos:         file_proto_mon_proto_enumTypes,
		MessageInfos:      file_proto_mon_proto_msgTypes,
	}.Build()
	File_proto_mon_proto = out.File
	file_proto_mon_proto_rawDesc = nil
	file_proto_mon_proto_goTypes = nil
	file_proto_mon_proto_depIdxs = nil
}