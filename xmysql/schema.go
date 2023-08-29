// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql

import (
	"fmt"
)

// Schema defines the representation of a database schema. It provides
// functionality to access the schema's contents.
type Schema struct {
	session *Session
	name    string
}

// newSchema instantiates a new Schema object using session. If name is the
// empty string, the current schema of session will be used.
func newSchema(session *Session, name string) (*Schema, error) {
	if session == nil {
		return nil, fmt.Errorf("session closed")
	}

	return &Schema{
		session: session,
		name:    name,
	}, nil
}

func (s *Schema) String() string {
	return fmt.Sprintf("<Schema:%s:%s>", s.name, s.session)
}

// Name returns the schema or database name.
func (s *Schema) Name() string {
	return s.name
}
