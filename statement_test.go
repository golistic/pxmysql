// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/golistic/xt"

	"github.com/golistic/pxmysql/mysqlerrors"
)

func TestStatement_Close(t *testing.T) {
	dsn := getTCPDSN("", "")
	db, err := sql.Open("mysqlpx", dsn)
	xt.OK(t, err)

	stmt := "SELECT ?"
	prep, err := db.Prepare(stmt)
	xt.OK(t, err)

	_, err = prep.Exec(3)
	xt.OK(t, err)

	xt.OK(t, prep.Close())

	_, err = prep.Exec(3)
	xt.KO(t, err)
	xt.Eq(t, "sql: statement is closed", err.Error())
}

func testOpenQueryRowsClose() ([]string, error) {
	dsn := getTCPDSN("", "")
	db, err := sql.Open("mysqlpx", dsn)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	stmt := `SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_SCHEMA`
	prep, err := db.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer func() { _ = prep.Close() }()

	rows, err := prep.QueryContext(context.Background(), "mysql")
	if err != nil {
		return nil, err
	}

	var tablesNames []string
	for rows.Next() {
		var name sql.NullString
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		if !name.Valid {
			return nil, fmt.Errorf("found entry with null as table name")
		}
		tablesNames = append(tablesNames, name.String)
	}

	return tablesNames, nil
}

func TestStatement_ExecContext(t *testing.T) {
	t.Run("respect timeout", func(t *testing.T) {
		dsn := getTCPDSN("", "")
		db, err := sql.Open("mysqlpx", dsn)
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err = db.ExecContext(ctx, "SELECT SLEEP(5)")
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, mysqlerrors.ErrContextDeadlineExceeded), err.Error())
	})
}

func BenchmarkStatement_QueryContext(b *testing.B) {
	b.Run("fetch tables from mysql database", func(b *testing.B) {
		if _, err := testOpenQueryRowsClose(); err != nil {
			b.Error(err)
		}
	})
}

func TestStatement_QueryContext(t *testing.T) {
	t.Run("has rows in result", func(t *testing.T) {
		tableNames, err := testOpenQueryRowsClose()
		xt.OK(t, err)

		sort.Strings(tableNames)
		sum := md5.Sum([]byte(strings.Join(tableNames, " ")))
		xt.Eq(t, "859173a1b7b8ef446282e772dcd3039b", hex.EncodeToString(sum[:]))
	})

	t.Run("respect timeout", func(t *testing.T) {
		dsn := getTCPDSN("", "")
		db, err := sql.Open("mysqlpx", dsn)
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err = db.QueryContext(ctx, "SELECT SLEEP(5)")
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, mysqlerrors.ErrContextDeadlineExceeded), err.Error())
	})

	t.Run("does not return sql.ErrNoRows", func(t *testing.T) {
		db, err := sql.Open("mysqlpx", getTCPDSN("", ""))
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		stmt := `SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_SCHEMA`
		rows, err := db.QueryContext(context.Background(), stmt, "_this_does_not_exists_")
		xt.OK(t, err)
		xt.Assert(t, !rows.Next(), "expected no rows")
	})

	t.Run("QueryRowContext does return sql.ErrNoRows", func(t *testing.T) {
		db, err := sql.Open("mysqlpx", getTCPDSN("", ""))
		xt.OK(t, err)
		defer func() { _ = db.Close() }()

		var name string
		stmt := `SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_SCHEMA`
		err = db.QueryRowContext(context.Background(), stmt, "_this_does_not_exists_").Scan(&name)
		xt.KO(t, err)
		xt.Assert(t, errors.Is(err, sql.ErrNoRows))
	})
}
