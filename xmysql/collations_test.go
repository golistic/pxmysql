// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql_test

import (
	"testing"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/xmysql"
)

func TestIsSupportedCollation(t *testing.T) {
	t.Run("using MySQL ID", func(t *testing.T) {
		xt.Assert(t, xmysql.IsSupportedCollation(241))
		xt.Assert(t, xmysql.IsSupportedCollation(uint64(241)))
	})

	t.Run("using MySQL name", func(t *testing.T) {
		xt.Assert(t, xmysql.IsSupportedCollation("utf8mb4_esperanto_ci"))
	})
}
