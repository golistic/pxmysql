// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"fmt"

	"github.com/golistic/xgo/xconv"
	"github.com/golistic/xgo/xsql"
)

// DataSource defines the configuration of the connection. It embeds xsql.DataSource
// and extends it with attributes defined in the Options.
type DataSource struct {
	xsql.DataSource

	UseTLS bool
}

// NewDataSource instantiates a DataSource using the Data Source Name (DSN).
func NewDataSource(name string) (DataSource, error) {
	xds, err := xsql.ParseDSN(name)
	if err != nil {
		return DataSource{}, err
	}

	ds := DataSource{
		DataSource: *xds,
		UseTLS:     false,
	}

	if err := ds.handleOptions(); err != nil {
		return DataSource{}, err
	}

	if err := ds.CheckValidity(); err != nil {
		return DataSource{}, fmt.Errorf("configuration not valid (%w)", err)
	}

	return ds, nil
}

// CheckValidity returns whether the DataSource has enough configuration to establish
// a connection. Needed are the address, protocol, and username.
func (ds *DataSource) CheckValidity() error {
	switch {
	case ds.Address == "":
		return fmt.Errorf("address missing")
	case ds.User == "":
		return fmt.Errorf("user missing")
	case ds.Protocol == "":
		return fmt.Errorf("protocol missing")
	default:
		return nil
	}
}

func (ds *DataSource) handleOptions() error {
	var err error
	useTLS := ds.Options.Get("useTLS")
	if useTLS != "" {
		ds.UseTLS, err = xconv.ParseBool(useTLS)
		if err != nil {
			return fmt.Errorf("invalid value for useTLS option (was %s)", useTLS)
		}
	}

	return nil
}
