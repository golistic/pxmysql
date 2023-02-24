//
// Copyright (c) 2015, 2022, Oracle and/or its affiliates.
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

//* @page mysqlx_protocol_messages Messages
//Topics in this section:
//
//-  @ref messages_Message_Structure
//-  @ref messages_Message_Sequence
//-  @ref messages_Common_Messages
//-  @ref Mysqlx::Connection "Connection"
//-  @ref Mysqlx::Session "Session"
//-  @ref Mysqlx::Expect "Expectations"
//-  @ref Mysqlx::Crud "CRUD"
//-  @ref Mysqlx::Sql "SQL"
//-  @ref Mysqlx::Resultset "Resultset"
//-  @ref Mysqlx::Expr "Expressions
//-  @ref Mysqlx::Datatypes "Datatypes"
//-  @ref Mysqlx::Notice "Notice"
//-  @ref Mysqlx::Prepare "Prepared statments"
//-  @ref Mysqlx::Cursor "Cursor"
//
//
//This section provides detailed information about how X %Protocol
//defines messages.
//
//@section messages_Message_Structure Message Structure
//
//Messages have a:
//-  4 byte *length* (little endian)
//-  1 byte *message type*
//-  a ``message_payload`` of length ``.length - 1``
//
//@par Mysqlx.Message
//Container of all messages that are exchanged between client and server.
//@n@b Parameters
//-  ``length`` -- length of the whole message
//-  ``message_type`` -- type of the ``message_payload``
//-  ``message_payload`` -- the message's payload encoded using
//[`Google Protobuf`](https://code.google.com/p/protobuf/) if
//not otherwise noted.
//
//@code{unparsed}
//struct Message {
//uint32          length;
//uint8           message_type;
//opaque          message_payload[Message.length - 1];
//};
//@endcode
//
//@note
//The ``message_payload`` is generated from the protobuf files using
//``protoc``:
//@code{shell}
//$ protoc --cpp_out=protodir mysqlx*.proto
//@endcode
//-  [``mysqlx.proto``]
//-  [``mysqlx_connection.proto``]
//-  [``mysqlx_session.proto``]
//-  [``mysqlx_crud.proto``]
//-  [``mysqlx_sql.proto``]
//-  [``mysqlx_resultset.proto``]
//-  [``mysqlx_expr.proto``]
//-  [``mysqlx_datatypes.proto``]
//-  [``mysqlx_expect.proto``]
//-  [``mysqlx_notice.proto``]
//
//@par
//
//@note
//The ``message_type`` can be taken from the
//@ref Mysqlx::ClientMessages for client-messages and from
//@ref Mysqlx::ServerMessages of server-side messages.
//@n In ``C++`` they are exposed in ``mysqlx.pb.h`` in the
//``ClientMessages`` class.
//@code{unparsed}
//ClientMessages.MsgCase.kMsgConGetCap
//ClientMessages.kMsgConGetCapFieldNumber
//@endcode
//
//
//@section messages_Message_Sequence Message Sequence
//
//Messages usually appear in a sequence. Each initial message (one
//referenced by @ref Mysqlx::ClientMessages) is
//associated with a set of possible following messages.
//
//A message sequence either:
//-  finishes successfully if it reaches its end-state or
//-  is aborted with a @ref Mysqlx::Error message
//
//At any time in between local @ref Mysqlx::Notice "Notices"
//may be sent by the server as part of the message sequence.
//
//Global @ref Mysqlx::Notice "Notices" may be sent by the  server at any time.
//
//
//@section messages_Common_Messages Common Messages
//
//@subsection messages_Error_Message Error Message
//
//After the client sent the initial message, the server may send a
//@ref Mysqlx::Error message at any time to terminate the
//current message sequence.

// tell protobuf 3.0 to use protobuf 2.x rules

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: mysqlx.proto

// ifdef PROTOBUF_LITE: option optimize_for = LITE_RUNTIME;

//*
//@namespace Mysqlx

package mysqlx

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ClientMessages_Type int32

