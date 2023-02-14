// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/golistic/xconv"
	"github.com/golistic/xstrings"
)

var reDSN = regexp.MustCompile(`(.*?)(?::(.*?))?@(\w+)\((.*?)\)(?:/(\w+))?/?(\?)?(.*)?`)

type DataSource struct {
	Driver   string
	User     string
	Password string
	Protocol string
	Address  string
	Schema   string
	UseTLS   bool
}

// ParseDSN parsers the name as a data source name (DSN).
func ParseDSN(name string) (*DataSource, error) {
	errMsg := "invalid data source name (%w)"

	m := reDSN.FindAllStringSubmatch(name, -1)
	if m == nil {
		return nil, fmt.Errorf(errMsg, fmt.Errorf("could not parse"))
	}

	protocol := strings.ToLower(m[0][3])
	if !(protocol == "unix" || protocol == "tcp") {
		return nil, fmt.Errorf(errMsg, fmt.Errorf("unsupported protocol '%s'", m[0][3]))
	}

	cfg := &DataSource{
		User:     m[0][1],
		Password: m[0][2],
		Protocol: protocol,
		Address:  m[0][4],
		Schema:   m[0][5],
	}

	if xstrings.SliceHas(m[0], "?") {
		query, err := url.ParseQuery(m[0][len(m[0])-1])
		if err != nil {
			return nil, fmt.Errorf(errMsg, fmt.Errorf("could not parse query part"))
		}
		if v, have := query["useTLS"]; have {
			if cfg.UseTLS, err = xconv.ParseBool(v[0]); err != nil {
				return nil, fmt.Errorf(errMsg, fmt.Errorf("invalid value for useTLS query parameter"))
			}
		}
	}

	return cfg, nil
}

func (d *DataSource) String() string {
	n := fmt.Sprintf("%s:%s@%s(%s)", d.User, d.Password, d.Protocol, d.Address)
	if d.Schema != "" {
		n += "/" + d.Schema
	} else {
		n += "/"
	}

	var queryParts []string
	if d.UseTLS {
		queryParts = append(queryParts, "useTLS=true")
	}

	if len(queryParts) > 0 {
		n += "?" + strings.Join(queryParts, "&")
	}

	return n
}
