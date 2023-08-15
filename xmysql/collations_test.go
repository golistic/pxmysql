// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestIsSupportedCollation(t *testing.T) {
	t.Run("using MySQL ID", func(t *testing.T) {
		xt.Assert(t, IsSupportedCollation(241))
		xt.Assert(t, IsSupportedCollation(uint64(241)))
	})

	t.Run("using MySQL name", func(t *testing.T) {
		xt.Assert(t, IsSupportedCollation("utf8mb4_esperanto_ci"))
	})
}
