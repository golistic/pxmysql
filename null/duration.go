// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Duration represents a time.Duration (MySQL TIME type) that may be NULL.
// This is not available in Go's sql package, and does not implement the Scanner interface.
// Note that the sql.NullTime is for timestamps (which includes dates).
type Duration struct {
	Duration time.Duration
	Valid    bool
}

var _ driver.Valuer = &Duration{}
var _ Nullable = &Duration{}

// Compare returns whether value compares with the nullable Duration.
// It returns:
// - true when Valid and stored Duration is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (nd Duration) Compare(value any) bool {
	if !nd.Valid && value != nil {
		return false
	}

	if value == nil {
		return nd.Valid == false
	}

	switch v := value.(type) {
	case time.Duration:
		return nd.Duration == v
	case *time.Duration:
		return nd.Duration == *v
	default:
		panic(fmt.Sprintf("value must be time.Duration or *time.Duration; not %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (nd Duration) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.Duration, nil
}
