// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"fmt"
	"testing"

	"github.com/golistic/xt"

	"github.com/golistic/pxmysql/decimal"
)

func TestDecimal_Compare(t *testing.T) {
	dec1 := decimal.MustNew("8.56")
	dec2 := decimal.MustNew("19.469")

	t.Run("Decimal", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value any
			exp   bool
		}{
			{
				n:     Decimal{Decimal: *dec1, Valid: true},
				value: *dec1,
				exp:   true,
			},
			{
				n:     Decimal{Decimal: *dec2, Valid: true},
				value: *dec1,
				exp:   false,
			},
			{
				n:     Decimal{Decimal: *dec1, Valid: true},
				value: dec1,
				exp:   true,
			},
			{
				n:     Decimal{Decimal: *dec2, Valid: true},
				value: dec1,
				exp:   false,
			},
			{
				n:     Decimal{Decimal: *dec1, Valid: false},
				value: nil,
				exp:   true,
			},
			{
				n:     Decimal{Decimal: *dec1, Valid: false},
				value: dec1,
				exp:   false,
			},
		}

		for _, c := range cases {
			t.Run("", func(t *testing.T) {
				xt.Eq(t, c.exp, c.n.Compare(c.value))
			})
		}
	})

	t.Run("panics if value type is not supported", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = Decimal{Decimal: *decimal.Zero, Valid: true}.Compare("str")
		})
	})
}

func TestDecimal_Value(t *testing.T) {
	pi := decimal.MustNew("3.14")

	t.Run("valid", func(t *testing.T) {
		nd := Decimal{Decimal: *pi, Valid: true}
		v, _ := nd.Value()
		d, ok := v.(decimal.Decimal)
		xt.Assert(t, ok, fmt.Sprintf("expected decimal.Decimal; got %T", v))
		xt.Eq(t, *pi, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nd := Decimal{Decimal: *pi, Valid: false}
		v, _ := nd.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
