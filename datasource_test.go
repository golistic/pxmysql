// Copyright (c) 2023, Geert JM Vanderkelen

package pxmysql

import (
	"strings"
	"testing"

	"github.com/golistic/xt"
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

		t.Run("using String-method, query part must be included", func(t *testing.T) {
			xt.Assert(t, strings.Contains(have.String(), "?useTLS=true"))
		})
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

		t.Run("using String-method, with useTLS false, it is not included", func(t *testing.T) {
			xt.Assert(t, !strings.Contains(have.String(), "?useTLS="))
		})
	})

	t.Run("no default schema with query string", func(t *testing.T) {
		var cases = map[string]string{
			"without slash": "scott:tiger@tcp(127.0.0.1:33060)/?useTLS=true",
			"with slash":    "scott:tiger@tcp(127.0.0.1:33060)?useTLS=true",
		}

		for name, dsn := range cases {
			t.Run(name, func(t *testing.T) {
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
		}
	})

	t.Run("no default schema without query string", func(t *testing.T) {
		var cases = map[string]string{
			"without slash": "scott:tiger@tcp(127.0.0.1:33060)/",
			"with slash":    "scott:tiger@tcp(127.0.0.1:33060)",
		}

		for name, dsn := range cases {
			t.Run(name, func(t *testing.T) {
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
		}
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
