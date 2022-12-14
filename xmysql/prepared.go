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
			pArgs[i] = nilAsScalar()
		case bool:
			pArgs[i] = boolAsScalar(v)
		case *bool:
			pArgs[i] = boolAsScalar(*v)
		case int:
			pArgs[i] = signedIntAsScalar(v)
		case int8:
			pArgs[i] = signedIntAsScalar(v)
		case int16:
			pArgs[i] = signedIntAsScalar(v)
		case int32:
			pArgs[i] = signedIntAsScalar(v)
		case int64:
			pArgs[i] = signedIntAsScalar(v)
		case uint:
			pArgs[i] = unsignedIntAsScalar(v)
		case uint8:
			pArgs[i] = unsignedIntAsScalar(v)
		case uint16:
			pArgs[i] = unsignedIntAsScalar(v)
		case uint32:
			pArgs[i] = unsignedIntAsScalar(v)
		case uint64:
			pArgs[i] = unsignedIntAsScalar(v)
		case *int:
			pArgs[i] = signedIntAsScalar(*v)
		case *int8:
			pArgs[i] = signedIntAsScalar(*v)
		case *int16:
			pArgs[i] = signedIntAsScalar(*v)
		case *int32:
			pArgs[i] = signedIntAsScalar(*v)
		case *int64:
			pArgs[i] = signedIntAsScalar(*v)
		case *uint:
			pArgs[i] = unsignedIntAsScalar(*v)
		case *uint8:
			pArgs[i] = unsignedIntAsScalar(*v)
		case *uint16:
			pArgs[i] = unsignedIntAsScalar(*v)
		case *uint32:
			pArgs[i] = unsignedIntAsScalar(*v)
		case *uint64:
			pArgs[i] = unsignedIntAsScalar(*v)
		case string:
			pArgs[i] = stringAsScalar(v)
		case *string:
			pArgs[i] = stringAsScalar(v)
		case []byte:
			pArgs[i] = byteSliceAsScalar(v)
		case float32:
			pArgs[i] = float32IntAsScalar(v)
		case *float32:
			pArgs[i] = float32IntAsScalar(*v)
		case float64:
			pArgs[i] = float64IntAsScalar(v)
		case *float64:
			pArgs[i] = float64IntAsScalar(*v)
		case decimal.Decimal:
			pArgs[i] = decimalAsScalar(v)
		case *decimal.Decimal:
			pArgs[i] = decimalAsScalar(*v)
		case time.Time:
			if pArgs[i], err = timeAsScalar(v, p.session.timeLocation.String()); err != nil {
				return nil, err
			}
		case *time.Time:
			if pArgs[i], err = timeAsScalar(*v, p.session.timeLocation.String()); err != nil {
				return nil, err
			}
		case []string:
			pArgs[i] = stringAsScalar(strings.Join(v, ","))
		default:
			return nil, fmt.Errorf("argument type '%T' not supported", a)
		}
	}

	if err := write(ctx, p.session, &mysqlxprepare.Execute{
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
	return p.session.deallocatePrepareStatement(ctx, p.result.stmtID)
}

// StatementID returns the statement ID.
func (p *Prepared) StatementID() uint32 {
	return p.result.stmtID
}

// NumPlaceholders returns the number of placeholder parameters.
func (p *Prepared) NumPlaceholders() int {
	return p.numPlaceholders
}
