// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"

	"github.com/golistic/xt"
)

func TestFloat32_Compare(t *testing.T) {
	t.Run("Float32", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value any
			exp   bool
		}{
			{
				n:     Float32{Float32: 8.56, Valid: true},
				value: float32(8.56),
				exp:   true,
			},
			{
				n:     Float32{Float32: 19.469, Valid: true},
				value: float32(8.56), // supposed to be not 19.469
				exp:   false,
			},
			{
				n:     Float32{Float32: 8.56, Valid: true},
				value: float32Ptr(8.56),
				exp:   true,
			},
			{
				n:     Float32{Float32: 19.469, Valid: true},
				value: float32Ptr(8.56), // supposed to be not 19.469
				exp:   false,
			},
			{
				n:     Float32{Float32: 898, Valid: false},
				value: nil,
				exp:   true,
			},
			{
				n:     Float32{Float32: 9, Valid: false},
				value: float32(9),
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
			_ = Float32{Float32: 0, Valid: true}.Compare("str")
		})
	})
}

func TestFloat32_Value(t *testing.T) {
	pi := float32(3.14)

	t.Run("valid", func(t *testing.T) {
		nt := Float32{Float32: pi, Valid: true}
		v, _ := nt.Value()
		d, ok := v.(float32)
		xt.Assert(t, ok, "expected float32")
		xt.Eq(t, pi, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nd := Float32{Float32: pi, Valid: false}
		v, _ := nd.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
