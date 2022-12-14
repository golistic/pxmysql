// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/geertjanvdk/xkit/xt"

	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/mysqlerrors"
)

func TestSQLDriver_Open(t *testing.T) {
	pwd := "aPassword"
	users := []string{"userfkEivks", "userFcae283"}

	for _, u := range users {
		_ = testContext.Server.DropUser(u)
		xt.OK(t, testContext.Server.CreateUser(u, pwd, testSchema, xxt.AuthPluginNative))
	}

	defer func() {
		for _, u := range users {
			_ = testContext.Server.DropUser(u)
		}
	}()

	t.Run("valid data source names", func(t *testing.T) {
		var cases = map[string]string{
			"no query": fmt.Sprintf("%s:%s@tcp(%s)/%s",
				users[0], pwd, testContext.XPluginAddr, testSchema),
			"no schema": fmt.Sprintf("%s:%s@tcp(%s)/?useTLS=true",
				users[1], pwd, testContext.XPluginAddr),
		}

		for cn, dsn := range cases {
			t.Run(cn, func(t *testing.T) {
				drv := &Driver{}
				_, err := drv.Open(dsn)
				xt.OK(t, err)
			})
		}
	})

	t.Run("retrieve LastInsertID after insert", func(t *testing.T) {
		tbl := "test_AFiek23eeF"
		q := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tbl)
		_, err := testContext.Server.ExecSQLStmt(q)
		xt.OK(t, err)

		_, err = testContext.Server.ExecSQLStmt("CREATE TABLE " + tbl +
			" (id INT AUTO_INCREMENT PRIMARY KEY, t1 INT)")
		xt.OK(t, err)

		dsn := getTCPDSN("", "")
		db, err := sql.Open("mysqlpx", dsn)

		t.Run("using Prepared Statement", func(t *testing.T) {
			xt.OK(t, err)
			q := fmt.Sprintf("INSERT INTO `%s` (`t1`) VALUES (?)", tbl)
			stmt, err := db.Prepare(q)
			xt.OK(t, err)
			res, err := stmt.Exec("45")
			xt.OK(t, err)
			lastID, err := res.LastInsertId()
			xt.OK(t, err)
			xt.Eq(t, 1, lastID)
		})

		t.Run("executing INSERT directly", func(t *testing.T) {
			xt.OK(t, err)
			q := fmt.Sprintf("INSERT INTO `%s` (`t1`) VALUES (?)", tbl)
			res, err := db.Exec(q, 46)
			xt.OK(t, err)
			lastID, err := res.LastInsertId()
			xt.OK(t, err)
			xt.Eq(t, 2, lastID)
		})
	})

	t.Run("using Unix socket", func(t *testing.T) {
		// runs app within Container; will not add to coverage
		app := "unix_socket"
		_, err := testContext.Builder.App(app)
		xt.OK(t, err)

		out, err := testContext.Server.ExecApp("/shared/builds/" + app)
		xt.OK(t, err)
		xt.Eq(t, testContext.Server.Version, string(out))
	})

	t.Run("unsupported protocol", func(t *testing.T) {
		drv := &Driver{}
		_, err := drv.Open("scott:tiger@UDP(localhost)/")
		xt.KO(t, err)
		xt.Eq(t, "unsupported protocol 'UDP'", errors.Unwrap(err).Error())
	})
}

func TestConnection_Ping(t *testing.T) {
	t.Run("using TCP", func(t *testing.T) {
		drv := &Driver{}
		db, err := drv.Open(getTCPDSN())
		xt.OK(t, err)
		xt.OK(t, db.(driver.Pinger).Ping(context.Background()))
		xt.Eq(t, "tcp", cnxType(t, db))
	})

	t.Run("using non-existing Unix socket", func(t *testing.T) {
		drv := &Driver{}
		os.TempDir()
		_, err := drv.Open("username:pwd@unix(_testdata/mysqlx.sock)/myschema")
		xt.KO(t, err)
		xt.Eq(t, mysqlerrors.ClientBadUnixSocket, err.(*mysqlerrors.Error).Code)
		var opErr *net.OpError
		xt.Assert(t, errors.As(err, &opErr))
	})

}
