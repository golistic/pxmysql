// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"bytes"
	"database/sql/driver"
	"fmt"
)

// Bytes represents a []byte (any MySQL BINARY types) that may be NULL.
// This is not available in Go's sql package, and does not implement the Scanner interface.
type Bytes struct {
	Bytes []byte
	Valid bool
}

var _ driver.Valuer = &Bytes{}
var _ Nullable = &Bytes{}

// Compare returns whether value compares with the nullable Bytes.
// It returns:
// - true when Valid and stored Bytes is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (n Bytes) Compare(value any) bool {
	if !n.Valid && value != nil {
		return false
	}

	if value == nil {
		return !n.Valid
	}

	switch v := value.(type) {
	case []byte:
		return bytes.Equal(n.Bytes, v)
	case *[]byte:
		return bytes.Equal(n.Bytes, *v)
	default:
		panic(fmt.Sprintf("value must be []byte or *[]byte; not %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (n Bytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}
