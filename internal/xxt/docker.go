// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Container struct {
	dockerExec string
	Name       string
}

func NewContainer(name string, dockerExec string) (*Container, error) {
	return &Container{
		dockerExec: dockerExec,
		Name:       name,
	}, nil
}

// CopyFileFromContainer copies a file from a Docker container.
func (c Container) CopyFileFromContainer(srcPath, dstPath string) error {
	args := []string{
		"cp", c.Name + ":" + srcPath, dstPath,
	}

	if _, err := c.run(args...); err != nil {
		return err
	}

	return nil
}

// getDockerCmd searches for the Docker executable in the directories named by
// the PATH environment variable.
func (c Container) getDockerCmd(output io.Writer, args ...string) (*exec.Cmd, error) {
	if c.dockerExec == "" || c.dockerExec[0] != '/' {
		dockerExec, err := exec.LookPath("docker")
		if err != nil {
			return nil, err
		}
		c.dockerExec = dockerExec
	}

	if output == nil {
		output = io.Discard
	}

	return &exec.Cmd{
		Path:   c.dockerExec,
		Args:   append([]string{c.dockerExec}, args...),
		Stdout: output,
		Stderr: output,
	}, nil
}

// run executes the docker command using provided arguments.
func (c Container) run(args ...string) ([]byte, error) {
	output := bytes.NewBuffer(nil)

	cmd, err := c.getDockerCmd(output, args...)
	if err != nil {
		return nil, err
	}

	err = cmd.Run()
	switch err.(type) {
	case *exec.ExitError:
		if err := getContainerExecError(output); err != nil {
			return nil, err
		}

	case error:
		return nil, err
	}

	buf, err := io.ReadAll(output)
	if err != nil {
		return nil, err
	}

	var res []byte
	for _, l := range strings.Split(string(buf), "\n") {
		if strings.Contains(l, "[Warning]") {
			continue
		}

		res = append(res, []byte(l)...)
	}

	return res, nil
}

// CheckRunning checks whether the container is running.
func (c Container) CheckRunning() error {
	args := []string{
		"inspect", "-f", "'{{.State.Running}}'", c.Name,
	}

	_, err := c.run(args...)
	return err
}

func getContainerExecError(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	for _, l := range strings.Split(string(buf), "\n") {
		if strings.Contains(l, "[Warning]") {
			continue
		}

		if strings.HasPrefix(l, "Error ") ||
			strings.HasPrefix(l, "error:") ||
			strings.Contains(l, "ERROR ") ||
			strings.Contains(l, "Error: ") ||
			strings.Contains(l, "[ERROR]") {
			return fmt.Errorf(l)
		}

		if strings.HasPrefix(l, "OCI runtime exec failed") {
			return fmt.Errorf(strings.Replace(l, "OCI runtime exec failed: ", "", -1))
		}
	}

	msg := bytes.Replace(buf, []byte("\n"), []byte("; "), -1)

	return fmt.Errorf(string(msg))
}
