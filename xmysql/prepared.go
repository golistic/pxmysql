// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxprepare"
	"github.com/golistic/pxmysql/xmysql/internal/scalars"
)

type Prepared struct {
	session         *Session
	result          *Result
	numPlaceholders int
}

// Execute the prepared statements replacing placeholders with args.
func (p *Prepared) Execute(ctx context.Context, args ...any) (*Result, error) {
	if p.session == nil || p.result == nil || p.result.stmtID == 0 {
		return nil, fmt.Errorf("not initialized")
	}

	pArgs := make([]*mysqlxdatatypes.Any, len(args))

	for i, arg := range args {
		var err error

		var a any
		switch v := arg.(type) {
		case driver.NamedValue:
			a = v.Value
		default:
			a = arg
		}

		// ridiculous type-switch; preventing using reflection
		switch v := a.(type) {
		case nil:
			pArgs[i] = scalars.Nil()
		case bool:
			pArgs[i] = scalars.Bool(v)
		case *bool:
			pArgs[i] = scalars.Bool(*v)
		case int:
			pArgs[i] = scalars.SignedInt(v)
		case int8:
			pArgs[i] = scalars.SignedInt(v)
		case int16:
			pArgs[i] = scalars.SignedInt(v)
		case int32:
			pArgs[i] = scalars.SignedInt(v)
		case int64:
			pArgs[i] = scalars.SignedInt(v)
		case uint:
			pArgs[i] = scalars.UnsignedInt(v)
		case uint8:
			pArgs[i] = scalars.UnsignedInt(v)
		case uint16:
			pArgs[i] = scalars.UnsignedInt(v)
		case uint32:
			pArgs[i] = scalars.UnsignedInt(v)
		case uint64:
			pArgs[i] = scalars.UnsignedInt(v)
		case *int:
			pArgs[i] = scalars.SignedInt(*v)
		case *int8:
			pArgs[i] = scalars.SignedInt(*v)
		case *int16:
			pArgs[i] = scalars.SignedInt(*v)
		case *int32:
			pArgs[i] = scalars.SignedInt(*v)
		case *int64:
			pArgs[i] = scalars.SignedInt(*v)
		case *uint:
			pArgs[i] = scalars.UnsignedInt(*v)
		case *uint8:
			pArgs[i] = scalars.UnsignedInt(*v)
		case *uint16:
			pArgs[i] = scalars.UnsignedInt(*v)
		case *uint32:
			pArgs[i] = scalars.UnsignedInt(*v)
		case *uint64:
			pArgs[i] = scalars.UnsignedInt(*v)
		case string:
			pArgs[i] = scalars.String(v)
		case *string:
			pArgs[i] = scalars.String(v)
		case []byte:
			pArgs[i] = scalars.Bytes(v)
		case float32:
			pArgs[i] = scalars.Float32(v)
		case *float32:
			pArgs[i] = scalars.Float32(*v)
		case float64:
			pArgs[i] = scalars.Float64(v)
		case *float64:
			pArgs[i] = scalars.Float64(*v)
		case decimal.Decimal:
			pArgs[i] = scalars.Decimal(v)
		case *decimal.Decimal:
			pArgs[i] = scalars.Decimal(*v)
		case time.Time:
			if pArgs[i], err = scalars.Time(v, p.session.TimeLocation().String()); err != nil {
				return nil, err
			}
		case *time.Time:
			if pArgs[i], err = scalars.Time(*v, p.session.TimeLocation().String()); err != nil {
				return nil, err
			}
		case []string:
			pArgs[i] = scalars.String(strings.Join(v, ","))
		default:
			return nil, fmt.Errorf("argument type '%T' not supported", a)
		}
	}

	if err := p.session.Write(ctx, &mysqlxprepare.Execute{
		StmtId: &p.result.stmtID,
		Args:   pArgs,
	}); err != nil {
		return nil, err
	}

	res, err := p.session.handleResult(ctx, func(r *Result) bool {
		return r.stmtOK
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Deallocate makes this prepared statement not usable any longer.
func (p *Prepared) Deallocate(ctx context.Context) error {
	return p.session.DeallocatePrepareStatement(ctx, p.result.stmtID)
}

// StatementID returns the statement ID.
func (p *Prepared) StatementID() uint32 {
	return p.result.stmtID
}

// NumPlaceholders returns the number of placeholder parameters.
func (p *Prepared) NumPlaceholders() int {
	return p.numPlaceholders
}