const (
	ClientMessages_CON_CAPABILITIES_GET       ClientMessages_Type = 1
	ClientMessages_CON_CAPABILITIES_SET       ClientMessages_Type = 2
	ClientMessages_CON_CLOSE                  ClientMessages_Type = 3
	ClientMessages_SESS_AUTHENTICATE_START    ClientMessages_Type = 4
	ClientMessages_SESS_AUTHENTICATE_CONTINUE ClientMessages_Type = 5
	ClientMessages_SESS_RESET                 ClientMessages_Type = 6
	ClientMessages_SESS_CLOSE                 ClientMessages_Type = 7
	ClientMessages_SQL_STMT_EXECUTE           ClientMessages_Type = 12
	ClientMessages_CRUD_FIND                  ClientMessages_Type = 17
	ClientMessages_CRUD_INSERT                ClientMessages_Type = 18
	ClientMessages_CRUD_UPDATE                ClientMessages_Type = 19
	ClientMessages_CRUD_DELETE                ClientMessages_Type = 20
	ClientMessages_EXPECT_OPEN                ClientMessages_Type = 24
	ClientMessages_EXPECT_CLOSE               ClientMessages_Type = 25
	ClientMessages_CRUD_CREATE_VIEW           ClientMessages_Type = 30
	ClientMessages_CRUD_MODIFY_VIEW           ClientMessages_Type = 31
	ClientMessages_CRUD_DROP_VIEW             ClientMessages_Type = 32
	ClientMessages_PREPARE_PREPARE            ClientMessages_Type = 40
	ClientMessages_PREPARE_EXECUTE            ClientMessages_Type = 41
	ClientMessages_PREPARE_DEALLOCATE         ClientMessages_Type = 42
	ClientMessages_CURSOR_OPEN                ClientMessages_Type = 43
	ClientMessages_CURSOR_CLOSE               ClientMessages_Type = 44
	ClientMessages_CURSOR_FETCH               ClientMessages_Type = 45
	ClientMessages_COMPRESSION                ClientMessages_Type = 46
)

// Enum value maps for ClientMessages_Type.
var (
	ClientMessages_Type_name = map[int32]string{
		1:  "CON_CAPABILITIES_GET",
		2:  "CON_CAPABILITIES_SET",
		3:  "CON_CLOSE",
		4:  "SESS_AUTHENTICATE_START",
		5:  "SESS_AUTHENTICATE_CONTINUE",
		6:  "SESS_RESET",
		7:  "SESS_CLOSE",
		12: "SQL_STMT_EXECUTE",
		17: "CRUD_FIND",
		18: "CRUD_INSERT",
		19: "CRUD_UPDATE",
		20: "CRUD_DELETE",
		24: "EXPECT_OPEN",
		25: "EXPECT_CLOSE",
		30: "CRUD_CREATE_VIEW",
		31: "CRUD_MODIFY_VIEW",
		32: "CRUD_DROP_VIEW",
		40: "PREPARE_PREPARE",
		41: "PREPARE_EXECUTE",
		42: "PREPARE_DEALLOCATE",
		43: "CURSOR_OPEN",
		44: "CURSOR_CLOSE",
		45: "CURSOR_FETCH",
		46: "COMPRESSION",
	}
	ClientMessages_Type_value = map[string]int32{
		"CON_CAPABILITIES_GET":       1,
		"CON_CAPABILITIES_SET":       2,
		"CON_CLOSE":                  3,
		"SESS_AUTHENTICATE_START":    4,
		"SESS_AUTHENTICATE_CONTINUE": 5,
		"SESS_RESET":                 6,
		"SESS_CLOSE":                 7,
		"SQL_STMT_EXECUTE":           12,
		"CRUD_FIND":                  17,
		"CRUD_INSERT":                18,
		"CRUD_UPDATE":                19,
		"CRUD_DELETE":                20,
		"EXPECT_OPEN":                24,
		"EXPECT_CLOSE":               25,
		"CRUD_CREATE_VIEW":           30,
		"CRUD_MODIFY_VIEW":           31,
		"CRUD_DROP_VIEW":             32,
		"PREPARE_PREPARE":            40,
		"PREPARE_EXECUTE":            41,
		"PREPARE_DEALLOCATE":         42,
		"CURSOR_OPEN":                43,
		"CURSOR_CLOSE":               44,
		"CURSOR_FETCH":               45,
		"COMPRESSION":                46,
	}
)

func (x ClientMessages_Type) Enum() *ClientMessages_Type {
	p := new(ClientMessages_Type)
	*p = x
	return p
}

func (x ClientMessages_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientMessages_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_proto_enumTypes[0].Descriptor()
}

func (ClientMessages_Type) Type() protoreflect.EnumType {
	return &file_mysqlx_proto_enumTypes[0]
}

