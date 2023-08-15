// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"fmt"
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestBytes_Compare(t *testing.T) {
	var cases = []struct {
		n     Nullable
		value any
		exp   bool
	}{
		{
			n:     Bytes{Bytes: []byte("Sakila"), Valid: true},
			value: bytesPtr([]byte("Sakila")),
			exp:   true,
		},
		{
			n:     Bytes{Bytes: []byte("Go gopher"), Valid: true},
			value: bytesPtr([]byte("Sakila")), // supposed to be not 'Go gopher'
			exp:   false,
		},
		{
			n:     Bytes{Bytes: []byte("Sakila"), Valid: true},
			value: []byte("Sakila"),
			exp:   true,
		},
		{
			n:     Bytes{Bytes: []byte("Go gopher"), Valid: true},
			value: []byte("Sakila"), // supposed to be not 'Go gopher'
			exp:   false,
		},
		{
			n:     Bytes{Bytes: nil, Valid: false},
			value: nil,
			exp:   true,
		},
		{
			n:     Bytes{Bytes: nil, Valid: false},
			value: []byte{},
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
			_ = Bytes{Bytes: nil, Valid: true}.Compare("str")
		})
	})
}

func TestBytes_Value(t *testing.T) {
	data := []byte("I am Data")

	t.Run("valid", func(t *testing.T) {
		nb := Bytes{Bytes: data, Valid: true}
		v, _ := nb.Value()
		d, ok := v.([]byte)
		xt.Assert(t, ok, fmt.Sprintf("expected []byte; got %T", v))
		xt.Eq(t, data, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nb := Bytes{Bytes: data, Valid: false}
		v, _ := nb.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
