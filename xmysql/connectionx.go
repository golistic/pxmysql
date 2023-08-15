// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/golistic/envs"

	"github.com/golistic/pxmysql/mysqlerrors"
)

const (
	defaultXPluginPort = "33060"
	defaultXPluginHost = "127.0.0.1"
)

// Connection manages the information and configuration on how to
// connect with a MySQL Server using the X Protocol (X Plugin).
type Connection struct {
	config       *ConnectConfig
	password     string
	timeLocation *time.Location
}

// NewConnection instantiates a new connection object.
// The password is not public available after the configuration has been
// stored.
func NewConnection(config *ConnectConfig) (*Connection, error) {
	cnx := &Connection{}

	if config != nil {
		if config.Password.Valid {
			cnx.password = config.Password.String
		}
		cnx.config = config.Clone()
	} else {
		cnx.config = DefaultConnectConfig.Clone()
	}

	if err := envs.OSEnviron(cnx.config); err != nil {
		return nil, fmt.Errorf("failed reading from environment (%w)", err)
	}

	if cnx.config.UnixSockAddr != "" {
		f, err := filepath.Abs(cnx.config.UnixSockAddr)
		if err == nil {
			var stat os.FileInfo
			stat, err = os.Stat(f)
			if err == nil {
				if stat.Mode().Type() != fs.ModeSocket {
					err = fmt.Errorf("not Unix domain socket")
				}
			}
		}

		if err != nil {
			return nil, mysqlerrors.New(mysqlerrors.ClientBadUnixSocket, cnx.config.UnixSockAddr, errors.Unwrap(err))
		}
		cnx.config.UnixSockAddr = f
	} else {
		h, p, err := net.SplitHostPort(cnx.config.Address)
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) {
			h = cnx.config.Address // on error h is empty
			p = defaultXPluginPort
		}
		if h == "" {
			h = defaultXPluginHost
		}
		cnx.config.Address = net.JoinHostPort(h, p)
	}

	if cnx.config.AuthMethod == "" {
		cnx.config.AuthMethod = DefaultConnectConfig.AuthMethod
	}

	if !supportedAuthMethods.Has(cnx.config.AuthMethod) {
		return nil, fmt.Errorf("unsupported authentication type '%s'", cnx.config.AuthMethod)
	}

	if cnx.config.TimeZoneName != "" {
		l, err := time.LoadLocation(cnx.config.TimeZoneName)
		if err != nil {
			return nil, fmt.Errorf("failed loading time location (%w)", err)
		}
		cnx.timeLocation = l
	} else {
		cnx.timeLocation = defaultTimeLocation
	}

	return cnx, nil
}

// NewSession instantiates a new Session which uses cnx.
func (cnx *Connection) NewSession(ctx context.Context) (*Session, error) {
	return newSession(ctx, cnx)
}

func (cnx *Connection) isReachable() bool {
	c, _ := net.DialTimeout("tcp", cnx.config.Address, time.Second)
	if c != nil {
		_ = c.Close()
		return true
	}

	return false
}

func (cnx *Connection) serverHostname() string {
	host, _, _ := net.SplitHostPort(cnx.config.Address) // error ignored; just return empty if not available
	return host
}
