// Copyright (c) 2022, Geert JM Vanderkelen

package mysqlerrors

import (
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestMySLQError_Error(t *testing.T) {
	t.Run("formatted MySQL errors", func(t *testing.T) {
		myErr := &Error{
			Message:  "Table 'test.no_such_table' doesn't exist", // shamelessly copy/pasted from MySQL docs
			Code:     1146,
			SQLState: "42S02",
			Severity: 1,
		}

		xt.Eq(t, "table 'test.no_such_table' doesn't exist [1146:42S02]", myErr.Error())
	})
}

func TestMySQLWarning_Error(t *testing.T) {
	t.Run("formatted as MySQL would do", func(t *testing.T) {
		myErr := &MySQLWarning{
			Message: "Data truncated for column 'b' at row 1", // shamelessly copy/pasted from MySQL docs
			Code:    1265,
			Level:   "Warning",
		}

		xt.Eq(t, "Warning 1265: Data truncated for column 'b' at row 1", myErr.Error())
	})
}
