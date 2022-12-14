// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"errors"
	"os"

	"github.com/golistic/pxmysql/mysqlerrors"
)

func handleError(err error) error {
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return mysqlerrors.ErrContextDeadlineExceeded
	}

	return err
}
