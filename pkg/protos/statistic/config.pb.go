// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: pkg/protos/statistic/config.proto

package statistic

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

type CloseConnsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conns []int64 `protobuf:"varint,1,rep,packed,name=conns,proto3" json:"conns,omitempty"`
}

func (x *CloseConnsReq) Reset() {
	*x = CloseConnsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_statistic_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseConnsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseConnsReq) ProtoMessage() {}

func (x *CloseConnsReq) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_statistic_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseConnsReq.ProtoReflect.Descriptor instead.
func (*CloseConnsReq) Descriptor() ([]byte, []int) {
	return file_pkg_protos_statistic_config_proto_rawDescGZIP(), []int{0}
}

func (x *CloseConnsReq) GetConns() []int64 {
	if x != nil {
		return x.Conns
	}
	return nil
}

type RateResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Download     string `protobuf:"bytes,1,opt,name=download,proto3" json:"download,omitempty"`
	Upload       string `protobuf:"bytes,2,opt,name=upload,proto3" json:"upload,omitempty"`
	DownloadRate string `protobuf:"bytes,3,opt,name=download_rate,json=downloadRate,proto3" json:"download_rate,omitempty"`
	UploadRate   string `protobuf:"bytes,4,opt,name=upload_rate,json=uploadRate,proto3" json:"upload_rate,omitempty"`
}

func (x *RateResp) Reset() {
	*x = RateResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_statistic_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateResp) ProtoMessage() {}

func (x *RateResp) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_statistic_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateResp.ProtoReflect.Descriptor instead.
func (*RateResp) Descriptor() ([]byte, []int) {
	return file_pkg_protos_statistic_config_proto_rawDescGZIP(), []int{1}
}

func (x *RateResp) GetDownload() string {
	if x != nil {
		return x.Download
	}
	return ""
}

func (x *RateResp) GetUpload() string {
	if x != nil {
		return x.Upload
	}
	return ""
}

func (x *RateResp) GetDownloadRate() string {
	if x != nil {
		return x.DownloadRate
	}
	return ""
}

func (x *RateResp) GetUploadRate() string {
	if x != nil {
		return x.UploadRate
	}
	return ""
}

type ConnResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Connections []*ConnRespConnection `protobuf:"bytes,1,rep,name=connections,proto3" json:"connections,omitempty"`
}

func (x *ConnResp) Reset() {
	*x = ConnResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_statistic_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnResp) ProtoMessage() {}

func (x *ConnResp) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_statistic_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnResp.ProtoReflect.Descriptor instead.
func (*ConnResp) Descriptor() ([]byte, []int) {
	return file_pkg_protos_statistic_config_proto_rawDescGZIP(), []int{2}
}

func (x *ConnResp) GetConnections() []*ConnRespConnection {
	if x != nil {
		return x.Connections
	}
	return nil
}

type ConnRespConnection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr   string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Id     int64  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Local  string `protobuf:"bytes,3,opt,name=local,proto3" json:"local,omitempty"`
	Remote string `protobuf:"bytes,4,opt,name=remote,proto3" json:"remote,omitempty"`
}

func (x *ConnRespConnection) Reset() {
	*x = ConnRespConnection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_statistic_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnRespConnection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnRespConnection) ProtoMessage() {}

func (x *ConnRespConnection) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_statistic_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnRespConnection.ProtoReflect.Descriptor instead.
func (*ConnRespConnection) Descriptor() ([]byte, []int) {
	return file_pkg_protos_statistic_config_proto_rawDescGZIP(), []int{2, 0}
}

func (x *ConnRespConnection) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *ConnRespConnection) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ConnRespConnection) GetLocal() string {
	if x != nil {
		return x.Local
	}
	return ""
}

func (x *ConnRespConnection) GetRemote() string {
	if x != nil {
		return x.Remote
	}
	return ""
}

var File_pkg_protos_statistic_config_proto protoreflect.FileDescriptor

