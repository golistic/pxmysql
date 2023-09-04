// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

type Department struct {
	Name string `json:"name"`
}

type Person struct {
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	Department Department `json:"department,omitempty"`
}

func TestExpr(t *testing.T) {
	t.Run("object with nested object", func(t *testing.T) {
		expr := Expr(&Person{Name: "Alice", Age: 36, Department: Department{Name: "Engineering"}})
		got := expr.Object

		xt.Eq(t, "name", *got.Fld[0].Key)
		xt.Eq(t, "Alice", string(got.Fld[0].Value.Literal.VString.Value))
		xt.Eq(t, "age", *got.Fld[1].Key)
		xt.Eq(t, int64(36), *got.Fld[1].Value.Literal.VSignedInt)
		xt.Eq(t, "department", *got.Fld[2].Key)
		xt.Eq(t, "Engineering", string(got.Fld[2].Value.Object.Fld[0].Value.Literal.VString.Value))
	})

	t.Run("array of objects", func(t *testing.T) {
		exp := []*Person{
			{Name: "Alice", Age: 36},
			{Name: "Bob", Age: 34},
		}

		expr := Expr(exp)
		array := expr.Array
		xt.Eq(t, 2, len(array.Value))

		for i, v := range array.Value {
			got := v.Object.Fld[0]
			xt.Eq(t, "name", *got.Key)
			xt.Eq(t, exp[i].Name, got.Value.Literal.VString.Value)
		}
	})
}
