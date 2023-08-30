// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"github.com/golistic/xgo/xstrings"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
)

type ObjectFields = []*mysqlxdatatypes.Object_ObjectField

func ObjectField[T ~string | ~bool | ~[]*mysqlxdatatypes.Object_ObjectField](key string, value T) *mysqlxdatatypes.Object_ObjectField {
	f := &mysqlxdatatypes.Object_ObjectField{
		Key: xstrings.Pointer(key),
	}

	switch v := any(value).(type) {
	case string:
		f.Value = String(v)
	case bool:
		f.Value = Bool(v)
	case []*mysqlxdatatypes.Object_ObjectField:
		f.Value = &mysqlxdatatypes.Any{
			Type: mysqlxdatatypes.Any_OBJECT.Enum(),
			Obj: &mysqlxdatatypes.Object{
				Fld: v,
			},
		}
	default:
		panic("unsupported value type")
	}

	return f
}
