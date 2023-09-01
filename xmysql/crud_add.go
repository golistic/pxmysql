// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golistic/xgo/xstrings"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxcrud"
	mysqlxexpr "github.com/golistic/pxmysql/internal/mysqlx/mysqlxexpr"
	"github.com/golistic/pxmysql/xmysql/xproto"
)

type cruder interface {
	Execute(ctx context.Context) error
	GetError() error
}

type adder interface {
	Add(object ...any) *Add
}

type Add struct {
	collection *Collection
	values     []any
	err        error
}

var (
	_ cruder = (*Add)(nil)
	_ adder  = (*Add)(nil)
)

func NewAdd(c *Collection) *Add {

	return &Add{collection: c}
}

// Add adds object to the queue.
func (a *Add) Add(objects ...any) *Add {

	for _, object := range objects {
		rt := reflect.TypeOf(object)
		if reflect.ValueOf(object).Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		if rt.Kind() != reflect.Struct {
			a.err = fmt.Errorf("unsupported object kind %s", rt.Kind())
		}

		a.values = append(a.values, object)
	}

	return a
}

func (a *Add) Execute(ctx context.Context) error {

	errBaseMsg := "adding to collection %s (%w)"

	if a.err != nil {
		return fmt.Errorf(errBaseMsg, a.collection.name, a.err)
	}

	rows := make([]*mysqlxcrud.Insert_TypedRow, len(a.values))

	for i, v := range a.values {
		rows[i] = &mysqlxcrud.Insert_TypedRow{
			Field: []*mysqlxexpr.Expr{
				{
					Type:   mysqlxexpr.Expr_OBJECT.Enum(),
					Object: xproto.StructExpr(v),
				},
			},
		}
	}

	msg := &mysqlxcrud.Insert{
		Collection: &mysqlxcrud.Collection{
			Name:   xstrings.Pointer(a.collection.Name()),
			Schema: xstrings.Pointer(a.collection.schema.Name()),
		},
		DataModel:  mysqlxcrud.DataModel_DOCUMENT.Enum(),
		Projection: nil,
		Row:        rows,
	}

	ses := a.collection.schema.GetSession()
	if err := ses.Write(ctx, msg); err != nil {
		return fmt.Errorf(errBaseMsg, a.collection.name, err)
	}

	_, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.stmtOK
	})
	if err != nil {
		return fmt.Errorf(errBaseMsg, a.collection.name, err)
	}

	return nil
}

func (a *Add) GetError() error {

	return a.err
}
