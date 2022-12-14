// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"fmt"
	"os"
	"regexp"
)

var reVersion = regexp.MustCompile(`(\d+)\.(\d+).(\d+)`)

const minMySQLVersion = 8000028
const minMySQLVersionStr = "8.0.28"

// container names are defined in _support/pxmysql-compose/docker-compose.yml
const (
	defaultDockerContainer   = "pxmysql.test.db"
	defaultDockerContainerGo = "pxmysql.test.go"
)

const (
	defaultDockerExec         = "docker"
	defaultDockerXPluginAddr  = "127.0.0.1:53360"
	defaultDockerMySQLAddr    = "127.0.0.1:53306"
	defaultDockerMySQLRootPwd = "rootpwd"
)

const (
	AuthPluginNative     = "mysql_native_password"
	AuthPluginCachedSha2 = "caching_sha2_password"
)

type TestContext struct {
	MySQLRootPwd string
	XPluginAddr  string
	MySQLAddr    string
	Server       *MySQLServer
	Builder      *GoBuilder
}

func New(schema string) (*TestContext, error) {
	dbContainerName := defaultDockerContainer
	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_CONTAINER"); have {
		dbContainerName = v
	}

	goContainerName := defaultDockerContainerGo
	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_CONTAINER_GO"); have {
		goContainerName = v
	}

	var err error
	tctx := &TestContext{
		MySQLRootPwd: defaultDockerMySQLRootPwd,
		XPluginAddr:  defaultDockerXPluginAddr,
		MySQLAddr:    defaultDockerMySQLAddr,
	}

	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_XPLUGIN_ADDR"); have {
		tctx.XPluginAddr = v
	}

	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_MYSQL_ADDR"); have {
		tctx.MySQLAddr = v
	}

	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_MYSQL_PWD"); have {
		tctx.MySQLRootPwd = v
	}

	dockerExec := defaultDockerExec
	if v, have := os.LookupEnv("PXMYSQL_TEST_DOCKER_EXEC"); have {
		dockerExec = v
	}

	dbContainer, err := NewContainer(dbContainerName, dockerExec)
	if err != nil {
		return nil, err
	}

	if err := dbContainer.CheckRunning(); err != nil {
		return nil, fmt.Errorf("make sure the Docker is available (set XMYSQL_TEST_DOCKER_EXEC?)"+
			" and container %s is running (%s)", dbContainerName, err)
	}

	goContainer, err := NewContainer(goContainerName, dockerExec)
	if err != nil {
		return nil, err
	}

	if err := goContainer.CheckRunning(); err != nil {
		return nil, fmt.Errorf("make sure the Docker is available (set XMYSQL_TEST_DOCKER_EXEC?)"+
			" and container %s is running (%s)", goContainerName, err)
	}

	tctx.Server, err = NewMySQLServer(tctx, dbContainer, schema)
	if err != nil {
		return nil, err
	}

	tctx.Builder, err = NewGoBuilder(tctx, goContainer)
	if err != nil {
		return nil, err
	}

	return tctx, nil
}
