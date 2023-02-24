// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"

	"github.com/golistic/xconv"
)

// Uint64 represents an uint64 (any MySQL unsigned integer type) that may be NULL.
// This is not available in Go's sql package, and does not implement the Scanner interface.
type Uint64 struct {
	Uint64 uint64
	Valid  bool
}

var _ driver.Valuer = &Uint64{}
var _ Nullable = &Uint64{}

// Compare returns whether value compares with the nullable Uint64.
// It returns:
// - true when Valid and stored Uint64 is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (ni Uint64) Compare(value any) bool {
	if !ni.Valid && value != nil {
		return false
	}

	if value == nil {
		return !ni.Valid
	}

	switch v := value.(type) {
	case uint64, uint, uint8, uint16, uint32:
		return xconv.UnsignedAsUint64(v) == ni.Uint64
	case *uint64, *uint, *uint8, *uint16, *uint32:
		return *xconv.UnsignedAsUint64Ptr(v) == ni.Uint64
	default:
		panic(fmt.Sprintf("value is of unsupported type %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (ni Uint64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Uint64, nil
}
