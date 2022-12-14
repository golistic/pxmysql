// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/golistic/pxmysql/internal/xxt"
)

var (
	testExitCode int
	testErr      error
	testSchema   = "pxmysql_tests"
	testContext  *xxt.TestContext
)

var (
	testMySQLMaxAllowedPacket = -1 // MySQL's @@max_allowed_packet
)

func testTearDown() {
	if testErr != nil {
		testExitCode = 1
		fmt.Println(testErr)
	}
}

func TestMain(m *testing.M) {
	defer func() { os.Exit(testExitCode) }()
	defer testTearDown()

	var err error
	if testContext, testErr = xxt.New(testSchema); err != nil {
		return
	}

	if err := testContext.Server.LoadSQLScript("base"); err != nil {
		testErr = fmt.Errorf("failed testing MySQL running in container %s (%s)",
			testContext.Server.Container.Name, err)
		return
	}

	if err := testContext.Server.Container.CopyFileFromContainer(
		"/etc/mysql/conf.d/ca.pem", "_testdata/mysql_ca.pem"); err != nil {
		testErr = fmt.Errorf("failed copying MySQL CA certificate from container %s (%s)",
			testContext.Server.Container.Name, err)
		return
	}

	if v, err := testContext.Server.Variable("global", "max_allowed_packet"); err != nil {
		testErr = fmt.Errorf("failed getting variable max_allowed_packet (%s)", err)
		return
	} else {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			testErr = fmt.Errorf("failed converting variable max_allowed_packet (%s)", err)
			return
		}
		testMySQLMaxAllowedPacket = int(n)
	}

	testExitCode = m.Run()
}
