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

// tell protobuf 3.0 to use protobuf 2.x rules

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: mysqlx_notice.proto

// ifdef PROTOBUF_LITE: option optimize_for = LITE_RUNTIME;

//*
//@namespace Mysqlx::Notice
//@brief A notice
//- is sent from the server to the client
//- may be global or relate to the current message sequence
//
//The server may send notices @ref Mysqlx::Notice::Frame
//to the client at any time.
//
//A notice can be:
//-  global (``.scope == GLOBAL``) or
//-  belong to the currently executed @ref messages_Message_Sequence
//(``.scope == LOCAL + message sequence is active``):
//
//@note
//If the Server sends a ``LOCAL`` notice while no message sequence is
//active, the Notice should be ignored.
//
//@par Tip
//For more information, see @ref mysqlx_protocol_notices "Notices".

package mysqlxnotice

import (
	_ "github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	mysqlxdatatypes "github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
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

// * scope of notice
type Frame_Scope int32

const (
	Frame_GLOBAL Frame_Scope = 1
	Frame_LOCAL  Frame_Scope = 2
)

// Enum value maps for Frame_Scope.
var (
	Frame_Scope_name = map[int32]string{
		1: "GLOBAL",
		2: "LOCAL",
	}
	Frame_Scope_value = map[string]int32{
		"GLOBAL": 1,
		"LOCAL":  2,
	}
)

func (x Frame_Scope) Enum() *Frame_Scope {
	p := new(Frame_Scope)
	*p = x
	return p
}

func (x Frame_Scope) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Frame_Scope) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_notice_proto_enumTypes[0].Descriptor()
}

func (Frame_Scope) Type() protoreflect.EnumType {
	return &file_mysqlx_notice_proto_enumTypes[0]
}

func (x Frame_Scope) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Frame_Scope) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Frame_Scope(num)
	return nil
}

// Deprecated: Use Frame_Scope.Descriptor instead.
func (Frame_Scope) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{0, 0}
}

// * type of notice payload
type Frame_Type int32

const (
	Frame_WARNING                         Frame_Type = 1
	Frame_SESSION_VARIABLE_CHANGED        Frame_Type = 2
	Frame_SESSION_STATE_CHANGED           Frame_Type = 3
	Frame_GROUP_REPLICATION_STATE_CHANGED Frame_Type = 4
	Frame_SERVER_HELLO                    Frame_Type = 5
)

// Enum value maps for Frame_Type.
var (
	Frame_Type_name = map[int32]string{
		1: "WARNING",
		2: "SESSION_VARIABLE_CHANGED",
		3: "SESSION_STATE_CHANGED",
		4: "GROUP_REPLICATION_STATE_CHANGED",
		5: "SERVER_HELLO",
	}
	Frame_Type_value = map[string]int32{
		"WARNING":                         1,
		"SESSION_VARIABLE_CHANGED":        2,
		"SESSION_STATE_CHANGED":           3,
		"GROUP_REPLICATION_STATE_CHANGED": 4,
		"SERVER_HELLO":                    5,
	}
)

func (x Frame_Type) Enum() *Frame_Type {
	p := new(Frame_Type)
	*p = x
	return p
}

func (x Frame_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Frame_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_notice_proto_enumTypes[1].Descriptor()
}

func (Frame_Type) Type() protoreflect.EnumType {
	return &file_mysqlx_notice_proto_enumTypes[1]
}

func (x Frame_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Frame_Type) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Frame_Type(num)
	return nil
}

// Deprecated: Use Frame_Type.Descriptor instead.
func (Frame_Type) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{0, 1}
}

type Warning_Level int32

const (
	Warning_NOTE    Warning_Level = 1
	Warning_WARNING Warning_Level = 2
	Warning_ERROR   Warning_Level = 3
)

// Enum value maps for Warning_Level.
var (
	Warning_Level_name = map[int32]string{
		1: "NOTE",
		2: "WARNING",
		3: "ERROR",
	}
	Warning_Level_value = map[string]int32{
		"NOTE":    1,
		"WARNING": 2,
		"ERROR":   3,
	}
)

func (x Warning_Level) Enum() *Warning_Level {
	p := new(Warning_Level)
	*p = x
	return p
}

