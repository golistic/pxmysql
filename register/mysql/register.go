// Copyright (c) 2023, Geert JM Vanderkelen

package mysql

import (
	"database/sql"

	"github.com/golistic/pxmysql"
)

const DriverName = "mysql"

func init() {
	sql.Register(DriverName, &pxmysql.Driver{})
}
