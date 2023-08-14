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
	db, err := sql.Open("pxmysql", dsn)
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
		rows, err := db.QueryContext(context.Background(), q, id)
		xt.OK(t, err)
		xt.Assert(t, !rows.Next())
	})
}

func TestConnection_ExecContext(t *testing.T) {
	t.Run("respect timeout", func(t *testing.T) {
		dsn := getTCPDSN()
		db, err := sql.Open("pxmysql", dsn)
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		_, err = db.ExecContext(ctx, "SELECT SLEEP(5)")
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, mysqlerrors.ErrContextDeadlineExceeded), err.Error())
	})

	t.Run("prepared statement should close using Query", func(t *testing.T) {
		dsn := getTCPDSN()
		db, err := sql.Open("pxmysql", dsn)
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		needle := "SDFIciwkdixks"
		stmt := "/* " + needle + " */ SELECT COUNT(*) FROM performance_schema.prepared_statements_instances WHERE SQL_TEXT LIKE ?"
		needleParam := "%" + needle + "%"

		for i := 0; i < 2; i++ {
			var got int
			xt.OK(t, db.QueryRowContext(context.Background(), stmt, needleParam).Scan(&got))
			xt.Eq(t, 1, got) // 1 because the query is seeing itself
		}
	})

	t.Run("prepared statement should close using Exec", func(t *testing.T) {
		dsn := getTCPDSN()
		db, err := sql.Open("pxmysql", dsn)
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		needle := "owicIOwidols"
		stmt := "/* " + needle + " */ DELETE FROM mysql.user WHERE user = 'nobody'"
		needleParam := "%" + needle + "%"

		selectStmt := "SELECT COUNT(*) FROM performance_schema.prepared_statements_instances WHERE SQL_TEXT LIKE ?"

		for i := 0; i < 2; i++ {
			var got int
			_, err := db.ExecContext(context.Background(), stmt)
			xt.OK(t, err)
			xt.OK(t, db.QueryRowContext(context.Background(), selectStmt, needleParam).Scan(&got))
			xt.Eq(t, 0, got)
		}
	})
}

func TestConnection_QueryContext(t *testing.T) {
	dsn := getTCPDSN()
	db, err := sql.Open("pxmysql", dsn)
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

func TestConnector_Connect(t *testing.T) {
	t.Run("server closing stale connection and reconnect", func(t *testing.T) {
		dsn := getTCPDSN()
		db, err := sql.Open("pxmysql", dsn)
		xt.OK(t, err)

		_, err = db.Exec("SET @@SESSION.mysqlx_wait_timeout = 2")
		xt.OK(t, err)

		var cnxID int
		xt.OK(t, db.QueryRow("SELECT CONNECTION_ID()").Scan(&cnxID))

		var n string
		var v string
		xt.OK(t, db.QueryRow("SHOW SESSION VARIABLES LIKE 'mysqlx_wait_timeout'").Scan(&n, &v))

		time.Sleep(3 * time.Second) // server should close connection

		var cnxIDAfter int
		xt.OK(t, db.QueryRow("SELECT CONNECTION_ID()").Scan(&cnxIDAfter))

		xt.Assert(t, cnxID != cnxIDAfter)
	})
}
