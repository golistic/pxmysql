// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestUint64_Compare(t *testing.T) {
	t.Run("Uint64", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value any
			exp   bool
		}{
			{
				n:     Uint64{Uint64: 8, Valid: true},
				value: uint64Ptr(8),
				exp:   true,
			},
			{
				n:     Uint64{Uint64: 19, Valid: true},
				value: uint64Ptr(8),
				exp:   false,
			},
			{
				n:     Uint64{Uint64: 8, Valid: true},
				value: uint(8),
				exp:   true,
			},
			{
				n:     Uint64{Uint64: 19, Valid: true},
				value: uint(8),
				exp:   false,
			},
			{
				n:     Uint64{Uint64: 898, Valid: false},
				value: nil,
				exp:   true,
			},
			{
				n:     Uint64{Uint64: 898, Valid: false},
				value: uint(898),
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
			_ = Uint64{Uint64: 0, Valid: true}.Compare("str")
		})
	})
}

func TestUint64_Value(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ni := Uint64{Uint64: 9, Valid: true}
		v, _ := ni.Value()
		d, ok := v.(uint64)
		xt.Assert(t, ok, "expected uint64")
		xt.Eq(t, 9, d)
	})

	t.Run("not valid", func(t *testing.T) {
		ni := Uint64{Uint64: 0, Valid: false}
		v, _ := ni.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
