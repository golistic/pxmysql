// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"reflect"
	"strings"

	"github.com/golistic/xgo/xstrings"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxexpr"
)

func Expr[T reflect.Value | any](value T) *mysqlxexpr.Expr {

	var rv reflect.Value

	switch v := any(value).(type) {
	case nil:
		return nil
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Slice:
		return &mysqlxexpr.Expr{
			Type:  mysqlxexpr.Expr_ARRAY.Enum(),
			Array: sliceExpr(rv),
		}
	case reflect.Struct:
		return &mysqlxexpr.Expr{
			Type:   mysqlxexpr.Expr_OBJECT.Enum(),
			Object: StructExpr(rv.Interface()),
		}
	default:
		return &mysqlxexpr.Expr{
			Type:    mysqlxexpr.Expr_LITERAL.Enum(),
			Literal: Scalar(rv),
		}
	}
}

func StructExpr(object any) *mysqlxexpr.Object {

	rv := reflect.Indirect(reflect.ValueOf(object))

	rt := reflect.TypeOf(object)
	if rt.Kind() == reflect.Pointer {
		rt = reflect.TypeOf(object).Elem()
	}

	obj := &mysqlxexpr.Object{}

	for i := 0; i < rv.NumField(); i++ {
		rvf := rv.Field(i)
		if !rvf.CanInterface() {
			continue
		}
		rtf := rt.Field(i)

		name := rtf.Name
		tag := rtf.Tag.Get("json")
		if tag != "" {
			if tag == "-" || (strings.HasSuffix(tag, ",omitempty") && rvf.IsZero()) {
				continue
			}
			name = strings.Replace(tag, ",omitempty", "", -1)
		}

		obj.Fld = append(obj.Fld, &mysqlxexpr.Object_ObjectField{
			Key:   xstrings.Pointer(name),
			Value: Expr(rvf),
		})
	}

	return obj
}

func sliceExpr(value reflect.Value) *mysqlxexpr.Array {
	array := &mysqlxexpr.Array{
		Value: make([]*mysqlxexpr.Expr, value.Len()),
	}

	for i := 0; i < value.Len(); i++ {
		array.Value[i] = Expr(value.Index(i))
	}

	return array
}
