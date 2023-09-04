// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql_test

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"testing"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/null"
	"github.com/golistic/pxmysql/xmysql"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func crudTestCollection(t *testing.T, name string) (*xmysql.Schema, *xmysql.Collection) {
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

	c, err := schema.CreateCollection(context.Background(), name)
	xt.OK(t, err)

	return schema, c
}

func TestCollection_Add(t *testing.T) {
	schema, coll := crudTestCollection(t, "person_2987dk8dj0s")

	t.Run("can only add struct", func(t *testing.T) {
		err := coll.Add(1).GetError()
		xt.KO(t, err)
		xt.Eq(t, "object must be struct", err.Error())
	})

	t.Run("object as pointer value", func(t *testing.T) {
		xt.OK(t, coll.Add(&Person{Name: "Alice"}).GetError())
	})

	t.Run("object as value", func(t *testing.T) {
		xt.OK(t, coll.Add(Person{Name: "Alice,c"}).GetError())
	})

	t.Run("execute stores data", func(t *testing.T) {
		xt.OK(t, coll.
			Add(&Person{Name: "Laurie", Age: 19}).
			Add(&Person{Name: "Nadya", Age: 54}, &Person{Name: "Lucas", Age: 32}).
			Execute(context.Background()))
		exp := []string{"Laurie", "Nadya", "Lucas"}
		sort.Strings(exp)

		ses := schema.GetSession()
		res, err := ses.ExecuteStatement(context.Background(), "SELECT doc FROM person_2987dk8dj0s")
		xt.OK(t, err)

		var got []string
		for _, row := range res.Rows {
			doc, ok := row.Values[0].(null.Bytes)
			xt.Assert(t, ok, "null.Bytes")
			p := Person{}
			xt.OK(t, json.Unmarshal(doc.Bytes, &p))
			got = append(got, p.Name)
		}
		sort.Strings(got)

		xt.Eq(t, exp, got)
	})

	t.Run("execute to return error stored by adding", func(t *testing.T) {
		adder := coll.Add(&Person{Name: "Laurie", Age: 19}).Add("something not OK")
		err := adder.Execute(context.Background())
		xt.KO(t, err)
		xt.Eq(t, "unsupported object kind string", errors.Unwrap(err).Error())
	})
}