func (x Warning_Level) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Warning_Level) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_notice_proto_enumTypes[2].Descriptor()
}

func (Warning_Level) Type() protoreflect.EnumType {
	return &file_mysqlx_notice_proto_enumTypes[2]
}

func (x Warning_Level) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Warning_Level) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Warning_Level(num)
	return nil
}

// Deprecated: Use Warning_Level.Descriptor instead.
func (Warning_Level) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{1, 0}
}

type SessionStateChanged_Parameter int32

const (
	SessionStateChanged_CURRENT_SCHEMA         SessionStateChanged_Parameter = 1
	SessionStateChanged_ACCOUNT_EXPIRED        SessionStateChanged_Parameter = 2
	SessionStateChanged_GENERATED_INSERT_ID    SessionStateChanged_Parameter = 3
	SessionStateChanged_ROWS_AFFECTED          SessionStateChanged_Parameter = 4
	SessionStateChanged_ROWS_FOUND             SessionStateChanged_Parameter = 5
	SessionStateChanged_ROWS_MATCHED           SessionStateChanged_Parameter = 6
	SessionStateChanged_TRX_COMMITTED          SessionStateChanged_Parameter = 7
	SessionStateChanged_TRX_ROLLEDBACK         SessionStateChanged_Parameter = 9
	SessionStateChanged_PRODUCED_MESSAGE       SessionStateChanged_Parameter = 10
	SessionStateChanged_CLIENT_ID_ASSIGNED     SessionStateChanged_Parameter = 11
	SessionStateChanged_GENERATED_DOCUMENT_IDS SessionStateChanged_Parameter = 12 // .. more to be added
)

// Enum value maps for SessionStateChanged_Parameter.
var (
	SessionStateChanged_Parameter_name = map[int32]string{
		1:  "CURRENT_SCHEMA",
		2:  "ACCOUNT_EXPIRED",
		3:  "GENERATED_INSERT_ID",
		4:  "ROWS_AFFECTED",
		5:  "ROWS_FOUND",
		6:  "ROWS_MATCHED",
		7:  "TRX_COMMITTED",
		9:  "TRX_ROLLEDBACK",
		10: "PRODUCED_MESSAGE",
		11: "CLIENT_ID_ASSIGNED",
		12: "GENERATED_DOCUMENT_IDS",
	}
	SessionStateChanged_Parameter_value = map[string]int32{
		"CURRENT_SCHEMA":         1,
		"ACCOUNT_EXPIRED":        2,
		"GENERATED_INSERT_ID":    3,
		"ROWS_AFFECTED":          4,
		"ROWS_FOUND":             5,
		"ROWS_MATCHED":           6,
		"TRX_COMMITTED":          7,
		"TRX_ROLLEDBACK":         9,
		"PRODUCED_MESSAGE":       10,
		"CLIENT_ID_ASSIGNED":     11,
		"GENERATED_DOCUMENT_IDS": 12,
	}
)

func (x SessionStateChanged_Parameter) Enum() *SessionStateChanged_Parameter {
	p := new(SessionStateChanged_Parameter)
	*p = x
	return p
}

func (x SessionStateChanged_Parameter) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SessionStateChanged_Parameter) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_notice_proto_enumTypes[3].Descriptor()
}

func (SessionStateChanged_Parameter) Type() protoreflect.EnumType {
	return &file_mysqlx_notice_proto_enumTypes[3]
}

func (x SessionStateChanged_Parameter) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *SessionStateChanged_Parameter) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = SessionStateChanged_Parameter(num)
	return nil
}

// Deprecated: Use SessionStateChanged_Parameter.Descriptor instead.
func (SessionStateChanged_Parameter) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{3, 0}
}

type GroupReplicationStateChanged_Type int32

const (
	GroupReplicationStateChanged_MEMBERSHIP_QUORUM_LOSS GroupReplicationStateChanged_Type = 1
	GroupReplicationStateChanged_MEMBERSHIP_VIEW_CHANGE GroupReplicationStateChanged_Type = 2
	GroupReplicationStateChanged_MEMBER_ROLE_CHANGE     GroupReplicationStateChanged_Type = 3
	GroupReplicationStateChanged_MEMBER_STATE_CHANGE    GroupReplicationStateChanged_Type = 4
)

