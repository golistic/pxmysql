// Copyright (c) 2022, Geert JM Vanderkelen

package pxmysql_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golistic/xgo/xt"

	_ "github.com/golistic/pxmysql/register" // registers pxmysql

	"github.com/golistic/pxmysql/internal/xxt"
)

var (
	testExitCode        int
	testErr             error
	testDockerContainer string
	testSchema          = "pxmysqldriver_tests"
	testContext         *xxt.TestContext
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

	if err := testContext.Server.Container.CopyFileFromContainer(
		"/etc/mysql/conf.d/ca.pem", "_testdata/mysql_ca.pem"); err != nil {
		testErr = fmt.Errorf("failed copying MySQL CA certificate from container %s (%s)",
			testDockerContainer, err)
		return
	}

	testExitCode = m.Run()
}

func getCredentials(credentials ...string) (string, string) {
	username := "root"
	password := testContext.MySQLRootPwd
	if len(credentials) > 0 && credentials[0] != "" {
		username = credentials[0]
	}
	if len(credentials) > 1 && credentials[1] != "" {
		password = credentials[1]
	}

	return username, password
}

func getTCPDSN(credentials ...string) string {
	username, password := getCredentials(credentials...)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?useTLS=yes", username, password, testContext.XPluginAddr,
		testSchema)
}

func cnxType(t *testing.T, db *sql.DB) string {
	t.Helper()

	q := "SELECT IF(HOST='localhost', 'unix', 'tcp') As CnxType " +
		"FROM performance_schema.processlist WHERE ID = CONNECTION_ID()"

	var ct string
	err := db.QueryRow(q).Scan(&ct)
	xt.OK(t, err)

	return ct
}
