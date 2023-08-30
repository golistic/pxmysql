// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"fmt"
	"slices"
)

type Collection struct {
	schema  *Schema
	session *Session
	name    string
}

// newCollection instantiates a new Collection object with schema.
func newCollection(schema *Schema, name string) (*Collection, error) {
	if schema == nil || schema.session == nil {
		return nil, fmt.Errorf("session closed")
	}

	if name == "" {
		return nil, fmt.Errorf("invalid name")
	}

	return &Collection{
		schema:  schema,
		session: schema.session,
		name:    name,
	}, nil
}

func (c *Collection) String() string {
	return fmt.Sprintf("<GetCollection:%s:%s>", c.name, c.schema)
}

// Name returns the collection.
func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) CheckExistence(ctx context.Context) error {
	names, err := c.schema.objectNames(ctx, ObjectCollection)
	if err != nil {
		return err
	}

	if _, ok := slices.BinarySearch(names, c.name); !ok {
		return ErrNotAvailable
	}

	return nil
}
