// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestFloat64_Compare(t *testing.T) {
	t.Run("Float64", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value any
			exp   bool
		}{
			{
				n:     Float64{Float64: 8.56, Valid: true},
				value: float64Ptr(8.56),
				exp:   true,
			},
			{
				n:     Float64{Float64: 19.469, Valid: true},
				value: float64Ptr(8.56), // supposed to be not 19.469
				exp:   false,
			},
			{
				n:     Float64{Float64: 8.56, Valid: true},
				value: 8.56,
				exp:   true,
			},
			{
				n:     Float64{Float64: 19.469, Valid: true},
				value: 8.56, // supposed to be not 19.469
				exp:   false,
			},
			{
				n:     Float64{Float64: 898, Valid: false},
				value: nil,
				exp:   true,
			},
			{
				n:     Float64{Float64: 898, Valid: false},
				value: 898,
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
			_ = Float64{Float64: 0, Valid: true}.Compare("str")
		})
	})
}

func TestFloat64_Value(t *testing.T) {
	pi := 3.14

	t.Run("valid", func(t *testing.T) {
		nt := Float64{Float64: pi, Valid: true}
		v, _ := nt.Value()
		d, ok := v.(float64)
		xt.Assert(t, ok, "expected float64")
		xt.Eq(t, pi, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nd := Float64{Float64: pi, Valid: false}
		v, _ := nd.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