// Enum value maps for GroupReplicationStateChanged_Type.
var (
	GroupReplicationStateChanged_Type_name = map[int32]string{
		1: "MEMBERSHIP_QUORUM_LOSS",
		2: "MEMBERSHIP_VIEW_CHANGE",
		3: "MEMBER_ROLE_CHANGE",
		4: "MEMBER_STATE_CHANGE",
	}
	GroupReplicationStateChanged_Type_value = map[string]int32{
		"MEMBERSHIP_QUORUM_LOSS": 1,
		"MEMBERSHIP_VIEW_CHANGE": 2,
		"MEMBER_ROLE_CHANGE":     3,
		"MEMBER_STATE_CHANGE":    4,
	}
)

func (x GroupReplicationStateChanged_Type) Enum() *GroupReplicationStateChanged_Type {
	p := new(GroupReplicationStateChanged_Type)
	*p = x
	return p
}

func (x GroupReplicationStateChanged_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GroupReplicationStateChanged_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_mysqlx_notice_proto_enumTypes[4].Descriptor()
}

func (GroupReplicationStateChanged_Type) Type() protoreflect.EnumType {
	return &file_mysqlx_notice_proto_enumTypes[4]
}

func (x GroupReplicationStateChanged_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *GroupReplicationStateChanged_Type) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = GroupReplicationStateChanged_Type(num)
	return nil
}

// Deprecated: Use GroupReplicationStateChanged_Type.Descriptor instead.
func (GroupReplicationStateChanged_Type) EnumDescriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{4, 0}
}

// *
// Common frame for all notices
//
// | “.type“                                         | Value |
// |---------------------------------------------------|------ |
// | @ref Mysqlx::Notice::Warning                      | 1     |
// | @ref Mysqlx::Notice::SessionVariableChanged       | 2     |
// | @ref Mysqlx::Notice::SessionStateChanged          | 3     |
// | @ref Mysqlx::Notice::GroupReplicationStateChanged | 4     |
// | @ref Mysqlx::Notice::ServerHello                  | 5     |
type Frame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * the type of the payload
	Type *uint32 `protobuf:"varint,1,req,name=type" json:"type,omitempty"`
	// * global or local notification
	Scope *Frame_Scope `protobuf:"varint,2,opt,name=scope,enum=Mysqlx.Notice.Frame_Scope,def=1" json:"scope,omitempty"`
	// * the payload of the notification
	Payload []byte `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
}

// Default values for Frame fields.
const (
	Default_Frame_Scope = Frame_GLOBAL
)

func (x *Frame) Reset() {
	*x = Frame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Frame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Frame) ProtoMessage() {}

func (x *Frame) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Frame.ProtoReflect.Descriptor instead.
func (*Frame) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{0}
}

func (x *Frame) GetType() uint32 {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return 0
}

func (x *Frame) GetScope() Frame_Scope {
	if x != nil && x.Scope != nil {
		return *x.Scope
	}
	return Default_Frame_Scope
}

func (x *Frame) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// *
// Server-side warnings and notes
//
// @par “.scope“ == “local“
// “.level“, “.code“ and “.msg“ map the content of:
// @code{sql}
// SHOW WARNINGS
// @endcode
//
// @par “.scope“ == “global“
// (undefined) Will be used for global, unstructured messages like:
// -  server is shutting down
// -  a node disconnected from group
// -  schema or table dropped
//
// | @ref Mysqlx::Notice::Frame Field  | Value                   |
// |-----------------------------------|-------------------------|
// | “.type“                         | 1                       |
// | “.scope“                        | “local“ or “global“ |
type Warning struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * Note or Warning
	Level *Warning_Level `protobuf:"varint,1,opt,name=level,enum=Mysqlx.Notice.Warning_Level,def=2" json:"level,omitempty"`
	// * warning code
	Code *uint32 `protobuf:"varint,2,req,name=code" json:"code,omitempty"`
	// * warning message
	Msg *string `protobuf:"bytes,3,req,name=msg" json:"msg,omitempty"`
}

// Default values for Warning fields.
const (
	Default_Warning_Level = Warning_WARNING
)

func (x *Warning) Reset() {
	*x = Warning{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Warning) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Warning) ProtoMessage() {}

func (x *Warning) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Warning.ProtoReflect.Descriptor instead.
func (*Warning) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{1}
}

func (x *Warning) GetLevel() Warning_Level {
	if x != nil && x.Level != nil {
		return *x.Level
	}
	return Default_Warning_Level
}

func (x *Warning) GetCode() uint32 {
	if x != nil && x.Code != nil {
		return *x.Code
	}
	return 0
}

func (x *Warning) GetMsg() string {
	if x != nil && x.Msg != nil {
		return *x.Msg
	}
	return ""
}

// *
// Notify clients about changes to the current session variables.
//
// Every change to a variable that is accessible through:
//
// @code{sql}
// SHOW SESSION VARIABLES
// @endcode
//
// | @ref Mysqlx::Notice::Frame  Field | Value    |
// |-----------------------------------|----------|
// | “.type“                         | 2        |
// | “.scope“                        | “local“|
type SessionVariableChanged struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * name of the variable
	Param *string `protobuf:"bytes,1,req,name=param" json:"param,omitempty"`
	// * the changed value of param
	Value *mysqlxdatatypes.Scalar `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (x *SessionVariableChanged) Reset() {
	*x = SessionVariableChanged{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionVariableChanged) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionVariableChanged) ProtoMessage() {}

