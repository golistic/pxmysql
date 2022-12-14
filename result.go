// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/golistic/pxmysql/xmysql"
)

type result struct {
	xpresult *xmysql.Result
}

var _ driver.Result = &result{}

// LastInsertId returns the database's auto-generated ID after,
// for example, an INSERT into a table with primary key.
// If this information is not available, zero (0) is returned.
func (r result) LastInsertId() (int64, error) {
	if r.xpresult != nil {
		lid := r.xpresult.LastInsertID()
		if lid > math.MaxInt64 {
			return 0, fmt.Errorf("LastInsertID overflowed max 64-bit unsigned integer")
		}
		return int64(lid), nil
	}
	return 0, nil
}

// RowsAffected returns the number of rows affected by the query.
// If this information is not available, zero (0) is returned.
func (r result) RowsAffected() (int64, error) {
	if r.xpresult != nil {
		affected := r.xpresult.RowsAffected()
		if affected > math.MaxInt64 {
			return 0, fmt.Errorf("RowsAffected overflowed max 64-bit unsigned integer")
		}
		return int64(affected), nil
	}
	return 0, nil
}
