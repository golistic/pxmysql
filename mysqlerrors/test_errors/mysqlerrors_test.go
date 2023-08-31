// Copyright (c) 2023, Geert JM Vanderkelen

package test_errors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/xmysql"
)

func TestMySQLErrors(t *testing.T) {
	t.Run("wrapped errors", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address: "127.0.0.40",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := xmysql.GetSession(ctx, config)
		xt.Eq(t, "unknown MySQL server host '127.0.0.40:33060' (i/o timeout) [2005:HY000]", err.Error())
		xt.Eq(t, "i/o timeout", errors.Unwrap(err).Error())
	})
}
