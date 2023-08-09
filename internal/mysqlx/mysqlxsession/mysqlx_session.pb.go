//
// Copyright (c) 2015, 2023, Oracle and/or its affiliates.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License, version 2.0,
// as published by the Free Software Foundation.
//
// This program is also distributed with certain software (including
// but not limited to OpenSSL) that is licensed under separate terms,
// as designated in a particular file or component or in included license
// documentation.  The authors of MySQL hereby grant you an additional
// permission to link the program and your derivative works with the
// separately licensed software that they have included with MySQL.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License, version 2.0, for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301  USA

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: mysqlx_session.proto

// ifdef PROTOBUF_LITE: option optimize_for = LITE_RUNTIME;

//*
//@namespace Mysqlx::Session
//@brief Messages to manage sessions.
//
//@startuml "Messages for Sessions"
//== session start ==
//Client -> Server: AuthenticateStart
//opt
//Server --> Client: AuthenticateContinue
//Client --> Server: AuthenticateContinue
//end
//alt
//Server --> Client: AuthenticateOk
//else
//Server --> Client: Error
//end
//...
//== session reset ==
//Client -> Server: Reset
//Server --> Client: Ok
//== session end ==
//Client -> Server: Close
//Server --> Client: Ok
//@enduml

package mysqlxsession

import (
	_ "github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
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

// *
// The initial message send from the client to the server to start
// the authentication process.
//
// @returns @ref Mysqlx::Session::AuthenticateContinue
type AuthenticateStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * authentication mechanism name
	MechName *string `protobuf:"bytes,1,req,name=mech_name,json=mechName" json:"mech_name,omitempty"`
	// * authentication data
	AuthData []byte `protobuf:"bytes,2,opt,name=auth_data,json=authData" json:"auth_data,omitempty"`
	// * initial response
	InitialResponse []byte `protobuf:"bytes,3,opt,name=initial_response,json=initialResponse" json:"initial_response,omitempty"`
}

func (x *AuthenticateStart) Reset() {
	*x = AuthenticateStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_session_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateStart) ProtoMessage() {}

func (x *AuthenticateStart) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_session_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateStart.ProtoReflect.Descriptor instead.
func (*AuthenticateStart) Descriptor() ([]byte, []int) {
	return file_mysqlx_session_proto_rawDescGZIP(), []int{0}
}

func (x *AuthenticateStart) GetMechName() string {
	if x != nil && x.MechName != nil {
		return *x.MechName
	}
	return ""
}

func (x *AuthenticateStart) GetAuthData() []byte {
	if x != nil {
		return x.AuthData
	}
	return nil
}

func (x *AuthenticateStart) GetInitialResponse() []byte {
	if x != nil {
		return x.InitialResponse
	}
	return nil
}

// *
// Send by client or server after an @ref Mysqlx::Session::AuthenticateStart
// to exchange more authentication data.
//
// @returns Mysqlx::Session::AuthenticateContinue
type AuthenticateContinue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * authentication data
	AuthData []byte `protobuf:"bytes,1,req,name=auth_data,json=authData" json:"auth_data,omitempty"`
}

func (x *AuthenticateContinue) Reset() {
	*x = AuthenticateContinue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_session_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateContinue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateContinue) ProtoMessage() {}

func (x *AuthenticateContinue) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_session_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateContinue.ProtoReflect.Descriptor instead.
func (*AuthenticateContinue) Descriptor() ([]byte, []int) {
	return file_mysqlx_session_proto_rawDescGZIP(), []int{1}
}

func (x *AuthenticateContinue) GetAuthData() []byte {
	if x != nil {
		return x.AuthData
	}
	return nil
}

// *
// Sent by the server after successful authentication.
type AuthenticateOk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * authentication data
	AuthData []byte `protobuf:"bytes,1,opt,name=auth_data,json=authData" json:"auth_data,omitempty"`
}

func (x *AuthenticateOk) Reset() {
	*x = AuthenticateOk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_session_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateOk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateOk) ProtoMessage() {}

func (x *AuthenticateOk) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_session_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateOk.ProtoReflect.Descriptor instead.
func (*AuthenticateOk) Descriptor() ([]byte, []int) {
	return file_mysqlx_session_proto_rawDescGZIP(), []int{2}
}

func (x *AuthenticateOk) GetAuthData() []byte {
	if x != nil {
		return x.AuthData
	}
	return nil
}

// *
// Reset the current session.
//
// @returns @ref Mysqlx::Ok
type Reset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * if is true the session will be reset, but stays authenticated; otherwise,
	// the session will be closed and needs to be authenticated again
	KeepOpen *bool `protobuf:"varint,1,opt,name=keep_open,json=keepOpen,def=0" json:"keep_open,omitempty"`
}

