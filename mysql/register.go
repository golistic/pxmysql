// Copyright (c) 2023, Geert JM Vanderkelen

package mysql

import (
	"database/sql"

	"github.com/golistic/pxmysql"
)

func init() {
	sql.Register("mysql", &pxmysql.Driver{})
}