func (x ClientMessages_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *ClientMessages_Type) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = ClientMessages_Type(num)
	return nil
}

// Deprecated: Use ClientMessages_Type.Descriptor instead.
func (ClientMessages_Type) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{0, 0}
}

type ServerMessages_Type int32

const (
	ServerMessages_OK                         ServerMessages_Type = 0
	ServerMessages_ERROR                      ServerMessages_Type = 1
	ServerMessages_CONN_CAPABILITIES          ServerMessages_Type = 2
	ServerMessages_SESS_AUTHENTICATE_CONTINUE ServerMessages_Type = 3
	ServerMessages_SESS_AUTHENTICATE_OK       ServerMessages_Type = 4
	// NOTICE has to stay at 11 forever
	ServerMessages_NOTICE                               ServerMessages_Type = 11
	ServerMessages_RESULTSET_COLUMN_META_DATA           ServerMessages_Type = 12
	ServerMessages_RESULTSET_ROW                        ServerMessages_Type = 13
	ServerMessages_RESULTSET_FETCH_DONE                 ServerMessages_Type = 14
	ServerMessages_RESULTSET_FETCH_SUSPENDED            ServerMessages_Type = 15
	ServerMessages_RESULTSET_FETCH_DONE_MORE_RESULTSETS ServerMessages_Type = 16
	ServerMessages_SQL_STMT_EXECUTE_OK                  ServerMessages_Type = 17
	ServerMessages_RESULTSET_FETCH_DONE_MORE_OUT_PARAMS ServerMessages_Type = 18
	ServerMessages_COMPRESSION                          ServerMessages_Type = 19
)

// Enum value maps for ServerMessages_Type.
var (
	ServerMessages_Type_name = map[int32]string{
		0:  "OK",
		1:  "ERROR",
		2:  "CONN_CAPABILITIES",
		3:  "SESS_AUTHENTICATE_CONTINUE",
		4:  "SESS_AUTHENTICATE_OK",
		11: "NOTICE",
		12: "RESULTSET_COLUMN_META_DATA",
		13: "RESULTSET_ROW",
		14: "RESULTSET_FETCH_DONE",
		15: "RESULTSET_FETCH_SUSPENDED",
		16: "RESULTSET_FETCH_DONE_MORE_RESULTSETS",
		17: "SQL_STMT_EXECUTE_OK",
		18: "RESULTSET_FETCH_DONE_MORE_OUT_PARAMS",
		19: "COMPRESSION",
	}
	ServerMessages_Type_value = map[string]int32{
		"OK":                                   0,
		"ERROR":                                1,
		"CONN_CAPABILITIES":                    2,
		"SESS_AUTHENTICATE_CONTINUE":           3,
		"SESS_AUTHENTICATE_OK":                 4,
		"NOTICE":                               11,
		"RESULTSET_COLUMN_META_DATA":           12,
		"RESULTSET_ROW":                        13,
		"RESULTSET_FETCH_DONE":                 14,
		"RESULTSET_FETCH_SUSPENDED":            15,
		"RESULTSET_FETCH_DONE_MORE_RESULTSETS": 16,
		"SQL_STMT_EXECUTE_OK":                  17,
		"RESULTSET_FETCH_DONE_MORE_OUT_PARAMS": 18,
		"COMPRESSION":                          19,
	}
)

func (x ServerMessages_Type) Enum() *ServerMessages_Type {
	p := new(ServerMessages_Type)
	*p = x
	return p
}

func (x ServerMessages_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerMessages_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_proto_enumTypes[1].Descriptor()
}

func (ServerMessages_Type) Type() protoreflect.EnumType {
	return &file_mysqlx_proto_enumTypes[1]
}

func (x ServerMessages_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *ServerMessages_Type) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = ServerMessages_Type(num)
	return nil
}

// Deprecated: Use ServerMessages_Type.Descriptor instead.
func (ServerMessages_Type) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{1, 0}
}

type Error_Severity int32

const (
	Error_ERROR Error_Severity = 0
	Error_FATAL Error_Severity = 1
)

// Enum value maps for Error_Severity.
var (
	Error_Severity_name = map[int32]string{
		0: "ERROR",
		1: "FATAL",
	}
	Error_Severity_value = map[string]int32{
		"ERROR": 0,
		"FATAL": 1,
	}
)

func (x Error_Severity) Enum() *Error_Severity {
	p := new(Error_Severity)
	*p = x
	return p
}

