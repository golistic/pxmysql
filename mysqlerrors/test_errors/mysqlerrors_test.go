// Copyright (c) 2023, Geert JM Vanderkelen

package test_errors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"

	"github.com/golistic/pxmysql/xmysql"
)

func TestMySQLErrors(t *testing.T) {
	t.Run("wrapped errors", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address: "127.0.0.40",
		}

		cnx, err := xmysql.NewConnection(config)
		xt.OK(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err = cnx.NewSession(ctx)
		xt.KO(t, err)
		xt.Eq(t, "unknown MySQL server host '127.0.0.40:33060' (i/o timeout) [2005:HY000]", err.Error())
		xt.Eq(t, "i/o timeout", errors.Unwrap(err).Error())
	})
}
