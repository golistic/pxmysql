// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"crypto/x509"
	"fmt"
	"os"
	"sync"
)

var serverCAPool *x509.CertPool
var muServerCAPool sync.RWMutex

func init() {
	serverCAPool = x509.NewCertPool()
}

func addServerCACert(certs []byte) error {
	muServerCAPool.Lock()
	defer muServerCAPool.Unlock()

	if ok := serverCAPool.AppendCertsFromPEM(certs); !ok {
		return fmt.Errorf("failed appending CA certificate to pool")
	}
	return nil
}

func addServerCACertFromFile(filename string) error {
	certs, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed reading server CA certificate (%w)", err)
	}

	return addServerCACert(certs)
}
