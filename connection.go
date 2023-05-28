// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/golistic/pxmysql/xmysql"
)

type connection struct {
	cfg     *xmysql.ConnectConfig
	cnx     *xmysql.Connection
	session *xmysql.Session
}

var (
	_ driver.Conn           = (*connection)(nil)
	_ driver.ConnBeginTx    = (*connection)(nil)
	_ driver.Pinger         = (*connection)(nil)
	_ driver.ExecerContext  = (*connection)(nil)
	_ driver.QueryerContext = (*connection)(nil)
)

func (c *connection) Prepare(query string) (driver.Stmt, error) {

	prep, err := c.session.PrepareStatement(context.Background(), query)
	if err != nil {
		return nil, err
	}

	s := &statement{
		prepared: prep,
	}

	return s, nil
}

func (c *connection) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

func (c *connection) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	q := "START TRANSACTION"
	if opts.ReadOnly {
		q += q + " READ ONLY"
	}

	if _, err := c.session.ExecuteStatement(ctx, q); err != nil {
		return nil, err
	}

	return &Transaction{session: c.session}, nil
}

func (c *connection) Ping(ctx context.Context) error {
	if c.cnx == nil {
		return fmt.Errorf("not connected (%w)", driver.ErrBadConn)
	}

	if c.session != nil {
		_, err := c.session.SessionID(ctx)
		if err != nil {
			return fmt.Errorf("ping failed (%w)", err)
		}
	}

	var err error
	c.session, err = c.cnx.NewSession(ctx)
	if err != nil {
		return fmt.Errorf("ping failed (%w)", err)
	}

	return nil
}

func (c *connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	prep, err := c.session.PrepareStatement(context.Background(), query)
	if err != nil {
		return nil, handleError(err)
	}

	stmt := &statement{
		prepared: prep,
	}

	return stmt.ExecContext(ctx, args)
}

func (c *connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	prep, err := c.session.PrepareStatement(context.Background(), query)
	if err != nil {
		return nil, handleError(err)
	}

	stmt := &statement{
		prepared: prep,
	}

	return stmt.QueryContext(ctx, args)
}
