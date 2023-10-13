// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"

	"github.com/golistic/pxmysql/decimal"
)

// Decimal represents a decimal.Decimal (MySQL DECIMAL type) that may be NULL.
// This is similar to types provided by Go's sql package, and does not implement the Scanner interface.
type Decimal struct {
	Decimal decimal.Decimal
	Valid   bool
}

var _ driver.Valuer = &Decimal{}
var _ Nullable = &Decimal{}

// Compare returns whether value compares with the nullable Decimal.
// It returns:
// - true when Valid and stored Decimal is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (nd Decimal) Compare(value any) bool {
	if !nd.Valid && value != nil {
		return false
	}

	if value == nil {
		return !nd.Valid
	}

	switch v := value.(type) {
	case decimal.Decimal:
		return nd.Decimal.Equal(v)
	case *decimal.Decimal:
		return nd.Decimal.Equal(*v)
	default:
		panic(fmt.Sprintf("value is of unsupported type %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (nd Decimal) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.Decimal, nil
}