func (x Error_Severity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error_Severity) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_proto_enumTypes[2].Descriptor()
}

func (Error_Severity) Type() protoreflect.EnumType {
	return &file_mysqlx_proto_enumTypes[2]
}

func (x Error_Severity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Error_Severity) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Error_Severity(num)
	return nil
}

// Deprecated: Use Error_Severity.Descriptor instead.
func (Error_Severity) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{3, 0}
}

// *
// IDs of messages that can be sent from client to the server.
//
// @note
// This message is never sent on the wire. It is only used to let “protoc“:
// -  generate constants
// -  check for uniqueness
type ClientMessages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClientMessages) Reset() {
	*x = ClientMessages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessages) ProtoMessage() {}

func (x *ClientMessages) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessages.ProtoReflect.Descriptor instead.
func (*ClientMessages) Descriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{0}
}

// *
// IDs of messages that can be sent from server to client.
//
// @note
// This message is never sent on the wire. It is only used to let “protoc“:
// -  generate constants
// -  check for uniqueness
type ServerMessages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerMessages) Reset() {
	*x = ServerMessages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessages) ProtoMessage() {}

func (x *ServerMessages) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessages.ProtoReflect.Descriptor instead.
func (*ServerMessages) Descriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{1}
}

type Ok struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg *string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (x *Ok) Reset() {
	*x = Ok{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ok) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ok) ProtoMessage() {}

func (x *Ok) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ok.ProtoReflect.Descriptor instead.
func (*Ok) Descriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{2}
}

func (x *Ok) GetMsg() string {
	if x != nil && x.Msg != nil {
		return *x.Msg
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * severity of the error message
	Severity *Error_Severity `protobuf:"varint,1,opt,name=severity,enum=Mysqlx.Error_Severity,def=0" json:"severity,omitempty"`
	// * error code
	Code *uint32 `protobuf:"varint,2,req,name=code" json:"code,omitempty"`
	// * SQL state
	SqlState *string `protobuf:"bytes,4,req,name=sql_state,json=sqlState" json:"sql_state,omitempty"`
	// * human-readable error message
	Msg *string `protobuf:"bytes,3,req,name=msg" json:"msg,omitempty"`
}

// Default values for Error fields.
const (
	Default_Error_Severity = Error_ERROR
)

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_mysqlx_proto_rawDescGZIP(), []int{3}
}

func (x *Error) GetSeverity() Error_Severity {
	if x != nil && x.Severity != nil {
		return *x.Severity
	}
	return Default_Error_Severity
}

func (x *Error) GetCode() uint32 {
	if x != nil && x.Code != nil {
		return *x.Code
	}
	return 0
}

func (x *Error) GetSqlState() string {
	if x != nil && x.SqlState != nil {
		return *x.SqlState
	}
	return ""
}

func (x *Error) GetMsg() string {
	if x != nil && x.Msg != nil {
		return *x.Msg
	}
	return ""
}

var file_mysqlx_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*ClientMessages_Type)(nil),
		Field:         100001,
		Name:          "Mysqlx.client_message_id",
		Tag:           "varint,100001,opt,name=client_message_id,enum=Mysqlx.ClientMessages_Type",
		Filename:      "mysqlx.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*ServerMessages_Type)(nil),
		Field:         100002,
		Name:          "Mysqlx.server_message_id",
		Tag:           "varint,100002,opt,name=server_message_id,enum=Mysqlx.ServerMessages_Type",
		Filename:      "mysqlx.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional Mysqlx.ClientMessages.Type client_message_id = 100001;
	E_ClientMessageId = &file_mysqlx_proto_extTypes[0]
	// optional Mysqlx.ServerMessages.Type server_message_id = 100002;
	E_ServerMessageId = &file_mysqlx_proto_extTypes[1]
)

var File_mysqlx_proto protoreflect.FileDescriptor

