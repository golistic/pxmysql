// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestInt64_Compare(t *testing.T) {
	t.Run("Int64", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value any
			exp   bool
		}{
			{
				n:     Int64{Int64: -8, Valid: true},
				value: int64Ptr(-8),
				exp:   true,
			},
			{
				n:     Int64{Int64: -19, Valid: true},
				value: int64Ptr(-8), // supposed to be not 19
				exp:   false,
			},
			{
				n:     Int64{Int64: -8, Valid: true},
				value: int64(-8),
				exp:   true,
			},
			{
				n:     Int64{Int64: -19, Valid: true},
				value: int64(-8), // supposed to be not 19
				exp:   false,
			},
			{
				n:     Int64{Int64: -898, Valid: false},
				value: nil,
				exp:   true,
			},
			{
				n:     Int64{Int64: -898, Valid: false},
				value: 123,
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
			_ = Int64{Int64: 0, Valid: true}.Compare("str")
		})
	})
}

func TestInt64_Value(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ni := Int64{Int64: 9, Valid: true}
		v, _ := ni.Value()
		d, ok := v.(int64)
		xt.Assert(t, ok, "expected int64")
		xt.Eq(t, 9, d)
	})

	t.Run("not valid", func(t *testing.T) {
		ni := Int64{Int64: 0, Valid: false}
		v, _ := ni.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
