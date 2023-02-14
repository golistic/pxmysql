// Copyright (c) 2023, Geert JM Vanderkelen

package pxmysql

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestParseDSN(t *testing.T) {
	t.Run("parse query string", func(t *testing.T) {
		dsn := "scott:tiger@tcp(127.0.0.1:33060)/test?useTLS=true"
		exp := &DataSource{
			User:     "scott",
			Password: "tiger",
			Protocol: "tcp",
			Address:  "127.0.0.1:33060",
			Schema:   "test",
			UseTLS:   true,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("no query string provided", func(t *testing.T) {
		dsn := "scott:tiger@tcp(127.0.0.1:33060)/test"
		exp := &DataSource{
			User:     "scott",
			Password: "tiger",
			Protocol: "tcp",
			Address:  "127.0.0.1:33060",
			Schema:   "test",
			UseTLS:   false,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("no default schema with query string", func(t *testing.T) {
		dsn := "scott:tiger@tcp(127.0.0.1:33060)?useTLS=true"
		exp := &DataSource{
			User:     "scott",
			Password: "tiger",
			Protocol: "tcp",
			Address:  "127.0.0.1:33060",
			Schema:   "",
			UseTLS:   true,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("no default schema without query string", func(t *testing.T) {
		dsn := "scott:tiger@tcp(127.0.0.1:33060)"
		exp := &DataSource{
			User:     "scott",
			Password: "tiger",
			Protocol: "tcp",
			Address:  "127.0.0.1:33060",
			Schema:   "",
			UseTLS:   false,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("no password", func(t *testing.T) {
		dsn := "scott@tcp(127.0.0.1:33060)/test"
		exp := &DataSource{
			User:     "scott",
			Password: "",
			Protocol: "tcp",
			Address:  "127.0.0.1:33060",
			Schema:   "test",
			UseTLS:   false,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("unix protocol", func(t *testing.T) {
		dsn := "scott:tiger@unix(/path/to/socket)/test"
		exp := &DataSource{
			User:     "scott",
			Password: "tiger",
			Protocol: "unix",
			Address:  "/path/to/socket",
			Schema:   "test",
			UseTLS:   false,
		}

		have, err := ParseDSN(dsn)
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})
}
