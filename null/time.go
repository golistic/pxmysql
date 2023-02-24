// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time represents as string (MySQL TIMESTAMP, DATETIME, and DATE types) that may be NULL.
// This is similar to sql.NullTime, and does not implement the Scanner interface.
type Time struct {
	Time  time.Time
	Valid bool
}

var _ driver.Valuer = &Time{}
var _ Nullable = &Time{}

// Compare returns whether value compares with the nullable Time.
// It returns:
// - true when Valid and stored Time is equal to value
// - true when not Valid and value is nil
// - false in any other case
func (nd Time) Compare(value any) bool {
	if !nd.Valid && value != nil {
		return false
	}

	if value == nil {
		return !nd.Valid
	}

	switch v := value.(type) {
	case time.Time:
		return nd.Time.Equal(v)
	case *time.Time:
		return nd.Time.Equal(*v)
	default:
		panic(fmt.Sprintf("value must be time.Time or *time.Time; not %T", value))
	}
}

// Value returns the value of n and implements the driver.Valuer
// as well as Nullable interface.
func (nd Time) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.Time, nil
}
