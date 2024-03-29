// Copyright (c) 2023, Geert JM Vanderkelen

package mysql

import (
	"database/sql"
	"testing"

	"github.com/golistic/xgo/xstrings"
	"github.com/golistic/xgo/xt"
)

func TestDriver_Open(t *testing.T) {
	t.Run("pxmysql is not registered", func(t *testing.T) {
		xt.Assert(t, !xstrings.SliceHas(sql.Drivers(), "pxmysql"), "expected driver pxmysql not to be registered")
	})

	t.Run("mysql is registered", func(t *testing.T) {
		xt.Assert(t, xstrings.SliceHas(sql.Drivers(), "mysql"), "expected driver mysql to be registered")
	})
}
