// Copyright (c) 2023, Geert JM Vanderkelen

package register

import (
	"database/sql"
	"testing"

	"github.com/golistic/xgo/xstrings"
	"github.com/golistic/xgo/xt"
)

func TestDriver_Open(t *testing.T) {
	t.Run("pxmysql is registered", func(t *testing.T) {
		xt.Assert(t, xstrings.SliceHas(sql.Drivers(), "pxmysql"), "expected driver pxmysql to be registered")
	})

	t.Run("mysql is not registered", func(t *testing.T) {
		xt.Assert(t, !xstrings.SliceHas(sql.Drivers(), "mysql"), "expected driver mysql not to be registered")
	})
}
