// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.24.3
// source: node/tag/tag.proto

package tag

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

type TagType int32

const (
	TagType_node   TagType = 0
	TagType_mirror TagType = 1
)

// Enum value maps for TagType.
var (
	TagType_name = map[int32]string{
		0: "node",
		1: "mirror",
	}
	TagType_value = map[string]int32{
		"node":   0,
		"mirror": 1,
	}
)

func (x TagType) Enum() *TagType {
	p := new(TagType)
	*p = x
	return p
}

func (x TagType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TagType) Descriptor() protoreflect.EnumDescriptor {
	return file_node_tag_tag_proto_enumTypes[0].Descriptor()
}

func (TagType) Type() protoreflect.EnumType {
	return &file_node_tag_tag_proto_enumTypes[0]
}

func (x TagType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TagType.Descriptor instead.
func (TagType) EnumDescriptor() ([]byte, []int) {
	return file_node_tag_tag_proto_rawDescGZIP(), []int{0}
}

type Tags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Type TagType  `protobuf:"varint,3,opt,name=type,proto3,enum=yuhaiin.tag.TagType" json:"type,omitempty"`
	Hash []string `protobuf:"bytes,2,rep,name=hash,proto3" json:"hash,omitempty"`
}

func (x *Tags) Reset() {
	*x = Tags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_tag_tag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tags) ProtoMessage() {}

func (x *Tags) ProtoReflect() protoreflect.Message {
	mi := &file_node_tag_tag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tags.ProtoReflect.Descriptor instead.
func (*Tags) Descriptor() ([]byte, []int) {
	return file_node_tag_tag_proto_rawDescGZIP(), []int{0}
}

func (x *Tags) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Tags) GetType() TagType {
	if x != nil {
		return x.Type
	}
	return TagType_node
}

func (x *Tags) GetHash() []string {
	if x != nil {
		return x.Hash
	}
	return nil
}

var File_node_tag_tag_proto protoreflect.FileDescriptor

var file_node_tag_tag_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x74, 0x61, 0x67, 0x2f, 0x74, 0x61, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69, 0x6e, 0x2e, 0x74, 0x61,
	0x67, 0x22, 0x57, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x29, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x79, 0x75, 0x68, 0x61,
	0x69, 0x69, 0x6e, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x74, 0x61, 0x67, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x2a, 0x20, 0x0a, 0x08, 0x74, 0x61,
	0x67, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x42, 0x32, 0x5a, 0x30,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x73, 0x75, 0x74, 0x6f,
	0x72, 0x75, 0x66, 0x61, 0x2f, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69, 0x6e, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x74, 0x61, 0x67,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_tag_tag_proto_rawDescOnce sync.Once
	file_node_tag_tag_proto_rawDescData = file_node_tag_tag_proto_rawDesc
)

func file_node_tag_tag_proto_rawDescGZIP() []byte {
	file_node_tag_tag_proto_rawDescOnce.Do(func() {
		file_node_tag_tag_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_tag_tag_proto_rawDescData)
	})
	return file_node_tag_tag_proto_rawDescData
}

var file_node_tag_tag_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_node_tag_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_node_tag_tag_proto_goTypes = []interface{}{
	(TagType)(0), // 0: yuhaiin.tag.tag_type
	(*Tags)(nil), // 1: yuhaiin.tag.tags
}
var file_node_tag_tag_proto_depIdxs = []int32{
	0, // 0: yuhaiin.tag.tags.type:type_name -> yuhaiin.tag.tag_type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_node_tag_tag_proto_init() }
func file_node_tag_tag_proto_init() {
	if File_node_tag_tag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_tag_tag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tags); i {
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
			RawDescriptor: file_node_tag_tag_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_tag_tag_proto_goTypes,
		DependencyIndexes: file_node_tag_tag_proto_depIdxs,
		EnumInfos:         file_node_tag_tag_proto_enumTypes,
		MessageInfos:      file_node_tag_tag_proto_msgTypes,
	}.Build()
	File_node_tag_tag_proto = out.File
	file_node_tag_tag_proto_rawDesc = nil
	file_node_tag_tag_proto_goTypes = nil
	file_node_tag_tag_proto_depIdxs = nil
}
