// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"database/sql/driver"
	"io"

	"github.com/golistic/pxmysql/null"
	"github.com/golistic/pxmysql/xmysql"
)

type rows struct {
	xpresult *xmysql.Result

	currRowIndex int
}

var _ driver.Rows = &rows{}

// Columns returns the names of the columns.
func (r *rows) Columns() []string {
	if r.xpresult == nil {
		return nil
	}

	cols := make([]string, len(r.xpresult.Columns))

	for i, c := range r.xpresult.Columns {
		cols[i] = string(c.Name)
	}

	return cols
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.xpresult == nil || r.currRowIndex >= len(r.xpresult.Rows) {
		return io.EOF
	}

	for i, value := range r.xpresult.Rows[r.currRowIndex].Values {
		if n, ok := value.(null.Nullable); ok {
			var err error
			dest[i], err = n.Value()
			if err != nil {
				return err
			}
		} else {
			dest[i] = value
		}
	}

	r.currRowIndex++
	return nil
}
