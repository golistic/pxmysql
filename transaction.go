// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/golistic/pxmysql/xmysql"
)

type Transaction struct {
	session *xmysql.Session
}

var _ driver.Tx = &Transaction{}

func (tx *Transaction) Commit() error {
	if tx.session == nil {
		return fmt.Errorf("not connected (%w)", driver.ErrBadConn)
	}

	if _, err := tx.session.ExecuteStatement(context.Background(), "COMMIT"); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) Rollback() error {
	if tx.session == nil {
		return fmt.Errorf("not connected (%w)", driver.ErrBadConn)
	}

	if _, err := tx.session.ExecuteStatement(context.Background(), "ROLLBACK"); err != nil {
		return err
	}

	return nil
}
