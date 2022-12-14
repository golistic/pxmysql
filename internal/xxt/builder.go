// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"fmt"
	"strings"
)

type GoBuilder struct {
	tctx      *TestContext
	Container *Container
}

func NewGoBuilder(tctx *TestContext, container *Container) (*GoBuilder, error) {
	gb := &GoBuilder{
		Container: container,
		tctx:      tctx,
		// Schema is stored at the end
	}

	_, err := gb.goVersion()
	if err != nil {
		return nil, fmt.Errorf("failed getting Go version (%w)", err)
	}

	return gb, nil
}

// App takes the application name which is located in the container's
// shared volume located at "/shared". The name is does not include the 'app_'
// prefix.
func (gb GoBuilder) App(name string) ([]byte, error) {
	args := []string{
		"exec", "-i", gb.Container.Name,
		"sh", "/shared/build", name,
	}

	return gb.Container.run(args...)
}

func (gb GoBuilder) goVersion() (string, error) {
	args := []string{
		"exec", "-i", gb.Container.Name,
		"go", "version",
	}

	buf, err := gb.Container.run(args...)
	if err != nil {
		return "", err
	}

	version := string(buf)
	version = strings.Replace(version, "go version go", "", 1)

	return version, nil
}
