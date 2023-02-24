// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golistic/xt"

	"github.com/golistic/pxmysql/mysqlerrors"
)

func TestConnection_Begin(t *testing.T) {
	dsn := getTCPDSN()
	db, err := sql.Open("mysqlpx", dsn)
	xt.OK(t, err)
	defer func() { _ = db.Close() }()

	tbl := "t29dkckiidk"

	_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tbl))
	xt.OK(t, err)
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE `%s` (id INT PRIMARY KEY, c1 INT)", tbl))
	xt.OK(t, err)

	t.Run("start transaction and commit", func(t *testing.T) {
		tx, err := db.Begin()
		xt.OK(t, err)

		id := 1
		exp := 123

		stmtInsert := fmt.Sprintf("INSERT INTO `%s` (id, c1) VALUES (?, ?)", tbl)
		result, err := tx.Exec(stmtInsert, id, exp)
		xt.OK(t, err)
		affected, err := result.RowsAffected()
		xt.OK(t, err)
		xt.Eq(t, 1, affected)

		xt.OK(t, tx.Commit())

		q := fmt.Sprintf("SELECT c1 FROM `%s` WHERE id = ?", tbl)
		var have int
		xt.OK(t, db.QueryRowContext(context.Background(), q, id).Scan(&have))
		xt.Eq(t, exp, have)
	})

	t.Run("start transaction and rollback", func(t *testing.T) {
		tx, err := db.Begin()
		xt.OK(t, err)

		id := 2
		value := 987

		stmtInsert := fmt.Sprintf("INSERT INTO `%s` (id, c1) VALUES (?, ?)", tbl)
		result, err := tx.Exec(stmtInsert, id, value)
		xt.OK(t, err)
		affected, err := result.RowsAffected()
		xt.OK(t, err)
		xt.Eq(t, 1, affected)

		xt.OK(t, tx.Rollback())

		q := fmt.Sprintf("SELECT c1 FROM `%s` WHERE id = ?", tbl)
		_, err = db.QueryContext(context.Background(), q, id)
		xt.KO(t, err)
		xt.Eq(t, err.Error(), sql.ErrNoRows.Error())
	})
}

func TestConnection_ExecContext(t *testing.T) {
	dsn := getTCPDSN()
	db, err := sql.Open("mysqlpx", dsn)
	xt.OK(t, err)
	defer func() { _ = db.Close() }()

	t.Run("respect timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := db.ExecContext(ctx, "SELECT SLEEP(5)")
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, mysqlerrors.ErrContextDeadlineExceeded), err.Error())
	})
}

func TestConnection_QueryContext(t *testing.T) {
	dsn := getTCPDSN()
	db, err := sql.Open("mysqlpx", dsn)
	xt.OK(t, err)
	defer func() { _ = db.Close() }()

	t.Run("respect timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := db.QueryContext(ctx, "SELECT SLEEP(5)")
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, mysqlerrors.ErrContextDeadlineExceeded), err.Error())
	})
}
