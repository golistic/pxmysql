// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql/driver"

	"github.com/golistic/pxmysql/xmysql"
)

type connector struct {
	dataSource DataSource
}

var _ driver.Connector = &connector{}

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {

	// dataSource at this point is valid
	config := &xmysql.ConnectConfig{
		UseTLS:     c.dataSource.UseTLS,
		AuthMethod: xmysql.AuthMethodAuto,
		Username:   c.dataSource.User,
		Schema:     c.dataSource.Schema,
	}
	config.SetPassword(c.dataSource.Password)

	switch c.dataSource.Protocol {
	case "unix":
		config.UnixSockAddr = c.dataSource.Address
	case "tcp":
		config.Address = c.dataSource.Address
	}

	cnx, err := xmysql.NewConnection(config)
	if err != nil {
		return nil, err
	}

	ses, err := cnx.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	return &connection{
		cfg:     config,
		cnx:     cnx,
		session: ses,
	}, nil
}

func (c connector) Driver() driver.Driver {
	return &Driver{}
}
