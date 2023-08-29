// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/interfaces"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxconnection"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxprepare"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsession"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsql"
	"github.com/golistic/pxmysql/mysqlerrors"
)

type ServerMessage struct {
	msgType int
	payload []byte
}

var _ interfaces.ServerMessager = &ServerMessage{}

var maxServerMessageType int32

func init() {
	for n := range mysqlx.ServerMessages_Type_name {
		if n > maxServerMessageType {
			maxServerMessageType = n
		}
	}
}

func (m *ServerMessage) Unmarshall(into proto.Message) error {
	if err := UnmarshalPartial(m.payload, into); err != nil {
		return fmt.Errorf("failed unmarshalling server message type %s (%w)",
			mysqlx.ServerMessages_Type(m.msgType).String(), err)
	}
	return nil
}

func (m *ServerMessage) ServerMessageType() mysqlx.ServerMessages_Type {
	return mysqlx.ServerMessages_Type(m.msgType)
}

func readMessage(r io.Reader) (*ServerMessage, error) {
	var header [5]byte
	if n, err := io.ReadFull(r, header[:]); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			err = os.ErrDeadlineExceeded
		}
		if errors.Is(err, io.EOF) && n < 5 {
			return nil, fmt.Errorf("broken pipe when reading (%w)", driver.ErrBadConn)
		}
		return nil, fmt.Errorf("failed reading message header (%w)", err)
	}

	if header[4] == 0x0a || int32(header[4]) > maxServerMessageType {
		return nil, mysqlerrors.New(2007)
	}

	msg := &ServerMessage{
		msgType: int(header[4]),
	}

	msg.payload = make([]byte, binary.LittleEndian.Uint32(header[0:4])-1)
	if _, err := io.ReadFull(r, msg.payload); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			err = os.ErrDeadlineExceeded
		}
		return nil, fmt.Errorf("failed reading message payload (%w)", err)
	}

	return msg, nil
}

func clientMessageType(msg proto.Message) (mysqlx.ClientMessages_Type, error) {
	// cases ordered as ClientMessage_Type constants
	switch msg.(type) {
	case *mysqlxconnection.CapabilitiesGet:
		return mysqlx.ClientMessages_CON_CAPABILITIES_GET, nil
	case *mysqlxconnection.CapabilitiesSet:
		return mysqlx.ClientMessages_CON_CAPABILITIES_SET, nil
	case *mysqlxconnection.Close:
		return mysqlx.ClientMessages_CON_CLOSE, nil

	case *mysqlxprepare.Execute:
		return mysqlx.ClientMessages_PREPARE_EXECUTE, nil
	case *mysqlxprepare.Prepare:
		return mysqlx.ClientMessages_PREPARE_PREPARE, nil
	case *mysqlxprepare.Deallocate:
		return mysqlx.ClientMessages_PREPARE_DEALLOCATE, nil

	case *mysqlxsession.AuthenticateStart:
		return mysqlx.ClientMessages_SESS_AUTHENTICATE_START, nil
	case *mysqlxsession.AuthenticateContinue:
		return mysqlx.ClientMessages_SESS_AUTHENTICATE_CONTINUE, nil
	case *mysqlxsession.Reset:
		return mysqlx.ClientMessages_SESS_RESET, nil
	case *mysqlxsession.Close:
		return mysqlx.ClientMessages_SESS_CLOSE, nil

	case *mysqlxsql.StmtExecute:
		return mysqlx.ClientMessages_SQL_STMT_EXECUTE, nil
	default:
		return 0, fmt.Errorf("unsupported message '%T'", msg)
	}
}