func (x *SessionVariableChanged) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionVariableChanged.ProtoReflect.Descriptor instead.
func (*SessionVariableChanged) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{2}
}

func (x *SessionVariableChanged) GetParam() string {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return ""
}

func (x *SessionVariableChanged) GetValue() *mysqlxdatatypes.Scalar {
	if x != nil {
		return x.Value
	}
	return nil
}

type SessionStateChanged struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * parameter key
	Param *SessionStateChanged_Parameter `protobuf:"varint,1,req,name=param,enum=Mysqlx.Notice.SessionStateChanged_Parameter" json:"param,omitempty"`
	// * updated value
	Value []*mysqlxdatatypes.Scalar `protobuf:"bytes,2,rep,name=value" json:"value,omitempty"`
}

func (x *SessionStateChanged) Reset() {
	*x = SessionStateChanged{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionStateChanged) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionStateChanged) ProtoMessage() {}

func (x *SessionStateChanged) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionStateChanged.ProtoReflect.Descriptor instead.
func (*SessionStateChanged) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{3}
}

func (x *SessionStateChanged) GetParam() SessionStateChanged_Parameter {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return SessionStateChanged_CURRENT_SCHEMA
}

func (x *SessionStateChanged) GetValue() []*mysqlxdatatypes.Scalar {
	if x != nil {
		return x.Value
	}
	return nil
}

// *
// Notify clients about group replication state changes
//
// | @ref Mysqlx::Notice::Frame Field  | Value      |
// |-----------------------------------|------------|
// |“.type“                          | 4          |
// |“.scope“                         | “global“ |
type GroupReplicationStateChanged struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * type of group replication event
	Type *uint32 `protobuf:"varint,1,req,name=type" json:"type,omitempty"`
	// * view identifier
	ViewId *string `protobuf:"bytes,2,opt,name=view_id,json=viewId" json:"view_id,omitempty"`
}

func (x *GroupReplicationStateChanged) Reset() {
	*x = GroupReplicationStateChanged{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupReplicationStateChanged) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupReplicationStateChanged) ProtoMessage() {}

func (x *GroupReplicationStateChanged) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupReplicationStateChanged.ProtoReflect.Descriptor instead.
func (*GroupReplicationStateChanged) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{4}
}

func (x *GroupReplicationStateChanged) GetType() uint32 {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return 0
}

func (x *GroupReplicationStateChanged) GetViewId() string {
	if x != nil && x.ViewId != nil {
		return *x.ViewId
	}
	return ""
}

// *
// Notify clients about connection to X Protocol server
//
// | @ref Mysqlx::Notice::Frame Field  | Value      |
// |-----------------------------------|------------|
// |“.type“                          | 5          |
// |“.scope“                         | “global“ |
type ServerHello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerHello) Reset() {
	*x = ServerHello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mysqlx_notice_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerHello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerHello) ProtoMessage() {}

