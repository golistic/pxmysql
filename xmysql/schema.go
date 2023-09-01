// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"fmt"
	"sort"

	"github.com/golistic/pxmysql/null"
	"github.com/golistic/pxmysql/xmysql/collection"
	"github.com/golistic/pxmysql/xmysql/xproto"
)

type ObjectKind string

const (
	ObjectCollection ObjectKind = "COLLECTION"
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

// GetSession returns the underlying session of s.
func (s *Schema) GetSession() *Session {
	return s.session
}

// GetCollection retrieve the collection using its name.
// To keep compatible with behavior seen it MySQL connectors, when the collection
// does not exist, by default, no error is returned. To return ErrNotAvailable instead,
// use the functional option collection.GetValidateExistence.
func (s *Schema) GetCollection(ctx context.Context, name string, options ...collection.GetOption) (*Collection, error) {

	c, err := newCollection(s, name)
	if err != nil {
		return nil, fmt.Errorf("getting collection (%w)", err)
	}

	opts := collection.NewGetOptions(options...)
	if opts.ValidateExistence {
		if err := c.CheckExistence(ctx); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// GetCollections retrieve all available collections (does not include views or tables).
func (s *Schema) GetCollections(ctx context.Context) ([]*Collection, error) {

	names, err := s.objectNames(ctx, ObjectCollection)
	if err != nil {
		return nil, fmt.Errorf("getting collections (%w)", err)
	}

	if len(names) == 0 {
		return nil, nil
	}

	collections := make([]*Collection, len(names))
	for i, name := range names {
		collections[i], err = newCollection(s, name)
		if err != nil {
			return nil, fmt.Errorf("getting collections (%w)", err)
		}
	}

	return collections, nil
}

// CreateCollection creates a new collection. If the functional option
// collection.CreateReuseExisting is used, no error is reported when collection
// already exists.
func (s *Schema) CreateCollection(ctx context.Context, name string,
	options ...collection.CreateOption) (*Collection, error) {

	c, err := newCollection(s, name)
	if err != nil {
		return nil, fmt.Errorf("creating collection (%w)", err)
	}

	opts := collection.NewCreateOptions(options...)

	args := xproto.CommandArgs(
		xproto.ObjectField("schema", s.name),
		xproto.ObjectField("name", name),
		xproto.ObjectField("options", xproto.ObjectFields{
			xproto.ObjectField("reuse_existing", opts.ReuseExisting),
		}),
	)

	_, err = s.session.ExecCommand(ctx, "create_collection", args)
	if err != nil {
		return nil, fmt.Errorf("creating collection (%w)", err)
	}

	return c, nil
}

// DropCollection drops the collection.
func (s *Schema) DropCollection(ctx context.Context, name string) error {

	args := xproto.CommandArgs(
		xproto.ObjectField("schema", s.name),
		xproto.ObjectField("name", name),
	)

	_, err := s.session.ExecCommand(ctx, "drop_collection", args)
	if err != nil {
		return fmt.Errorf("dropping collection (%w)", err)
	}

	return nil
}

func (s *Schema) objectNames(ctx context.Context, kind ObjectKind) ([]string, error) {

	args := xproto.CommandArgs(
		xproto.ObjectField("schema", s.name),
	)

	if err := s.session.Write(ctx, xproto.Command("list_objects", args)); err != nil {
		return nil, err
	}

	res, err := s.session.handleResult(ctx, func(r *Result) bool {
		return r.stmtOK
	})
	if err != nil {
		return nil, err
	}

	if len(res.Rows) == 0 {
		return nil, nil
	}

	var names []string
	for _, row := range res.Rows {
		objType, ok := row.Values[1].(string)
		if !ok || objType != string(kind) {
			continue
		}

		name, ok := row.Values[0].(null.String)
		if !ok || !name.Valid {
			continue
		}

		names = append(names, name.String)
	}

	if len(names) == 0 {
		return nil, nil
	}

	sort.Strings(names)

	return names, nil
}
