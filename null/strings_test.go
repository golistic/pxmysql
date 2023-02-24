// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"fmt"
	"testing"

	"github.com/golistic/xt"
)

func TestStrings_Compare(t *testing.T) {
	mascots1 := []string{"Sakila", "Go gopher"}
	mascots2 := []string{"Sakila", "Duke"}
	mascots3 := []string{"Sakila"}

	var cases = []struct {
		n     Nullable
		value any
		exp   bool
	}{
		{
			n:     Strings{Strings: mascots1, Valid: true},
			value: mascots1,
			exp:   true,
		},
		{
			n:     Strings{Strings: mascots1, Valid: true},
			value: mascots2,
			exp:   false,
		},
		{
			n:     Strings{Strings: mascots1, Valid: true},
			value: &mascots1,
			exp:   true,
		},
		{
			n:     Strings{Strings: mascots1, Valid: true},
			value: &mascots3,
			exp:   false,
		},
		{
			n:     Strings{Strings: nil, Valid: false},
			value: nil,
			exp:   true,
		},
		{
			n:     Strings{Strings: []string{}, Valid: false},
			value: nil,
			exp:   true,
		},
		{
			n:     Strings{Strings: []string{}, Valid: false},
			value: []string{""},
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
			_ = Strings{Strings: nil, Valid: true}.Compare("str")
		})
	})
}

func TestStrings_Value(t *testing.T) {
	data := []string{"Sakila", "Go gopher"}

	t.Run("valid", func(t *testing.T) {
		ns := Strings{Strings: data, Valid: true}
		v, _ := ns.Value()
		d, ok := v.([]string)
		xt.Assert(t, ok, fmt.Sprintf("expected []string; got %T", v))
		xt.Eq(t, data, d)
	})

	t.Run("not valid", func(t *testing.T) {
		ns := Strings{Strings: data, Valid: false}
		v, _ := ns.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
