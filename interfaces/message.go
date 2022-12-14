// Copyright (c) 2022, Geert JM Vanderkelen

package interfaces

import (
	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
)

type ServerMessager interface {
	Unmarshall(message proto.Message) error
	ServerMessageType() mysqlx.ServerMessages_Type
}
