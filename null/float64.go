// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
)

// Float64 represents a float64 (any MySQL float/double type) that may be NULL.
// This is similar to sql.NullFloat64, and does not implement the Scanner interface.
type Float64 struct {
	Float64 float64
	Valid   bool
}

var _ driver.Valuer = &Float64{}
var _ Nullable = &Float64{}

// Compare returns whether value compares with the nullable Float64.
// It returns:
// - true when Valid and stored Float64 is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (nf Float64) Compare(value any) bool {
	if !nf.Valid && value != nil {
		return false
	}

	if value == nil {
		return nf.Valid == false
	}

	switch v := value.(type) {
	case float64:
		return v == nf.Float64
	case *float64:
		return *v == nf.Float64
	default:
		panic(fmt.Sprintf("value is of unsupported type %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (nf Float64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}