var file_mysqlx_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfc, 0x03, 0x0a, 0x0e, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0xe9, 0x03, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x50, 0x41,
	0x42, 0x49, 0x4c, 0x49, 0x54, 0x49, 0x45, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x10, 0x01, 0x12, 0x18,
	0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x49,
	0x45, 0x53, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4e, 0x5f,
	0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x45, 0x53, 0x53, 0x5f,
	0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x54, 0x41,
	0x52, 0x54, 0x10, 0x04, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53, 0x5f, 0x41, 0x55, 0x54,
	0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x49, 0x4e,
	0x55, 0x45, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x45, 0x53, 0x53, 0x5f, 0x52, 0x45, 0x53,
	0x45, 0x54, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x45, 0x53, 0x53, 0x5f, 0x43, 0x4c, 0x4f,
	0x53, 0x45, 0x10, 0x07, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x51, 0x4c, 0x5f, 0x53, 0x54, 0x4d, 0x54,
	0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x45, 0x10, 0x0c, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x52,
	0x55, 0x44, 0x5f, 0x46, 0x49, 0x4e, 0x44, 0x10, 0x11, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x52, 0x55,
	0x44, 0x5f, 0x49, 0x4e, 0x53, 0x45, 0x52, 0x54, 0x10, 0x12, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x52,
	0x55, 0x44, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x13, 0x12, 0x0f, 0x0a, 0x0b, 0x43,
	0x52, 0x55, 0x44, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x14, 0x12, 0x0f, 0x0a, 0x0b,
	0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x10, 0x18, 0x12, 0x10, 0x0a,
	0x0c, 0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x19, 0x12,
	0x14, 0x0a, 0x10, 0x43, 0x52, 0x55, 0x44, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x56,
	0x49, 0x45, 0x57, 0x10, 0x1e, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x52, 0x55, 0x44, 0x5f, 0x4d, 0x4f,
	0x44, 0x49, 0x46, 0x59, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x10, 0x1f, 0x12, 0x12, 0x0a, 0x0e, 0x43,
	0x52, 0x55, 0x44, 0x5f, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x10, 0x20, 0x12,
	0x13, 0x0a, 0x0f, 0x50, 0x52, 0x45, 0x50, 0x41, 0x52, 0x45, 0x5f, 0x50, 0x52, 0x45, 0x50, 0x41,
	0x52, 0x45, 0x10, 0x28, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x52, 0x45, 0x50, 0x41, 0x52, 0x45, 0x5f,
	0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x45, 0x10, 0x29, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x52, 0x45,
	0x50, 0x41, 0x52, 0x45, 0x5f, 0x44, 0x45, 0x41, 0x4c, 0x4c, 0x4f, 0x43, 0x41, 0x54, 0x45, 0x10,
	0x2a, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x55, 0x52, 0x53, 0x4f, 0x52, 0x5f, 0x4f, 0x50, 0x45, 0x4e,
	0x10, 0x2b, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x55, 0x52, 0x53, 0x4f, 0x52, 0x5f, 0x43, 0x4c, 0x4f,
	0x53, 0x45, 0x10, 0x2c, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x55, 0x52, 0x53, 0x4f, 0x52, 0x5f, 0x46,
	0x45, 0x54, 0x43, 0x48, 0x10, 0x2d, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x4d, 0x50, 0x52, 0x45,
	0x53, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x2e, 0x22, 0xf3, 0x02, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0xe0, 0x02, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x4e, 0x4e, 0x5f, 0x43,
	0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x49, 0x45, 0x53, 0x10, 0x02, 0x12, 0x1e, 0x0a,
	0x1a, 0x53, 0x45, 0x53, 0x53, 0x5f, 0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41,
	0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x49, 0x4e, 0x55, 0x45, 0x10, 0x03, 0x12, 0x18, 0x0a,
	0x14, 0x53, 0x45, 0x53, 0x53, 0x5f, 0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41,
	0x54, 0x45, 0x5f, 0x4f, 0x4b, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x54, 0x49, 0x43,
	0x45, 0x10, 0x0b, 0x12, 0x1e, 0x0a, 0x1a, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x45, 0x54,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x4d, 0x45, 0x54, 0x41, 0x5f, 0x44, 0x41, 0x54,
	0x41, 0x10, 0x0c, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x45, 0x54,
	0x5f, 0x52, 0x4f, 0x57, 0x10, 0x0d, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54,
	0x53, 0x45, 0x54, 0x5f, 0x46, 0x45, 0x54, 0x43, 0x48, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x10, 0x0e,
	0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x45, 0x54, 0x5f, 0x46, 0x45,
	0x54, 0x43, 0x48, 0x5f, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x0f, 0x12,
	0x28, 0x0a, 0x24, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x45, 0x54, 0x5f, 0x46, 0x45, 0x54,
	0x43, 0x48, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x5f, 0x4d, 0x4f, 0x52, 0x45, 0x5f, 0x52, 0x45, 0x53,
	0x55, 0x4c, 0x54, 0x53, 0x45, 0x54, 0x53, 0x10, 0x10, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x51, 0x4c,
	0x5f, 0x53, 0x54, 0x4d, 0x54, 0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x45, 0x5f, 0x4f, 0x4b,
	0x10, 0x11, 0x12, 0x28, 0x0a, 0x24, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x45, 0x54, 0x5f,
	0x46, 0x45, 0x54, 0x43, 0x48, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x5f, 0x4d, 0x4f, 0x52, 0x45, 0x5f,
	0x4f, 0x55, 0x54, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x53, 0x10, 0x12, 0x12, 0x0f, 0x0a, 0x0b,
	0x43, 0x4f, 0x4d, 0x50, 0x52, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x13, 0x22, 0x1c, 0x0a,
	0x02, 0x4f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x73, 0x67, 0x3a, 0x04, 0x90, 0xea, 0x30, 0x00, 0x22, 0xad, 0x01, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x39, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x3a,
	0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0d, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x71, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x73, 0x71, 0x6c, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x02, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x73, 0x67, 0x22, 0x20, 0x0a, 0x08, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x41,
	0x54, 0x41, 0x4c, 0x10, 0x01, 0x3a, 0x04, 0x90, 0xea, 0x30, 0x01, 0x3a, 0x6a, 0x0a, 0x11, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xa1, 0x8d, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x4d, 0x79, 0x73, 0x71,
	0x6c, 0x78, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x3a, 0x6a, 0x0a, 0x11, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x12, 0x1f, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x8d,
	0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x64, 0x42, 0x19, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c,
	0x2e, 0x63, 0x6a, 0x2e, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
}

