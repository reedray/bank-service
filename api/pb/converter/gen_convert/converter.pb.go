// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: converter.proto

package gen_convert

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

type Money struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount       string `protobuf:"bytes,1,opt,name=amount,proto3" json:"amount,omitempty"`
	CurrencyCode string `protobuf:"bytes,2,opt,name=currencyCode,proto3" json:"currencyCode,omitempty"`
}

func (x *Money) Reset() {
	*x = Money{}
	if protoimpl.UnsafeEnabled {
		mi := &file_converter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Money) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Money) ProtoMessage() {}

func (x *Money) ProtoReflect() protoreflect.Message {
	mi := &file_converter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Money.ProtoReflect.Descriptor instead.
func (*Money) Descriptor() ([]byte, []int) {
	return file_converter_proto_rawDescGZIP(), []int{0}
}

func (x *Money) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Money) GetCurrencyCode() string {
	if x != nil {
		return x.CurrencyCode
	}
	return ""
}

var File_converter_proto protoreflect.FileDescriptor

var file_converter_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x43, 0x0a, 0x05, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x32, 0x51, 0x0a, 0x0e, 0x43, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x12, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x6f, 0x6e,
	0x65, 0x79, 0x1a, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x42, 0x0f, 0x5a,
	0x0d, 0x2e, 0x2f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_converter_proto_rawDescOnce sync.Once
	file_converter_proto_rawDescData = file_converter_proto_rawDesc
)

func file_converter_proto_rawDescGZIP() []byte {
	file_converter_proto_rawDescOnce.Do(func() {
		file_converter_proto_rawDescData = protoimpl.X.CompressGZIP(file_converter_proto_rawDescData)
	})
	return file_converter_proto_rawDescData
}

var file_converter_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_converter_proto_goTypes = []interface{}{
	(*Money)(nil), // 0: converter.protobuf.Money
}
var file_converter_proto_depIdxs = []int32{
	0, // 0: converter.protobuf.ConvertService.Convert:input_type -> converter.protobuf.Money
	0, // 1: converter.protobuf.ConvertService.Convert:output_type -> converter.protobuf.Money
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_converter_proto_init() }
func file_converter_proto_init() {
	if File_converter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_converter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Money); i {
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
			RawDescriptor: file_converter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_converter_proto_goTypes,
		DependencyIndexes: file_converter_proto_depIdxs,
		MessageInfos:      file_converter_proto_msgTypes,
	}.Build()
	File_converter_proto = out.File
	file_converter_proto_rawDesc = nil
	file_converter_proto_goTypes = nil
	file_converter_proto_depIdxs = nil
}
