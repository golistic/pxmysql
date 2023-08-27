// Copyright (c) 2023, Geert JM Vanderkelen

package pxmysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golistic/xgo/xt"
)

func TestRows_Next(t *testing.T) {
	db, err := sql.Open("pxmysql", getTCPDSN("", ""))
	xt.OK(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	t.Run("time.Time", func(t *testing.T) {
		tbl := "test_data_types_null_time"
		_, err := db.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` (id int, ts DATETIME NULL)", tbl))
		xt.OK(t, err)

		_, err = db.ExecContext(ctx, fmt.Sprintf(
			"INSERT INTO `%s` (id, ts) VALUE (1, NOW()),(2, NULL)", tbl))
		xt.OK(t, err)

		stmt := fmt.Sprintf("SELECT ts FROM `%s` WHERE id = ?", tbl)

		var ts time.Time
		xt.OK(t, db.QueryRowContext(ctx, stmt, 1).Scan(&ts))

		var tsNull sql.NullTime
		xt.OK(t, db.QueryRowContext(ctx, stmt, 2).Scan(&tsNull))
	})
}