var (
	file_mysqlx_proto_rawDescOnce sync.Once
	file_mysqlx_proto_rawDescData = file_mysqlx_proto_rawDesc
)

func file_mysqlx_proto_rawDescGZIP() []byte {
	file_mysqlx_proto_rawDescOnce.Do(func() {
		file_mysqlx_proto_rawDescData = protoimpl.X.CompressGZIP(file_mysqlx_proto_rawDescData)
	})
	return file_mysqlx_proto_rawDescData
}

var file_mysqlx_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_mysqlx_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mysqlx_proto_goTypes = []interface{}{
	(ClientMessages_Type)(0),            // 0: Mysqlx.ClientMessages.Type
	(ServerMessages_Type)(0),            // 1: Mysqlx.ServerMessages.Type
	(Error_Severity)(0),                 // 2: Mysqlx.Error.Severity
	(*ClientMessages)(nil),              // 3: Mysqlx.ClientMessages
	(*ServerMessages)(nil),              // 4: Mysqlx.ServerMessages
	(*Ok)(nil),                          // 5: Mysqlx.Ok
	(*Error)(nil),                       // 6: Mysqlx.Error
	(*descriptorpb.MessageOptions)(nil), // 7: google.protobuf.MessageOptions
}
var file_mysqlx_proto_depIdxs = []int32{
	2, // 0: Mysqlx.Error.severity:type_name -> Mysqlx.Error.Severity
	7, // 1: Mysqlx.client_message_id:extendee -> google.protobuf.MessageOptions
	7, // 2: Mysqlx.server_message_id:extendee -> google.protobuf.MessageOptions
	0, // 3: Mysqlx.client_message_id:type_name -> Mysqlx.ClientMessages.Type
	1, // 4: Mysqlx.server_message_id:type_name -> Mysqlx.ServerMessages.Type
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	3, // [3:5] is the sub-list for extension type_name
	1, // [1:3] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_mysqlx_proto_init() }
func file_mysqlx_proto_init() {
	if File_mysqlx_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mysqlx_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessages); i {
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
		file_mysqlx_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessages); i {
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
		file_mysqlx_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ok); i {
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
		file_mysqlx_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_mysqlx_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   4,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_mysqlx_proto_goTypes,
		DependencyIndexes: file_mysqlx_proto_depIdxs,
		EnumInfos:         file_mysqlx_proto_enumTypes,
		MessageInfos:      file_mysqlx_proto_msgTypes,
		ExtensionInfos:    file_mysqlx_proto_extTypes,
	}.Build()
	File_mysqlx_proto = out.File
	file_mysqlx_proto_rawDesc = nil
	file_mysqlx_proto_goTypes = nil
	file_mysqlx_proto_depIdxs = nil
}