var file_pkg_protos_statistic_config_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x11, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69, 0x6e, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x27, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x6e, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6e, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x6e, 0x6e, 0x73, 0x22, 0x85, 0x01, 0x0a,
	0x09, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x23,
	0x0a, 0x0d, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52,
	0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x72, 0x61,
	0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x61, 0x74, 0x65, 0x22, 0xb6, 0x01, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x12, 0x49, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69,
	0x6e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x63, 0x6f, 0x6e, 0x6e,
	0x5f, 0x72, 0x65, 0x73, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x5e, 0x0a,
	0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x32, 0xdb, 0x01,
	0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3d, 0x0a,
	0x05, 0x63, 0x6f, 0x6e, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c,
	0x2e, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69, 0x6e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x12, 0x48, 0x0a, 0x0a,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x12, 0x22, 0x2e, 0x79, 0x75, 0x68,
	0x61, 0x69, 0x69, 0x6e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73,
	0x74, 0x69, 0x63, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c, 0x2e, 0x79, 0x75,
	0x68, 0x61, 0x69, 0x69, 0x6e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x2e,
	0x72, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x30, 0x01, 0x42, 0x33, 0x5a, 0x31, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x73, 0x75, 0x74, 0x6f, 0x72,
	0x75, 0x66, 0x61, 0x2f, 0x79, 0x75, 0x68, 0x61, 0x69, 0x69, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_protos_statistic_config_proto_rawDescOnce sync.Once
	file_pkg_protos_statistic_config_proto_rawDescData = file_pkg_protos_statistic_config_proto_rawDesc
)

func file_pkg_protos_statistic_config_proto_rawDescGZIP() []byte {
	file_pkg_protos_statistic_config_proto_rawDescOnce.Do(func() {
		file_pkg_protos_statistic_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protos_statistic_config_proto_rawDescData)
	})
	return file_pkg_protos_statistic_config_proto_rawDescData
}

var file_pkg_protos_statistic_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_protos_statistic_config_proto_goTypes = []interface{}{
	(*CloseConnsReq)(nil),      // 0: yuhaiin.statistic.close_conns_req
	(*RateResp)(nil),           // 1: yuhaiin.statistic.rate_resp
	(*ConnResp)(nil),           // 2: yuhaiin.statistic.conn_resp
	(*ConnRespConnection)(nil), // 3: yuhaiin.statistic.conn_resp.connection
	(*emptypb.Empty)(nil),      // 4: google.protobuf.Empty
}
var file_pkg_protos_statistic_config_proto_depIdxs = []int32{
	3, // 0: yuhaiin.statistic.conn_resp.connections:type_name -> yuhaiin.statistic.conn_resp.connection
	4, // 1: yuhaiin.statistic.connections.conns:input_type -> google.protobuf.Empty
	0, // 2: yuhaiin.statistic.connections.close_conn:input_type -> yuhaiin.statistic.close_conns_req
	4, // 3: yuhaiin.statistic.connections.statistic:input_type -> google.protobuf.Empty
	2, // 4: yuhaiin.statistic.connections.conns:output_type -> yuhaiin.statistic.conn_resp
	4, // 5: yuhaiin.statistic.connections.close_conn:output_type -> google.protobuf.Empty
	1, // 6: yuhaiin.statistic.connections.statistic:output_type -> yuhaiin.statistic.rate_resp
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_protos_statistic_config_proto_init() }
func file_pkg_protos_statistic_config_proto_init() {
	if File_pkg_protos_statistic_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_protos_statistic_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloseConnsReq); i {
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
		file_pkg_protos_statistic_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateResp); i {
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
		file_pkg_protos_statistic_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnResp); i {
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
		file_pkg_protos_statistic_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnRespConnection); i {
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
			RawDescriptor: file_pkg_protos_statistic_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_protos_statistic_config_proto_goTypes,
		DependencyIndexes: file_pkg_protos_statistic_config_proto_depIdxs,
		MessageInfos:      file_pkg_protos_statistic_config_proto_msgTypes,
	}.Build()
	File_pkg_protos_statistic_config_proto = out.File
	file_pkg_protos_statistic_config_proto_rawDesc = nil
	file_pkg_protos_statistic_config_proto_goTypes = nil
	file_pkg_protos_statistic_config_proto_depIdxs = nil
}