// Default values for Reset fields.
const (
	Default_Reset_KeepOpen = bool(false)
)

func (x *Reset) Reset() {
	*x = Reset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_session_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reset) ProtoMessage() {}

func (x *Reset) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_session_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reset.ProtoReflect.Descriptor instead.
func (*Reset) Descriptor() ([]byte, []int) {
	return file_mysqlx_session_proto_rawDescGZIP(), []int{3}
}

func (x *Reset) GetKeepOpen() bool {
	if x != nil && x.KeepOpen != nil {
		return *x.KeepOpen
	}
	return Default_Reset_KeepOpen
}

// *
// Close the current session.
//
// @returns @ref Mysqlx::Ok
type Close struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Close) Reset() {
	*x = Close{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_session_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Close) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close) ProtoMessage() {}

func (x *Close) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_session_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close.ProtoReflect.Descriptor instead.
func (*Close) Descriptor() ([]byte, []int) {
	return file_mysqlx_session_proto_rawDescGZIP(), []int{4}
}

var File_mysqlx_session_proto protoreflect.FileDescriptor

var file_mysqlx_session_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x0c, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7e, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x63,
	0x68, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65,
	0x63, 0x68, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x29, 0x0a, 0x10, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x3a, 0x04,
	0x88, 0xea, 0x30, 0x04, 0x22, 0x3d, 0x0a, 0x14, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0c, 0x52,
	0x08, 0x61, 0x75, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x3a, 0x08, 0x88, 0xea, 0x30, 0x05, 0x90,
	0xea, 0x30, 0x03, 0x22, 0x33, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x4f, 0x6b, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x44, 0x61,
	0x74, 0x61, 0x3a, 0x04, 0x90, 0xea, 0x30, 0x04, 0x22, 0x31, 0x0a, 0x05, 0x52, 0x65, 0x73, 0x65,
	0x74, 0x12, 0x22, 0x0a, 0x09, 0x6b, 0x65, 0x65, 0x70, 0x5f, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x08, 0x6b, 0x65, 0x65,
	0x70, 0x4f, 0x70, 0x65, 0x6e, 0x3a, 0x04, 0x88, 0xea, 0x30, 0x06, 0x22, 0x0d, 0x0a, 0x05, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x3a, 0x04, 0x88, 0xea, 0x30, 0x07, 0x42, 0x19, 0x0a, 0x17, 0x63, 0x6f,
	0x6d, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x2e, 0x63, 0x6a, 0x2e, 0x78, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66,
}

var (
	file_mysqlx_session_proto_rawDescOnce sync.Once
	file_mysqlx_session_proto_rawDescData = file_mysqlx_session_proto_rawDesc
)

func file_mysqlx_session_proto_rawDescGZIP() []byte {
	file_mysqlx_session_proto_rawDescOnce.Do(func() {
		file_mysqlx_session_proto_rawDescData = protoimpl.X.CompressGZIP(file_mysqlx_session_proto_rawDescData)
	})
	return file_mysqlx_session_proto_rawDescData
}

var file_mysqlx_session_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_mysqlx_session_proto_goTypes = []interface{}{
	(*AuthenticateStart)(nil),    // 0: Mysqlx.Session.AuthenticateStart
	(*AuthenticateContinue)(nil), // 1: Mysqlx.Session.AuthenticateContinue
	(*AuthenticateOk)(nil),       // 2: Mysqlx.Session.AuthenticateOk
	(*Reset)(nil),                // 3: Mysqlx.Session.Reset
	(*Close)(nil),                // 4: Mysqlx.Session.Close
}
var file_mysqlx_session_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mysqlx_session_proto_init() }
func file_mysqlx_session_proto_init() {
	if File_mysqlx_session_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mysqlx_session_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateStart); i {
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
		file_mysqlx_session_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateContinue); i {
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
		file_mysqlx_session_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateOk); i {
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
		file_mysqlx_session_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reset); i {
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
		file_mysqlx_session_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Close); i {
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
			RawDescriptor: file_mysqlx_session_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mysqlx_session_proto_goTypes,
		DependencyIndexes: file_mysqlx_session_proto_depIdxs,
		MessageInfos:      file_mysqlx_session_proto_msgTypes,
	}.Build()
	File_mysqlx_session_proto = out.File
	file_mysqlx_session_proto_rawDesc = nil
	file_mysqlx_session_proto_goTypes = nil
	file_mysqlx_session_proto_depIdxs = nil
}
