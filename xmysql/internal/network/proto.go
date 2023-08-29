// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

// UnmarshalPartial parses the wire-format message in b and places the result in m.
// The provided message must be mutable (e.g., a non-nil pointer to a message).
// This is the same function as proto.Unmarshall except that AllowPartial option set to true.
func UnmarshalPartial(b []byte, m proto.Message) error {
	return proto.UnmarshalOptions{
		RecursionLimit: protowire.DefaultRecursionLimit,
		AllowPartial:   true,
	}.Unmarshal(b, m)
}
