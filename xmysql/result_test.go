// Copyright (c) 2022, 2023, Geert JM Vanderkelen

package xmysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/xmysql"
)

func TestResult_FetchRow(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	tbl := "bulk_fidiEfiS223"

	ses, err := xmysql.GetSession(context.Background(), config)
	xt.OK(t, err)
	xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

	createTable := fmt.Sprintf(
		"CREATE TABLE `%s` (id INT AUTO_INCREMENT PRIMARY KEY, c1 VARCHAR(30) NOT NULL)", tbl)

	_, err = ses.ExecuteStatement(context.Background(), fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tbl))
	xt.OK(t, err)

	_, err = ses.ExecuteStatement(context.Background(), createTable)
	xt.OK(t, err)

	nrRows := 100
	for i := 0; i < nrRows; i++ {
		_, err = ses.ExecuteStatement(context.Background(),
			fmt.Sprintf("INSERT INTO `%s` (c1) VALUES (?)", tbl), fmt.Sprintf("data%d", i+1))
		xt.OK(t, err)
	}

	t.Run("fetch", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

		mUse := xxt.NewMemoryUse()
		res, err := ses.ExecuteStatement(context.Background(),
			fmt.Sprintf("SELECT * FROM `%s` ORDER BY id", tbl))
		xt.OK(t, err)
		xt.Eq(t, nrRows, len(res.Rows))

		rowCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		for i := 1; res.Row != nil; i++ {
			id := res.Row.Values[0].(int64)
			xt.Eq(t, i, id)
			xt.Eq(t, fmt.Sprintf("data%d", i), res.Row.Values[1].(string))

			err = res.FetchRow(rowCtx)
			xt.OK(t, err)
		}
		mUse.Stop()

		// keep allocations in check (if nrRows changes, this will obviously go up)
		xt.Assert(t, mUse.DiffAlloc() < 35000)
	})
}
