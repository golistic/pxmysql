// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"crypto/x509"
	"fmt"
	"os"
	"sync"
)

var ServerCAPool *x509.CertPool
var muServerCAPool sync.RWMutex

func init() {
	ServerCAPool = x509.NewCertPool()
}

func addServerCACert(certs []byte) error {

	muServerCAPool.Lock()
	defer muServerCAPool.Unlock()

	if ok := ServerCAPool.AppendCertsFromPEM(certs); !ok {
		return fmt.Errorf("appending CA certificate to pool")
	}
	return nil
}

func AddServerCACertFromFile(filename string) error {

	certs, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading server CA certificate (%w)", err)
	}

	return addServerCACert(certs)
}
