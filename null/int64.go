// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"

	"github.com/golistic/xconv"
)

// Int64 represents an int64 (any MySQL signed integral type) that may be NULL.
// This is similar to sql.NullInt64, and does not implement the Scanner interface.
type Int64 struct {
	Int64 int64
	Valid bool
}

var _ driver.Valuer = &Int64{}
var _ Nullable = &Int64{}

// Compare returns whether value compares with the nullable Duration.
// It returns:
// - true when Valid and stored Duration is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (ni Int64) Compare(value any) bool {
	if !ni.Valid && value != nil {
		return false
	}

	if value == nil {
		return ni.Valid == false
	}

	switch v := value.(type) {
	case int64, int, int8, int16, int32:
		return xconv.SignedAsInt64(v) == ni.Int64
	case *int64, *int, *int8, *int16, *int32:
		return *xconv.SignedAsInt64Ptr(v) == ni.Int64
	default:
		panic(fmt.Sprintf("value is of unsupported type %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (ni Int64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}
