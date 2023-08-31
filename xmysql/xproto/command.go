// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"github.com/golistic/xgo/xstrings"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsql"
	"github.com/golistic/pxmysql/xmysql/internal/network"
)

func Command(command string, args *mysqlxdatatypes.Any) *mysqlxsql.StmtExecute {
	return &mysqlxsql.StmtExecute{
		Namespace: xstrings.Pointer(network.NamespaceMySQLx),
		Stmt:      []byte(command),
		Args:      []*mysqlxdatatypes.Any{args},
	}
}

func CommandArgs(fields ...*mysqlxdatatypes.Object_ObjectField) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_OBJECT.Enum(),
		Obj: &mysqlxdatatypes.Object{
			Fld: fields,
		},
	}
}
