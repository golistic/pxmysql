// Copyright (c) 2022, 2023, Geert JM Vanderkelen

package pxmysql

import (
	"context"
	"database/sql/driver"
)

type Driver struct{}

var (
	_ driver.Driver        = &Driver{}
	_ driver.DriverContext = &Driver{}
)

// Open returns a new connection to the MySQL database using MySQL X Protocol.
func (d *Driver) Open(name string) (driver.Conn, error) {
	c, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}

	return c.Connect(context.Background())
}

// OpenConnector returns a connector which will be used by sql.DB to open a connection
// to the MySQL database using MySQL X Protocol.
// This will be used instead of the Open-method (which actually uses this method).
func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	ds, err := NewDataSource(name)
	if err != nil {
		return nil, err
	}

	return &connector{
		dataSource: ds,
	}, nil
}
