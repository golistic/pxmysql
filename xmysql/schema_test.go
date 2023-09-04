// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/golistic/xgo/xstrings"
	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/xmysql"
	"github.com/golistic/pxmysql/xmysql/collection"
)

func TestSchema_GetSession(t *testing.T) {

	t.Run("session of schema returned", func(t *testing.T) {

		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ctx := context.Background()

		exp, err := xmysql.GetSession(ctx, config)
		xt.OK(t, err)

		for i := 0; i < 10; i++ {
			schema, err := exp.GetSchema(ctx)
			xt.OK(t, err)

			got := schema.GetSession()
			xt.Assert(t, got != nil, "expected not nil")
			xt.Eq(t, exp, got)
		}
	})

	t.Run("no session returns nil", func(t *testing.T) {

		xt.Eq(t, nil, (&xmysql.Schema{}).GetSession())
	})
}

func TestSchema_GetCollections(t *testing.T) {

	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	t.Run("all collections", func(t *testing.T) {

		xt.OK(t, testContext.Server.LoadSQLScript("schema_collections"))

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		exp := []string{"collection_wic28skwixkd", "collection_weux73293jsnsj"}
		sort.Strings(exp)

		schema, err := ses.GetSchemaWithName(context.Background(), "pxmysql_tests")
		xt.OK(t, err)

		collections, err := schema.GetCollections(context.Background())
		xt.OK(t, err)

		xt.Assert(t, len(collections) >= len(exp), fmt.Sprintf("expected at least %d", len(exp)))

		var got []string
		for _, s := range collections {
			got = append(got, s.Name())
		}
		sort.Strings(got)

		xt.Assert(t, func(exp, got []string) bool {
			if len(exp) > len(got) {
				return false
			}
			for _, l := range exp {
				if !xstrings.SliceHas(got, l) {
					return false
				}
			}
			return true
		}(exp, got))
	})
}

func TestSchema_CreateCollection(t *testing.T) {

	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
		Schema:   "pxmysql_tests",
	}
	config.SetPassword(xxt.UserNativePwd)

	ses, err := xmysql.GetSession(context.Background(), config)
	xt.OK(t, err)

	schema, err := ses.GetSchema(context.Background())
	xt.OK(t, err)

	t.Run("name is required", func(t *testing.T) {

		ctx := context.Background()
		_, err := schema.CreateCollection(ctx, "")
		xt.KO(t, err)
		xt.Eq(t, "creating collection (invalid name)", err.Error())
	})

	t.Run("create and drop", func(t *testing.T) {

		ctx := context.Background()
		name := "ciwejkuwmidi2938x"
		c, err := schema.CreateCollection(ctx, name)
		xt.OK(t, err)
		xt.Eq(t, name, c.Name())

		t.Run("check existence", func(t *testing.T) {
			c, err := schema.GetCollection(ctx, name, collection.GetValidateExistence())
			xt.OK(t, err)
			xt.Eq(t, name, c.Name())

			t.Run("drop", func(t *testing.T) {
				err := schema.DropCollection(ctx, name)
				xt.OK(t, err)

				err = c.CheckExistence(ctx)
				xt.KO(t, err)
				xt.Assert(t, errors.Is(err, xmysql.ErrNotAvailable))
			})
		})
	})

	t.Run("reuse existing", func(t *testing.T) {

		ctx := context.Background()

		name := "eovwo28373"
		_, err := schema.CreateCollection(ctx, name)
		xt.OK(t, err)

		_, err = schema.CreateCollection(ctx, name)
		xt.KO(t, err)
		xt.Eq(t, "table 'eovwo28373' already exists [1050:42S01]", errors.Unwrap(err).Error())

		_, err = schema.CreateCollection(ctx, name, collection.CreateReuseExisting())
		xt.OK(t, err)
	})
}
