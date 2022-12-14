// Copyright (c) 2022, Geert JM Vanderkelen

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/geertjanvdk/xkit/xpath"
)

const (
	mysqlVersion = "8.0.28"
	baseURL      = "https://raw.githubusercontent.com/mysql/mysql-server/mysql-" + mysqlVersion + "/plugin/x/protocol/protobuf/"
)

const pkgxmysql = "github.com/golistic/pxmysql"

const nrOfFilesAtLeast = 5

var protoPath = path.Join("internal", "mysqlx")

func main() {
	if err := generate(); err != nil {
		exitWithErr(err)
	}
}

func exitWithErr(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}

func checkExecLocation() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed getting working directory (%w)", err)
	}

	needles := []string{".git", protoPath}

	for _, n := range needles {
		if !xpath.IsDir(path.Join(d, n)) {
			return "", fmt.Errorf("must execute within root of xmysql repository")
		}
	}

	return d, nil
}

func fetchFile(name string) ([]byte, error) {
	u := baseURL + name
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed opening URL downloading %s (%w)", name, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed opening URL downloading %s (HTTP status %d)", name, resp.StatusCode)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading body downloading file %s (%w)", name, err)
	}

	return body, nil
}

func protoFiles() ([]string, error) {
	fileData, err := fetchFile("source_files.cmake")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, nrOfFilesAtLeast)
	for _, l := range bytes.Split(fileData, []byte("\n")) {
		l = bytes.TrimSpace(l)
		if len(l) == 0 || l[0] == '#' ||
			!(bytes.HasPrefix(l, []byte("mysqlx")) && bytes.HasSuffix(l, []byte(".proto"))) {
			continue
		}
		files = append(files, string(l))
	}

	return files, nil
}

func downloadFile(dir, filename string) error {
	fileData, err := fetchFile(filename)
	if err != nil {
		return err
	}
	fp, err := os.OpenFile(path.Join(dir, filename), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed opening file %s (%w)", filename, err)
	}
	if _, err := fp.Write(fileData); err != nil {
		_ = fp.Close()
		return fmt.Errorf("failed writing to file %s (%w)", filename, err)
	}
	_ = fp.Close()

	return nil
}

func generate() error {
	wd, err := checkExecLocation()
	if err != nil {
		return err
	}

	protoc, err := exec.LookPath("protoc")
	if err != nil {
		return fmt.Errorf("protoc executable not available")
	}

	files, err := protoFiles()
	if err != nil {
		return err
	}

	args := []string{protoc, "--proto_path=" + protoPath,
		"--go_out=.",
		"--go_opt=paths=import",
		"--go_opt=module=github.com/golistic/pxmysql",
	}

	for _, f := range files {
		if err := downloadFile(protoPath, f); err != nil {
			return err
		}

		m := strings.Replace(f, ".proto", "", 1)
		m = strings.Replace(m, "_", "", -1)
		args = append(args, fmt.Sprintf("--go_opt=M%s=%s/%s/%s", f, pkgxmysql, protoPath, m))
	}

	args = append(args, files...)

	output := bytes.NewBuffer(nil)

	cmd := exec.Cmd{
		Dir:    wd,
		Path:   protoc,
		Args:   args,
		Stdout: output,
		Stderr: output,
	}

	err = cmd.Run()
	switch err.(type) {
	case *exec.ExitError:
		return fmt.Errorf("execution of protoc failed: %s", output.String())
	case error:
		return fmt.Errorf("could not run protoc")
	}

	for _, f := range files {
		_ = os.Remove(path.Join(protoPath, f))
	}

	infoFile := path.Join(protoPath, "info.md")
	fp, err := os.OpenFile(infoFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed opening file %s (%w)", infoFile, err)
	}
	defer func() { _ = fp.Close() }()

	_, err = fp.WriteString(fmt.Sprintf("Generated from MySQL Server %s at %s.\n", mysqlVersion,
		time.Now().UTC().Format(time.RFC3339)))
	if err != nil {
		return fmt.Errorf("failed writing to file %s (%w)", infoFile, err)
	}

	return nil
}
