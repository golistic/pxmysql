// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"

	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/interfaces"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxconnection"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxprepare"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsession"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsql"
	"github.com/golistic/pxmysql/mysqlerrors"
)

type serverMessage struct {
	msgType int
	payload []byte
}

var _ interfaces.ServerMessager = &serverMessage{}

var maxServerMessageType int32

func init() {
	for n := range mysqlx.ServerMessages_Type_name {
		if n > maxServerMessageType {
			maxServerMessageType = n
		}
	}
}

func (m *serverMessage) Unmarshall(into proto.Message) error {
	if err := UnmarshalPartial(m.payload, into); err != nil {
		return fmt.Errorf("failed unmarshalling server message type %s (%w)",
			mysqlx.ServerMessages_Type(m.msgType).String(), err)
	}
	return nil
}

func (m *serverMessage) ServerMessageType() mysqlx.ServerMessages_Type {
	return mysqlx.ServerMessages_Type(m.msgType)
}

func readMessage(r io.Reader) (*serverMessage, error) {
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

	msg := &serverMessage{
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

func write(ctx context.Context, session *Session, msg proto.Message) error {
	if session == nil || session.conn == nil {
		return fmt.Errorf("not connected (%w)", driver.ErrBadConn)
	}

	msgType, err := clientMessageType(msg)
	if err != nil {
		return err
	}

	deadline, _ := ctx.Deadline()
	if err := session.conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("failed setting write deadline (%w)", err)
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed marshalling protobuf message (%w)", err)
	}

	if session.maxAllowedPacket > 0 && len(b) > session.maxAllowedPacket {
		return mysqlerrors.New(mysqlerrors.ClientNetPacketTooLarge)
	}

	var header [5]byte
	binary.LittleEndian.PutUint32(header[:], uint32(len(b))+1) // +1 is final \x00

	header[4] = byte(msgType)

	buf := &net.Buffers{header[:], b}
	_, err = buf.WriteTo(session.conn)
	switch {
	case errors.Is(err, syscall.EPIPE):
		return fmt.Errorf("broken pipe when writing (%w)", driver.ErrBadConn)
	case err != nil:
		return fmt.Errorf("failed sending message (%w)", err)
	}

	trace("w", msg)

	return nil
}

func read(ctx context.Context, conn net.Conn) (*serverMessage, error) {
	deadline, _ := ctx.Deadline()

	if err := conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("failed setting read deadline (%w)", err)
	}

	msg, err := readMessage(conn)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		} else if _, ok := err.(*mysqlerrors.Error); ok {
			return nil, err
		}

		return nil, err
	}

	trace("r", msg)

	return msg, nil

}
