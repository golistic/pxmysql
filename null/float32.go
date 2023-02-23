// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
)

// Float32 represents a float32 (any MySQL FLOAT type) that may be NULL.
// This is similar to sql.NullFloat64, and does not implement the Scanner interface.
type Float32 struct {
	Float32 float32
	Valid   bool
}

var _ driver.Valuer = &Float32{}
var _ Nullable = &Float32{}

// Compare returns whether value compares with the nullable Float32.
// It returns:
// - true when Valid and stored Float32 is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (nf Float32) Compare(value any) bool {
	if !nf.Valid && value != nil {
		return false
	}

	if value == nil {
		return !nf.Valid
	}

	switch v := value.(type) {
	case float32:
		return v == nf.Float32
	case *float32:
		return *v == nf.Float32
	default:
		panic(fmt.Sprintf("value is of unsupported type %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (nf Float32) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float32, nil
}
