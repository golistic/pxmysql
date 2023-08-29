// Copyright (c) 2023, Geert JM Vanderkelen

package statements_test

import (
	"testing"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/xmysql/internal/statements"
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
			xt.Eq(t, c.exp, statements.PlaceholderIndexes(statements.Placeholder, c.stmt))
		})
	}
}
