// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"time"
)

type CtxKey struct{}

var CtxTimeLocation = &CtxKey{}

var defaultTimeLocation = time.UTC

// SetContextTimeLocation sets the time location used when decoding MySQL DATETIME and
// TIMESTAMP to Go `time.Time` objects. If l is nil, it is unset, and default will
// be used.
func SetContextTimeLocation(ctx context.Context, l *time.Location) context.Context {
	return context.WithValue(ctx, CtxTimeLocation, l)
}

// ContextTimeLocation retrieves the time location set in context used when decoding
// MySQL DATETIME and TIMESTAMP to Go `time.Time`. If none is defined in context,
// or a none `*time.Location` was found, the default will be returned.
func ContextTimeLocation(ctx context.Context) *time.Location {
	if v := ctx.Value(CtxTimeLocation); v != nil {
		if l, ok := v.(*time.Location); ok {
			return l
		}
	}

	return defaultTimeLocation
}
