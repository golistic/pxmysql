// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestContextTimeLocation(t *testing.T) {
	t.Run("no time location in context", func(t *testing.T) {
		xt.Eq(t, DefaultTimeLocation.String(), ContextTimeLocation(context.Background()).String())
	})
}
