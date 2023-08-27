// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"context"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"

	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/mysqlerrors"
)

// Write writes protobuf msg using conn to the server.
func Write(ctx context.Context, conn net.Conn, msg proto.Message, maxAllowedPacket int) error {

	if conn == nil {
		return fmt.Errorf("not connected (%w)", driver.ErrBadConn)
	}

	msgType, err := clientMessageType(msg)
	if err != nil {
		return err
	}

	deadline, _ := ctx.Deadline()
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("failed setting write deadline (%w)", err)
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed marshalling protobuf message (%w)", err)
	}

	if maxAllowedPacket > 0 && len(b) > maxAllowedPacket {
		return mysqlerrors.New(mysqlerrors.ClientNetPacketTooLarge)
	}

	var header [5]byte
	binary.LittleEndian.PutUint32(header[:], uint32(len(b))+1) // +1 is final \x00

	header[4] = byte(msgType)

	buf := &net.Buffers{header[:], b}
	_, err = buf.WriteTo(conn)
	switch {
	case errors.Is(err, syscall.EPIPE):
		return fmt.Errorf("broken pipe when writing (%w)", driver.ErrBadConn)
	case err != nil:
		return fmt.Errorf("failed sending message (%w)", err)
	}

	Trace("w", msg)

	return nil
}
