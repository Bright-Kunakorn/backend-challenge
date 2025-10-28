// Code generated manually to accompany proto/user.proto. DO NOT EDIT.

package userpb

import (
	reflect "reflect"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// User represents a user projection used in responses.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email     string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt string `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*User) Descriptor() ([]byte, []int) {
	return file_proto_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

// CreateUserRequest contains fields required to create a new user.
type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_user_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// CreateUserResponse wraps the created user and issued token.
type CreateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User  *User  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CreateUserResponse) Reset() {
	*x = CreateUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserResponse) ProtoMessage() {}

func (x *CreateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_user_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *CreateUserResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// GetUserRequest fetches a user by ID.
type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetUserResponse contains a user projection.
type GetUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_proto_user_proto protoreflect.FileDescriptor

var file_proto_user_proto_rawDescOnce sync.Once
var file_proto_user_proto_rawDescData = buildProtoUserRawDesc()

func file_proto_user_proto_rawDescGZIP() []byte {
	file_proto_user_proto_rawDescOnce.Do(func() {
		file_proto_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_user_proto_rawDescData)
	})
	return file_proto_user_proto_rawDescData
}

var file_proto_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_user_proto_goTypes = []interface{}{
	(*User)(nil),               // 0: user.v1.User
	(*CreateUserRequest)(nil),  // 1: user.v1.CreateUserRequest
	(*CreateUserResponse)(nil), // 2: user.v1.CreateUserResponse
	(*GetUserRequest)(nil),     // 3: user.v1.GetUserRequest
	(*GetUserResponse)(nil),    // 4: user.v1.GetUserResponse
}
var file_proto_user_proto_depIdxs = []int32{
	0, // 0: user.v1.CreateUserResponse.user:type_name -> user.v1.User
	0, // 1: user.v1.GetUserResponse.user:type_name -> user.v1.User
	1, // 2: user.v1.UserService.CreateUser:input_type -> user.v1.CreateUserRequest
	2, // 3: user.v1.UserService.CreateUser:output_type -> user.v1.CreateUserResponse
	3, // 4: user.v1.UserService.GetUser:input_type -> user.v1.GetUserRequest
	4, // 5: user.v1.UserService.GetUser:output_type -> user.v1.GetUserResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_user_proto_init() }

func file_proto_user_proto_init() {
	if File_proto_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserRequest); i {
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
		file_proto_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserResponse); i {
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
		file_proto_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserRequest); i {
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
		file_proto_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserResponse); i {
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
			RawDescriptor: file_proto_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_user_proto_goTypes,
		DependencyIndexes: file_proto_user_proto_depIdxs,
		MessageInfos:      file_proto_user_proto_msgTypes,
	}.Build()
	File_proto_user_proto = out.File
	file_proto_user_proto_rawDesc = nil
	file_proto_user_proto_goTypes = nil
	file_proto_user_proto_depIdxs = nil
}

var file_proto_user_proto_rawDesc []byte

func buildProtoUserRawDesc() []byte {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    strPtr("proto/user.proto"),
		Package: strPtr("user.v1"),
		Syntax:  strPtr("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: strPtr("backend-challenge/proto/userpb;userpb"),
		},
	}

	fd.MessageType = []*descriptorpb.DescriptorProto{
		buildUserMessage(),
		buildCreateUserRequestMessage(),
		buildCreateUserResponseMessage(),
		buildGetUserRequestMessage(),
		buildGetUserResponseMessage(),
	}

	fd.Service = []*descriptorpb.ServiceDescriptorProto{
		buildUserServiceDescriptor(),
	}

	b, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}

	file_proto_user_proto_rawDesc = b
	return b
}

func buildUserMessage() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: strPtr("User"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     strPtr("id"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("id"),
			},
			{
				Name:     strPtr("name"),
				Number:   int32Ptr(2),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("name"),
			},
			{
				Name:     strPtr("email"),
				Number:   int32Ptr(3),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("email"),
			},
			{
				Name:     strPtr("created_at"),
				Number:   int32Ptr(4),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("createdAt"),
			},
		},
	}
}

func buildCreateUserRequestMessage() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: strPtr("CreateUserRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     strPtr("name"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("name"),
			},
			{
				Name:     strPtr("email"),
				Number:   int32Ptr(2),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("email"),
			},
			{
				Name:     strPtr("password"),
				Number:   int32Ptr(3),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("password"),
			},
		},
	}
}

func buildCreateUserResponseMessage() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: strPtr("CreateUserResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     strPtr("user"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: strPtr(".user.v1.User"),
				JsonName: strPtr("user"),
			},
			{
				Name:     strPtr("token"),
				Number:   int32Ptr(2),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("token"),
			},
		},
	}
}

func buildGetUserRequestMessage() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: strPtr("GetUserRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     strPtr("id"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: strPtr("id"),
			},
		},
	}
}

func buildGetUserResponseMessage() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: strPtr("GetUserResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     strPtr("user"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: strPtr(".user.v1.User"),
				JsonName: strPtr("user"),
			},
		},
	}
}

func buildUserServiceDescriptor() *descriptorpb.ServiceDescriptorProto {
	return &descriptorpb.ServiceDescriptorProto{
		Name: strPtr("UserService"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       strPtr("CreateUser"),
				InputType:  strPtr(".user.v1.CreateUserRequest"),
				OutputType: strPtr(".user.v1.CreateUserResponse"),
			},
			{
				Name:       strPtr("GetUser"),
				InputType:  strPtr(".user.v1.GetUserRequest"),
				OutputType: strPtr(".user.v1.GetUserResponse"),
			},
		},
	}
}

func strPtr(s string) *string { return &s }
func int32Ptr(v int32) *int32 { return &v }
