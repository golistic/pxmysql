// Copyright (c) 2023, Geert JM Vanderkelen

package register

import (
	"database/sql"

	"github.com/golistic/pxmysql"
)

const DriverName = "pxmysql"

func init() {
	sql.Register(DriverName, &pxmysql.Driver{})
}