func (x *ServerHello) ProtoReflect() protoreflect.Message {
	mi := &file_mysqlx_notice_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerHello.ProtoReflect.Descriptor instead.
func (*ServerHello) Descriptor() ([]byte, []int) {
	return file_mysqlx_notice_proto_rawDescGZIP(), []int{5}
}

var File_mysqlx_notice_proto protoreflect.FileDescriptor

var file_mysqlx_notice_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x63, 0x65, 0x1a, 0x0c, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x16, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9b, 0x02, 0x0a, 0x05, 0x46,
	0x72, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x02,
	0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78,
	0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x3a, 0x06, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x52, 0x05, 0x73, 0x63, 0x6f,
	0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x1e, 0x0a, 0x05,
	0x53, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x10, 0x02, 0x22, 0x83, 0x01, 0x0a,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47,
	0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x56, 0x41,
	0x52, 0x49, 0x41, 0x42, 0x4c, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x44, 0x10, 0x02,
	0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x44, 0x10, 0x03, 0x12, 0x23, 0x0a, 0x1f, 0x47,
	0x52, 0x4f, 0x55, 0x50, 0x5f, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x44, 0x10, 0x04,
	0x12, 0x10, 0x0a, 0x0c, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x48, 0x45, 0x4c, 0x4c, 0x4f,
	0x10, 0x05, 0x3a, 0x04, 0x90, 0xea, 0x30, 0x0b, 0x22, 0x97, 0x01, 0x0a, 0x07, 0x57, 0x61, 0x72,
	0x6e, 0x69, 0x6e, 0x67, 0x12, 0x3b, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x4e, 0x6f, 0x74,
	0x69, 0x63, 0x65, 0x2e, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x3a, 0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0d, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x29, 0x0a, 0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c,
	0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41,
	0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x03, 0x22, 0x5e, 0x0a, 0x16, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x72,
	0x69, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x05, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x12, 0x2e, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0xff, 0x02, 0x0a, 0x13, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x42, 0x0a, 0x05, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x2c, 0x2e, 0x4d, 0x79, 0x73, 0x71,
	0x6c, 0x78, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x2e, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x2e,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x4d, 0x79, 0x73, 0x71, 0x6c, 0x78, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x53, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xf3,
	0x01, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x0e,
	0x43, 0x55, 0x52, 0x52, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x43, 0x48, 0x45, 0x4d, 0x41, 0x10, 0x01,
	0x12, 0x13, 0x0a, 0x0f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x45, 0x58, 0x50, 0x49,
	0x52, 0x45, 0x44, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x54,
	0x45, 0x44, 0x5f, 0x49, 0x4e, 0x53, 0x45, 0x52, 0x54, 0x5f, 0x49, 0x44, 0x10, 0x03, 0x12, 0x11,
	0x0a, 0x0d, 0x52, 0x4f, 0x57, 0x53, 0x5f, 0x41, 0x46, 0x46, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10,
	0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x4f, 0x57, 0x53, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10,
	0x05, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x4f, 0x57, 0x53, 0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x45,
	0x44, 0x10, 0x06, 0x12, 0x11, 0x0a, 0x0d, 0x54, 0x52, 0x58, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x49,
	0x54, 0x54, 0x45, 0x44, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x52, 0x58, 0x5f, 0x52, 0x4f,
	0x4c, 0x4c, 0x45, 0x44, 0x42, 0x41, 0x43, 0x4b, 0x10, 0x09, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x52,
	0x4f, 0x44, 0x55, 0x43, 0x45, 0x44, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x0a,
	0x12, 0x16, 0x0a, 0x12, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x5f, 0x49, 0x44, 0x5f, 0x41, 0x53,
	0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x0b, 0x12, 0x1a, 0x0a, 0x16, 0x47, 0x45, 0x4e, 0x45,
	0x52, 0x41, 0x54, 0x45, 0x44, 0x5f, 0x44, 0x4f, 0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x49,
	0x44, 0x53, 0x10, 0x0c, 0x22, 0xbc, 0x01, 0x0a, 0x1c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x76, 0x69, 0x65,
	0x77, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x69, 0x65, 0x77,
	0x49, 0x64, 0x22, 0x6f, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x4d, 0x45,
	0x4d, 0x42, 0x45, 0x52, 0x53, 0x48, 0x49, 0x50, 0x5f, 0x51, 0x55, 0x4f, 0x52, 0x55, 0x4d, 0x5f,
	0x4c, 0x4f, 0x53, 0x53, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52,
	0x53, 0x48, 0x49, 0x50, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45,
	0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c,
	0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x45,
	0x4d, 0x42, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47,
	0x45, 0x10, 0x04, 0x22, 0x0d, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x42, 0x19, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x2e,
	0x63, 0x6a, 0x2e, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
}

