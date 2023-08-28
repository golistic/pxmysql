// Copyright (c) 2022, 2023, Geert JM Vanderkelen

package xmysql

type CollationID int

type Collation struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	CharSet string `json:"charSet"`
}

// IsSupportedCollation returns whether c is a valid/supported collation. Note that
// MySQL Protocol X only supports the utf8mb4 character set, and consequently, only collations of utf8mb4.
// Argument c can be the internal MYSQL ID or MySQL name.
func IsSupportedCollation[T string | uint64 | int](c T) bool {
	var have bool
	switch v := any(c).(type) {
	case string:
		_, have = Collations[v]
	case uint64:
		_, have = collationIDs[v]
	case int:
		_, have = collationIDs[uint64(v)]
	}
	return have
}
