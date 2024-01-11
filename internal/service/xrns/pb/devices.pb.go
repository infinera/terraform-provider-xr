// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.3
// source: devices.proto

package pb

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

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_devices_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_devices_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_devices_proto_rawDescGZIP(), []int{0}
}

func (x *Device) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Device) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetDeviceByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetDeviceByNameRequest) Reset() {
	*x = GetDeviceByNameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_devices_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDeviceByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceByNameRequest) ProtoMessage() {}

func (x *GetDeviceByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_devices_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceByNameRequest.ProtoReflect.Descriptor instead.
func (*GetDeviceByNameRequest) Descriptor() ([]byte, []int) {
	return file_devices_proto_rawDescGZIP(), []int{1}
}

func (x *GetDeviceByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetAllDevicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllDevicesRequest) Reset() {
	*x = GetAllDevicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_devices_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllDevicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllDevicesRequest) ProtoMessage() {}

func (x *GetAllDevicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_devices_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllDevicesRequest.ProtoReflect.Descriptor instead.
func (*GetAllDevicesRequest) Descriptor() ([]byte, []int) {
	return file_devices_proto_rawDescGZIP(), []int{2}
}

type GetDeviceCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDeviceCountRequest) Reset() {
	*x = GetDeviceCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_devices_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDeviceCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceCountRequest) ProtoMessage() {}

func (x *GetDeviceCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_devices_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceCountRequest.ProtoReflect.Descriptor instead.
func (*GetDeviceCountRequest) Descriptor() ([]byte, []int) {
	return file_devices_proto_rawDescGZIP(), []int{3}
}

type GetDeviceCountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetDeviceCountResponse) Reset() {
	*x = GetDeviceCountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_devices_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDeviceCountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceCountResponse) ProtoMessage() {}

func (x *GetDeviceCountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_devices_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceCountResponse.ProtoReflect.Descriptor instead.
func (*GetDeviceCountResponse) Descriptor() ([]byte, []int) {
	return file_devices_proto_rawDescGZIP(), []int{4}
}

func (x *GetDeviceCountResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_devices_proto protoreflect.FileDescriptor

var file_devices_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x78, 0x72, 0x6e, 0x73, 0x22, 0x2c, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x2c, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x2e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x42, 0x2b, 0x5a, 0x29, 0x62, 0x69, 0x74, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x2e,
	0x69, 0x6e, 0x66, 0x69, 0x6e, 0x65, 0x72, 0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x72,
	0x76, 0x65, 0x6c, 0x2f, 0x69, 0x70, 0x6d, 0x2d, 0x78, 0x72, 0x6e, 0x73, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_devices_proto_rawDescOnce sync.Once
	file_devices_proto_rawDescData = file_devices_proto_rawDesc
)

func file_devices_proto_rawDescGZIP() []byte {
	file_devices_proto_rawDescOnce.Do(func() {
		file_devices_proto_rawDescData = protoimpl.X.CompressGZIP(file_devices_proto_rawDescData)
	})
	return file_devices_proto_rawDescData
}

var file_devices_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_devices_proto_goTypes = []interface{}{
	(*Device)(nil),                 // 0: xrns.Device
	(*GetDeviceByNameRequest)(nil), // 1: xrns.GetDeviceByNameRequest
	(*GetAllDevicesRequest)(nil),   // 2: xrns.GetAllDevicesRequest
	(*GetDeviceCountRequest)(nil),  // 3: xrns.GetDeviceCountRequest
	(*GetDeviceCountResponse)(nil), // 4: xrns.GetDeviceCountResponse
}
var file_devices_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_devices_proto_init() }
func file_devices_proto_init() {
	if File_devices_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_devices_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
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
		file_devices_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDeviceByNameRequest); i {
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
		file_devices_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllDevicesRequest); i {
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
		file_devices_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDeviceCountRequest); i {
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
		file_devices_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDeviceCountResponse); i {
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
			RawDescriptor: file_devices_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_devices_proto_goTypes,
		DependencyIndexes: file_devices_proto_depIdxs,
		MessageInfos:      file_devices_proto_msgTypes,
	}.Build()
	File_devices_proto = out.File
	file_devices_proto_rawDesc = nil
	file_devices_proto_goTypes = nil
	file_devices_proto_depIdxs = nil
}
