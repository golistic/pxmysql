// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"fmt"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxconnection"
)

// ServerCapabilities holds the capabilities returned by the server.
type ServerCapabilities struct {
	TLS            bool
	AuthMechanisms []string
}

// NewServerCapabilitiesFromMessage instantiates a new ServerCapabilities object
// using a message returned by MySQL Server's X Plugin.
func NewServerCapabilitiesFromMessage(msg *serverMessage) (*ServerCapabilities, error) {
	capabilities := &mysqlxconnection.Capabilities{}
	if err := msg.Unmarshall(capabilities); err != nil {
		return nil, fmt.Errorf("message was not mysqlxconnection.Capabilities")
	}

	sc := &ServerCapabilities{}

	for _, c := range capabilities.Capabilities {
		switch c.GetName() {
		case "tls":
			sc.TLS = c.Value.Scalar.GetVBool()
		case "authentication.mechanisms":
			sc.AuthMechanisms = []string{}
			for _, m := range c.Value.Array.Value {
				sc.AuthMechanisms = append(sc.AuthMechanisms, string(m.Scalar.GetVString().GetValue()))
			}
		}
	}

	return sc, nil
}
