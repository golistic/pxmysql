// Copyright (c) 2022, Geert JM Vanderkelen

package null

import "database/sql/driver"

type Nullable interface {
	Compare(any) bool
	Value() (driver.Value, error)
}

// Compare checks value against the value stored with Nullable.
// It returns true when:
// - nullable is valid and value of nullable is equal to given value,
// - or when nullable is not valid (SQL NULL) and value is nil
//
// Panics when value cannot be used with given Nullable n.
func Compare(n Nullable, value any) bool {
	return n.Compare(value)
}
