// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/golistic/pxmysql/mysqlerrors"
)

// Read reads a message from the network connection conn.
// If no deadline is set in ctx, a default will be used.
func Read(ctx context.Context, conn net.Conn) (*ServerMessage, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Second)
	}

	if err := conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("setting read deadline (%w)", err)
	}

	msg, err := readMessage(conn)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}

		var myErr *mysqlerrors.Error
		if errors.As(err, &myErr) {
			return nil, myErr
		}

		return nil, err
	}

	Trace("r", msg)

	return msg, nil
}
