// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestPlaceholderIndexes(t *testing.T) {
	var cases = []struct {
		stmt string
		exp  []int
	}{
		{
			stmt: `SELECT ?`,
			exp:  []int{7},
		},
		{
			stmt: `SELECT ?, '?', "?"`,
			exp:  []int{7},
		},
		{
			stmt: `SELECT ?, '?', "?", ?`,
			exp:  []int{7, 20},
		},

		{
			stmt: `SELECT ?, '?', "?", ?, "'?'", ?`,
			exp:  []int{7, 20, 30},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			xt.Eq(t, c.exp, placeholderIndexes(stmtPlaceholder, c.stmt))
		})
	}
}
