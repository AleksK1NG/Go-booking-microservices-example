// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: user.proto

//protoc --go_out=plugins=grpc:. *.proto

package userService

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID    string                 `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	FirstName string                 `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty"`
	LastName  string                 `protobuf:"bytes,3,opt,name=LastName,proto3" json:"LastName,omitempty"`
	Email     string                 `protobuf:"bytes,4,opt,name=Email,proto3" json:"Email,omitempty"`
	Avatar    string                 `protobuf:"bytes,5,opt,name=Avatar,proto3" json:"Avatar,omitempty"`
	Role      string                 `protobuf:"bytes,6,opt,name=Role,proto3" json:"Role,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *User) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *User) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *User) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type GetByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
}

func (x *GetByIDResponse) Reset() {
	*x = GetByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDResponse) ProtoMessage() {}

func (x *GetByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDResponse.ProtoReflect.Descriptor instead.
func (*GetByIDResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *GetByIDResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
}

func (x *GetByIDRequest) Reset() {
	*x = GetByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDRequest) ProtoMessage() {}

func (x *GetByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDRequest.ProtoReflect.Descriptor instead.
func (*GetByIDRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *GetByIDRequest) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type GetByIDsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*User `protobuf:"bytes,1,rep,name=Users,proto3" json:"Users,omitempty"`
}

func (x *GetByIDsRes) Reset() {
	*x = GetByIDsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIDsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDsRes) ProtoMessage() {}

func (x *GetByIDsRes) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDsRes.ProtoReflect.Descriptor instead.
func (*GetByIDsRes) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetByIDsRes) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetByIDsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UsersIDs []string `protobuf:"bytes,1,rep,name=UsersIDs,proto3" json:"UsersIDs,omitempty"`
}

func (x *GetByIDsReq) Reset() {
	*x = GetByIDsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIDsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDsReq) ProtoMessage() {}

func (x *GetByIDsReq) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDsReq.ProtoReflect.Descriptor instead.
func (*GetByIDsReq) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetByIDsReq) GetUsersIDs() []string {
	if x != nil {
		return x.UsersIDs
	}
	return nil
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x75, 0x73,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x02, 0x0a, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x46,
	0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x46, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x61, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x61, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x41,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x38, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x38, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25,
	0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x55, 0x73, 0x65, 0x72, 0x22, 0x28, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22,
	0x36, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x73, 0x52, 0x65, 0x73, 0x12, 0x27,
	0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x22, 0x29, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x42, 0x79,
	0x49, 0x44, 0x73, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x73, 0x49,
	0x44, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x73, 0x49,
	0x44, 0x73, 0x32, 0xa0, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4a, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x49,
	0x44, 0x12, 0x1b, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x45,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x42, 0x79, 0x49, 0x44, 0x73, 0x12,
	0x18, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65,
	0x74, 0x42, 0x79, 0x49, 0x44, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x73,
	0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_proto_goTypes = []interface{}{
	(*User)(nil),                  // 0: userService.User
	(*GetByIDResponse)(nil),       // 1: userService.GetByIDResponse
	(*GetByIDRequest)(nil),        // 2: userService.GetByIDRequest
	(*GetByIDsRes)(nil),           // 3: userService.GetByIDsRes
	(*GetByIDsReq)(nil),           // 4: userService.GetByIDsReq
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_user_proto_depIdxs = []int32{
	5, // 0: userService.User.CreatedAt:type_name -> google.protobuf.Timestamp
	5, // 1: userService.User.UpdatedAt:type_name -> google.protobuf.Timestamp
	0, // 2: userService.GetByIDResponse.User:type_name -> userService.User
	0, // 3: userService.GetByIDsRes.Users:type_name -> userService.User
	2, // 4: userService.UserService.GetUserByID:input_type -> userService.GetByIDRequest
	4, // 5: userService.UserService.GetUsersByIDs:input_type -> userService.GetByIDsReq
	1, // 6: userService.UserService.GetUserByID:output_type -> userService.GetByIDResponse
	3, // 7: userService.UserService.GetUsersByIDs:output_type -> userService.GetByIDsRes
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIDResponse); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIDRequest); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIDsRes); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIDsReq); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserServiceClient interface {
	GetUserByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error)
	GetUsersByIDs(ctx context.Context, in *GetByIDsReq, opts ...grpc.CallOption) (*GetByIDsRes, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetUserByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error) {
	out := new(GetByIDResponse)
	err := c.cc.Invoke(ctx, "/userService.UserService/GetUserByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUsersByIDs(ctx context.Context, in *GetByIDsReq, opts ...grpc.CallOption) (*GetByIDsRes, error) {
	out := new(GetByIDsRes)
	err := c.cc.Invoke(ctx, "/userService.UserService/GetUsersByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
type UserServiceServer interface {
	GetUserByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error)
	GetUsersByIDs(context.Context, *GetByIDsReq) (*GetByIDsRes, error)
}

// UnimplementedUserServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (*UnimplementedUserServiceServer) GetUserByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByID not implemented")
}
func (*UnimplementedUserServiceServer) GetUsersByIDs(context.Context, *GetByIDsReq) (*GetByIDsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersByIDs not implemented")
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_GetUserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userService.UserService/GetUserByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByID(ctx, req.(*GetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUsersByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUsersByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userService.UserService/GetUsersByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUsersByIDs(ctx, req.(*GetByIDsReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "userService.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserByID",
			Handler:    _UserService_GetUserByID_Handler,
		},
		{
			MethodName: "GetUsersByIDs",
			Handler:    _UserService_GetUsersByIDs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}