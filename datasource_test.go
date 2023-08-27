// Copyright (c) 2023, Geert JM Vanderkelen

package pxmysql

import (
	"errors"
	"testing"

	"github.com/golistic/xgo/xsql"
	"github.com/golistic/xgo/xt"
)

func TestDataSource_IsZero(t *testing.T) {
	t.Run("zero when username missing", func(t *testing.T) {
		_, err := NewDataSource(":pwd@tcp(127.0.0.1)")
		xt.KO(t, err)
		xt.Eq(t, "user missing", errors.Unwrap(err).Error())
	})

	t.Run("zero when address missing", func(t *testing.T) {
		ds := DataSource{
			DataSource: xsql.DataSource{
				Driver:   "pxmysql",
				User:     "user",
				Password: "",
				Protocol: "tcp",
				Address:  "",
				Schema:   "",
				Options:  nil,
			},
		}
		err := ds.CheckValidity()
		xt.KO(t, err)
		xt.Eq(t, "address missing", err.Error())
	})

	t.Run("protocol missing", func(t *testing.T) {
		ds := DataSource{
			DataSource: xsql.DataSource{
				Driver:   "pxmysql",
				User:     "user",
				Password: "",
				Protocol: "",
				Address:  "127.0.0.1",
				Schema:   "",
				Options:  nil,
			},
		}
		err := ds.CheckValidity()
		xt.KO(t, err)
		xt.Eq(t, "protocol missing", err.Error())
	})
}

func TestNewDataSource(t *testing.T) {
	t.Run("invalid useTLS option value", func(t *testing.T) {
		_, err := NewDataSource("user:pwd@tcp(127.0.0.1)/?useTLS=nope")
		xt.KO(t, err)
		xt.Eq(t, "invalid value for useTLS option (was nope)", err.Error())
	})
}