var (
	file_mysqlx_notice_proto_rawDescOnce sync.Once
	file_mysqlx_notice_proto_rawDescData = file_mysqlx_notice_proto_rawDesc
)

func file_mysqlx_notice_proto_rawDescGZIP() []byte {
	file_mysqlx_notice_proto_rawDescOnce.Do(func() {
		file_mysqlx_notice_proto_rawDescData = protoimpl.X.CompressGZIP(file_mysqlx_notice_proto_rawDescData)
	})
	return file_mysqlx_notice_proto_rawDescData
}

var file_mysqlx_notice_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_mysqlx_notice_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_mysqlx_notice_proto_goTypes = []interface{}{
	(Frame_Scope)(0),                       // 0: Mysqlx.Notice.Frame.Scope
	(Frame_Type)(0),                        // 1: Mysqlx.Notice.Frame.Type
	(Warning_Level)(0),                     // 2: Mysqlx.Notice.Warning.Level
	(SessionStateChanged_Parameter)(0),     // 3: Mysqlx.Notice.SessionStateChanged.Parameter
	(GroupReplicationStateChanged_Type)(0), // 4: Mysqlx.Notice.GroupReplicationStateChanged.Type
	(*Frame)(nil),                          // 5: Mysqlx.Notice.Frame
	(*Warning)(nil),                        // 6: Mysqlx.Notice.Warning
	(*SessionVariableChanged)(nil),         // 7: Mysqlx.Notice.SessionVariableChanged
	(*SessionStateChanged)(nil),            // 8: Mysqlx.Notice.SessionStateChanged
	(*GroupReplicationStateChanged)(nil),   // 9: Mysqlx.Notice.GroupReplicationStateChanged
	(*ServerHello)(nil),                    // 10: Mysqlx.Notice.ServerHello
	(*mysqlxdatatypes.Scalar)(nil),         // 11: Mysqlx.Datatypes.Scalar
}
var file_mysqlx_notice_proto_depIdxs = []int32{
	0,  // 0: Mysqlx.Notice.Frame.scope:type_name -> Mysqlx.Notice.Frame.Scope
	2,  // 1: Mysqlx.Notice.Warning.level:type_name -> Mysqlx.Notice.Warning.Level
	11, // 2: Mysqlx.Notice.SessionVariableChanged.value:type_name -> Mysqlx.Datatypes.Scalar
	3,  // 3: Mysqlx.Notice.SessionStateChanged.param:type_name -> Mysqlx.Notice.SessionStateChanged.Parameter
	11, // 4: Mysqlx.Notice.SessionStateChanged.value:type_name -> Mysqlx.Datatypes.Scalar
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_mysqlx_notice_proto_init() }
func file_mysqlx_notice_proto_init() {
	if File_mysqlx_notice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mysqlx_notice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Frame); i {
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
		file_mysqlx_notice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Warning); i {
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
		file_mysqlx_notice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionVariableChanged); i {
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
		file_mysqlx_notice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionStateChanged); i {
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
		file_mysqlx_notice_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupReplicationStateChanged); i {
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
		file_mysqlx_notice_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerHello); i {
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
			RawDescriptor: file_mysqlx_notice_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mysqlx_notice_proto_goTypes,
		DependencyIndexes: file_mysqlx_notice_proto_depIdxs,
		EnumInfos:         file_mysqlx_notice_proto_enumTypes,
		MessageInfos:      file_mysqlx_notice_proto_msgTypes,
	}.Build()
	File_mysqlx_notice_proto = out.File
	file_mysqlx_notice_proto_rawDesc = nil
	file_mysqlx_notice_proto_goTypes = nil
	file_mysqlx_notice_proto_depIdxs = nil
}
