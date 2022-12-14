// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"database/sql"
)

// ConnectConfig manages the configuration of a connection to a MySQL server.
type ConnectConfig struct {
	Address             string
	UnixSockAddr        string
	Username            string
	Password            sql.NullString
	Schema              string
	UseTLS              bool
	AuthMethod          AuthMethodType
	TLSServerCACertPath string `envVar:"PXMYSQL_CA_CERT"`
	TimeZoneName        string
}

// DefaultConnectConfig is the default configuration used if none is provided
// when a Connection is instantiated.
var DefaultConnectConfig = &ConnectConfig{
	Address:    "127.0.0.1:33060", // note that the port number is of X Plugin
	Username:   "root",
	Password:   sql.NullString{String: "", Valid: true},
	Schema:     "",
	UseTLS:     false,
	AuthMethod: AuthMethodAuto,
}

// Clone duplicates other, but leaves the password nil. The caller must
// save the password.
func (cfg *ConnectConfig) Clone() *ConnectConfig {
	return &ConnectConfig{
		Address:             cfg.Address,
		UnixSockAddr:        cfg.UnixSockAddr,
		Username:            cfg.Username,
		Password:            sql.NullString{Valid: false},
		Schema:              cfg.Schema,
		UseTLS:              cfg.UseTLS,
		AuthMethod:          cfg.AuthMethod,
		TLSServerCACertPath: cfg.TLSServerCACertPath,
		TimeZoneName:        cfg.TimeZoneName,
	}
}

// SetPassword sets the password within cfg. If no password is provided,
// the Password-field of cfg will be the SQL NULL string (sql.NullString set
// as invalid).
// Panics when p has more than 1 element.
func (cfg *ConnectConfig) SetPassword(p ...string) *ConnectConfig {
	switch len(p) {
	case 1:
		cfg.Password = sql.NullString{
			String: p[0],
			Valid:  true,
		}
	case 0:
		cfg.Password = sql.NullString{Valid: false}
	default:
		panic("accepting only 1 optional string")
	}

	return cfg
}
