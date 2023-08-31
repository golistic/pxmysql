// Copyright (c) 2023, Geert JM Vanderkelen

package statements_test

import (
	"strconv"
	"testing"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/xmysql/internal/statements"
)

func TestQuoteValue(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		var cases = []struct {
			got string
			exp string
		}{
			{
				got: "Gopher",
				exp: "'Gopher'",
			},
			{
				got: "'Gopher'",
				exp: `'\'Gopher\''`,
			},
			{
				got: "'poop'; DROP TABLE gophers",
				exp: `'\'poop\'; DROP TABLE gophers'`,
			},
			{
				got: "🐰",
				exp: `'🐰'`,
			},
		}

		for i, c := range cases {
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				got, err := statements.QuoteValue(c.got)
				xt.OK(t, err)
				xt.Eq(t, c.exp, got)
			})
		}
	})
}
