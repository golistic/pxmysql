// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
)

// String represents as string (any MySQL CHAR-kind of data type) that may be NULL.
// This is similar to sql.NullString, and does not implement the Scanner interface.
type String struct {
	String string
	Valid  bool
}

var _ driver.Valuer = &String{}
var _ Nullable = &String{}

// Compare returns whether value compares with the nullable String.
// It returns:
// - true when Valid and stored String is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (ns String) Compare(value any) bool {
	if !ns.Valid && value != nil {
		return false
	}

	if value == nil {
		return ns.Valid == false
	}

	switch v := value.(type) {
	case string:
		return ns.String == v
	case *string:
		return ns.String == *v
	default:
		panic(fmt.Sprintf("value must be string or *string; not %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (ns String) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}
