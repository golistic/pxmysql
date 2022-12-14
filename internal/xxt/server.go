// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/golistic/xstrings"
)

type MySQLServer struct {
	tctx      *TestContext
	Container *Container
	Schema    string
	Version   string
}

func NewMySQLServer(tctx *TestContext, container *Container, schema string) (*MySQLServer, error) {
	server := &MySQLServer{
		Container: container,
		tctx:      tctx,
		// Schema is stored at the end
	}

	if _, err := server.ExecSQLStmt("DROP SCHEMA IF EXISTS " + schema); err != nil {
		return nil, err
	}
	if _, err := server.ExecSQLStmt("CREATE SCHEMA " + schema); err != nil {
		return nil, err
	}

	errMsg := "failed getting MySQL version running in container %s (%s)"
	if output, err := server.ExecSQLStmt("SELECT VERSION()"); err != nil {
		return nil, NewTestErr(err, errMsg, container.Name, err)
	} else {
		parts := reVersion.FindAllStringSubmatch(string(output), -1)
		if parts == nil {
			return nil, NewTestErr(err, errMsg, container.Name, "reVersion")
		}
		// simplistic way of checking the MySQL version, but works..
		maj, err := strconv.ParseInt(parts[0][1], 10, 64)
		if err != nil {
			return nil, NewTestErr(nil, errMsg, container.Name)
		}
		min, err := strconv.ParseInt(parts[0][2], 10, 64)
		if err != nil {
			return nil, NewTestErr(nil, errMsg, container.Name)

		}
		patch, err := strconv.ParseInt(parts[0][3], 10, 64)
		if err != nil {
			return nil, NewTestErr(nil, errMsg, container.Name)
		}
		v := maj*1000000 + min*1000 + patch
		if v < minMySQLVersion {
			return nil, NewTestErr(fmt.Errorf("MySQL version must be %s or greater", minMySQLVersionStr),
				errMsg, container.Name)
		}

		server.Version = fmt.Sprintf("%d.%d.%d", maj, min, patch)
	}

	server.Schema = schema

	return server, nil
}

// ExecSQLStmt executes the SQL stmt using the mysql CLI within the container.
// This is not SQL-injection safe and is only used for testing.
func (my MySQLServer) ExecSQLStmt(stmt string) ([]byte, error) {
	args := []string{
		"exec", "-i", my.Container.Name,
		"mysql", "-uroot", "-p" + my.tctx.MySQLRootPwd, "-NB", "-e", stmt,
	}

	if my.Schema != "" {
		args = append(args, []string{"-D", my.Schema}...)
	}

	return my.Container.run(args...)
}

// LoadSQLScript executes the statements from files provided as scripts
// using the mysql CLI within the container.
func (my MySQLServer) LoadSQLScript(scripts ...string) error {
	args := []string{
		"exec", "-i", my.Container.Name,
		"mysql", "-uroot", "-p" + my.tctx.MySQLRootPwd,
	}

	stderr := bytes.NewBuffer(nil)

	cmd, err := my.Container.getDockerCmd(stderr, args...)
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	var capturedErr error
	go func() {
		defer func() {
			_ = stdin.Close()
			wg.Done()
		}()

		for _, s := range scripts {
			if !strings.HasSuffix(s, ".sql") {
				s += ".sql"
			}
			p := path.Join("_testdata", s)
			sql, err := os.ReadFile(p)
			if err != nil {
				capturedErr = fmt.Errorf("failed reading SQL script %s (%s)", p, err)
				break
			}

			if _, err := io.WriteString(stdin, string(sql)); err != nil {
				capturedErr = fmt.Errorf("failed writing SQL script to STDIN %s (%s)", p, err)
				break
			}
		}
	}()

	wg.Add(1)
	if err := cmd.Run(); err != nil {
		if err := getContainerExecError(stderr); err != nil {
			capturedErr = err
		}
	}

	wg.Wait()
	return capturedErr
}

func (my MySQLServer) FlushPrivileges() error {
	args := []string{
		"exec", "-i", my.Container.Name,
		"mysqladmin", "-uroot", "-p" + my.tctx.MySQLRootPwd, "flush-privileges",
	}

	_, err := my.Container.run(args...)
	return err
}

func (my MySQLServer) Variable(scope, variable string) (string, error) {
	if !(scope == "global" || scope == "session") {
		panic("scope must be one of 'session' or 'global'")
	}

	output, err := my.ExecSQLStmt(fmt.Sprintf("SELECT @@%s.%s", scope, variable))
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (my MySQLServer) CreateUser(username, password, schema, authPlugin string) error {
	// this function is not SQL-injection-safe; only used for testing

	authPlugins := []string{AuthPluginNative, AuthPluginCachedSha2}

	if !xstrings.SliceHas(authPlugins, authPlugin) {
		panic("unsupported authMethod")
	}

	createUser := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED WITH %s BY '%s'",
		username, authPlugin, password)
	grant := fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%%'", schema, username)

	if _, err := my.ExecSQLStmt(createUser); err != nil {
		return err
	}
	if _, err := my.ExecSQLStmt(grant); err != nil {
		return err
	}

	return nil
}

func (my MySQLServer) DropUser(username string) error {
	// this function is not SQL-injection-safe; only used for testing

	dropUser := fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", username)

	if _, err := my.ExecSQLStmt(dropUser); err != nil {
		return err
	}
	return nil
}

// ExecApp runs the application within the container found at path and returns its output.
func (my MySQLServer) ExecApp(path string) ([]byte, error) {
	args := []string{
		"exec", "-i", my.Container.Name,
		path,
	}

	return my.Container.run(args...)
}
