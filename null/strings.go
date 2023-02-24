// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
	"sort"
)

// Strings represents a []string (slice of strings), for example used for MySQL ENUM
// type, that may be NULL.
// This is not available in the Go's sql package, and does not implement the Scanner interface.
type Strings struct {
	Strings []string
	Valid   bool
}

var _ driver.Valuer = &Strings{}
var _ Nullable = &Strings{}

// Compare returns whether value compares with the nullable Strings.
// It returns:
// - true when Valid and stored Strings is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (ns Strings) Compare(value any) bool {
	if !ns.Valid && value != nil {
		return false
	}

	if value == nil {
		return !ns.Valid
	}

	equal := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}

		ac := make([]string, len(a))
		copy(ac, a)
		sort.Strings(ac)

		bc := make([]string, len(b))
		copy(bc, b)
		sort.Strings(bc)

		for i, v := range ac {
			if v != bc[i] {
				return false
			}
		}
		return true
	}

	switch v := value.(type) {
	case []string:
		return equal(ns.Strings, v)
	case *[]string:
		return equal(ns.Strings, *v)
	default:
		panic(fmt.Sprintf("value must be []strings or []*strings; not %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (ns Strings) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Strings, nil
}
