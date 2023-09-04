// Copyright (c) 2023, Geert JM Vanderkelen

package xproto_test

import (
	"fmt"
	"testing"

	"github.com/golistic/xgo/xstrings"
	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
	"github.com/golistic/pxmysql/xmysql/xproto"
)

func TestScalar(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		var nilString *string

		var cases = []struct {
			have any
			exp  *mysqlxdatatypes.Scalar
		}{
			{
				have: "gopher",
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_STRING.Enum(),
					VString: &mysqlxdatatypes.Scalar_String{
						Value: []byte("gopher"),
					},
				},
			},
			{
				have: "",
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_STRING.Enum(),
					VString: &mysqlxdatatypes.Scalar_String{
						Value: []byte(""),
					},
				},
			},
			{
				have: xstrings.Pointer("gopher"),
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_STRING.Enum(),
					VString: &mysqlxdatatypes.Scalar_String{
						Value: []byte("gopher"),
					},
				},
			},
			{
				have: nilString,
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
				},
			},
			{
				have: nil,
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
				},
			},
			{
				have: []byte("gopher"),
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_OCTETS.Enum(),
					VOctets: &mysqlxdatatypes.Scalar_Octets{
						Value: []byte("gopher"),
					},
				},
			},
			{
				have: []byte{},
				exp: &mysqlxdatatypes.Scalar{
					Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
				},
			},
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("%T", c.have), func(t *testing.T) {
				got := xproto.Scalar(c.have)
				xt.Eq(t, c.exp, got)
			})
		}
	})

}
