// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: data.proto

package stream

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PeerData_OS int32

const (
	PeerData_LINUX   PeerData_OS = 0
	PeerData_DARWIN  PeerData_OS = 1
	PeerData_WINDOWS PeerData_OS = 2
	PeerData_ANDROID PeerData_OS = 3
	PeerData_IOS     PeerData_OS = 4
	PeerData_OTHER   PeerData_OS = 5
)

// Enum value maps for PeerData_OS.
var (
	PeerData_OS_name = map[int32]string{
		0: "LINUX",
		1: "DARWIN",
		2: "WINDOWS",
		3: "ANDROID",
		4: "IOS",
		5: "OTHER",
	}
	PeerData_OS_value = map[string]int32{
		"LINUX":   0,
		"DARWIN":  1,
		"WINDOWS": 2,
		"ANDROID": 3,
		"IOS":     4,
		"OTHER":   5,
	}
)

func (x PeerData_OS) Enum() *PeerData_OS {
	p := new(PeerData_OS)
	*p = x
	return p
}

func (x PeerData_OS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PeerData_OS) Descriptor() protoreflect.EnumDescriptor {
	return file_data_proto_enumTypes[0].Descriptor()
}

func (PeerData_OS) Type() protoreflect.EnumType {
	return &file_data_proto_enumTypes[0]
}

func (x PeerData_OS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PeerData_OS.Descriptor instead.
func (PeerData_OS) EnumDescriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{0, 0}
}

type PeerData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Os        PeerData_OS `protobuf:"varint,1,opt,name=os,proto3,enum=stream.PeerData_OS" json:"os,omitempty"`
	PublicKey []byte      `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *PeerData) Reset() {
	*x = PeerData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerData) ProtoMessage() {}

func (x *PeerData) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerData.ProtoReflect.Descriptor instead.
func (*PeerData) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{0}
}

func (x *PeerData) GetOs() PeerData_OS {
	if x != nil {
		return x.Os
	}
	return PeerData_LINUX
}

func (x *PeerData) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

type ClipboardData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsImage  bool   `protobuf:"varint,1,opt,name=is_image,json=isImage,proto3" json:"is_image,omitempty"`
	DataSize uint32 `protobuf:"fixed32,2,opt,name=data_size,json=dataSize,proto3" json:"data_size,omitempty"`
	Time     int64  `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Data     []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ClipboardData) Reset() {
	*x = ClipboardData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClipboardData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClipboardData) ProtoMessage() {}

func (x *ClipboardData) ProtoReflect() protoreflect.Message {
	mi := &file_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClipboardData.ProtoReflect.Descriptor instead.
func (*ClipboardData) Descriptor() ([]byte, []int) {
	return file_data_proto_rawDescGZIP(), []int{1}
}

func (x *ClipboardData) GetIsImage() bool {
	if x != nil {
		return x.IsImage
	}
	return false
}

func (x *ClipboardData) GetDataSize() uint32 {
	if x != nil {
		return x.DataSize
	}
	return 0
}

func (x *ClipboardData) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *ClipboardData) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_data_proto protoreflect.FileDescriptor

var file_data_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x22, 0x99, 0x01, 0x0a, 0x08, 0x50, 0x65, 0x65, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x23, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x2e,
	0x4f, 0x53, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x49, 0x0a, 0x02, 0x4f, 0x53, 0x12, 0x09, 0x0a, 0x05, 0x4c,
	0x49, 0x4e, 0x55, 0x58, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x41, 0x52, 0x57, 0x49, 0x4e,
	0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x49, 0x4e, 0x44, 0x4f, 0x57, 0x53, 0x10, 0x02, 0x12,
	0x0b, 0x0a, 0x07, 0x41, 0x4e, 0x44, 0x52, 0x4f, 0x49, 0x44, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03,
	0x49, 0x4f, 0x53, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x05,
	0x22, 0x6f, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x70, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x07, 0x52,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6e, 0x74, 0x73, 0x64, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x63, 0x6c, 0x69, 0x70, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_proto_rawDescOnce sync.Once
	file_data_proto_rawDescData = file_data_proto_rawDesc
)

func file_data_proto_rawDescGZIP() []byte {
	file_data_proto_rawDescOnce.Do(func() {
		file_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_proto_rawDescData)
	})
	return file_data_proto_rawDescData
}

var file_data_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_data_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_data_proto_goTypes = []interface{}{
	(PeerData_OS)(0),      // 0: stream.PeerData.OS
	(*PeerData)(nil),      // 1: stream.PeerData
	(*ClipboardData)(nil), // 2: stream.ClipboardData
}
var file_data_proto_depIdxs = []int32{
	0, // 0: stream.PeerData.os:type_name -> stream.PeerData.OS
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_data_proto_init() }
func file_data_proto_init() {
	if File_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerData); i {
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
		file_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClipboardData); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_proto_goTypes,
		DependencyIndexes: file_data_proto_depIdxs,
		EnumInfos:         file_data_proto_enumTypes,
		MessageInfos:      file_data_proto_msgTypes,
	}.Build()
	File_data_proto = out.File
	file_data_proto_rawDesc = nil
	file_data_proto_goTypes = nil
	file_data_proto_depIdxs = nil
}
