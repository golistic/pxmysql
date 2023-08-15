// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/mysqlerrors"
)

func NewTestErr(err error, format string, a ...any) error {
	if err != nil {
		format += " (" + err.Error() + ")"
	}
	return fmt.Errorf(format, a...)
}

func AssertMySQLError(t *testing.T, err error, code int) {
	t.Helper()

	xt.KO(t, err)
	var errMySQL *mysqlerrors.Error
	xt.Assert(t, errors.As(err, &errMySQL))
	xt.Eq(t, code, errMySQL.Code)
}
