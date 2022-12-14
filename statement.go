// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/golistic/pxmysql/xmysql"
)

var closeTimeout = time.Second

type statement struct {
	prepared *xmysql.Prepared
	result   *xmysql.Result
}

var (
	_ driver.Stmt             = &statement{}
	_ driver.StmtQueryContext = &statement{}
	_ driver.StmtExecContext  = &statement{}
)

func (s *statement) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()

	if err := s.prepared.Deallocate(ctx); err != nil {
		return err
	} else {
		s.prepared = nil
		s.result = nil
	}
	return nil
}

// NumInput returns the number of placeholders.
func (s *statement) NumInput() int {
	if s.prepared != nil {
		return s.prepared.NumPlaceholders()
	}

	return 0
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
// Deprecated: use ExecContext instead.
func (s *statement) Exec(args []driver.Value) (driver.Result, error) {
	named := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		named[i].Name = ""
		named[i].Ordinal = i + 1
		named[i].Value = arg
	}
	return s.ExecContext(context.Background(), named)
}

func (s *statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	execArgs := make([]any, len(args))

	for i, a := range args {
		execArgs[i] = a.Value
	}

	execResult, err := s.prepared.Execute(ctx, execArgs...)
	if err != nil {
		return nil, handleError(err)
	}

	res := &result{
		xpresult: execResult,
	}

	return res, nil
}

// Query executes a query that may return rows, such as a SELECT.
// Deprecated: use QueryContext instead.
func (s *statement) Query(args []driver.Value) (driver.Rows, error) {
	named := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		named[i].Name = ""
		named[i].Ordinal = i + 1
		named[i].Value = arg
	}
	return s.QueryContext(context.Background(), named)
}

// QueryContext executes a query that may return rows, such as a SELECT.
func (s *statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	execArgs := make([]any, len(args))

	for i, a := range args {
		execArgs[i] = a
	}

	execResult, err := s.prepared.Execute(ctx, execArgs...)
	if err != nil {
		return nil, handleError(err)
	}

	if len(execResult.Rows) == 0 {
		return nil, sql.ErrNoRows
	}

	r := &rows{
		xpresult: execResult,
	}

	return r, nil
}
