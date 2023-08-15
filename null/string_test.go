// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestString_Compare(t *testing.T) {
	var cases = []struct {
		n     Nullable
		value any
		exp   bool
	}{
		{
			n:     String{String: "Sakila", Valid: true},
			value: "Sakila",
			exp:   true,
		},
		{
			n:     String{String: "Sakila", Valid: true},
			value: "Go gopher", // supposed to not include 'Go gopher'
			exp:   false,
		},
		{
			n:     String{String: "Sakila", Valid: true},
			value: stringPtr("Sakila"),
			exp:   true,
		},
		{
			n:     String{String: "Sakila", Valid: true},
			value: stringPtr("Go gopher"), // supposed to not include 'Go gopher'
			exp:   false,
		},
		{
			n:     String{Valid: false},
			value: nil,
			exp:   true,
		},
		{
			n:     String{Valid: false},
			value: "",
			exp:   false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			xt.Eq(t, c.exp, c.n.Compare(c.value))
		})
	}

	t.Run("panics if value type is not supported", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = String{Valid: true}.Compare(123)
		})
	})
}

func TestString_Value(t *testing.T) {
	str := "String!"

	t.Run("valid", func(t *testing.T) {
		ns := String{String: str, Valid: true}
		v, _ := ns.Value()
		d, ok := v.(string)
		xt.Assert(t, ok, "expected string")
		xt.Eq(t, str, d)
	})

	t.Run("not valid", func(t *testing.T) {
		ns := String{String: str, Valid: false}
		v, _ := ns.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
