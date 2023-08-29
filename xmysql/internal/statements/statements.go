// Copyright (c) 2023, Geert JM Vanderkelen

package statements

import (
	"bytes"
	"fmt"

	"github.com/golistic/xgo/xmath"
)

const Placeholder = '?'

// SubstitutePlaceholders replaces the placeholders within stmt with respective element of args.
func SubstitutePlaceholders(stmt string, args ...any) (string, error) {

	placeholders := PlaceholderIndexes(Placeholder, stmt)
	if len(placeholders) != len(args) {
		return "", fmt.Errorf("need %d placeholder(s); found %d)", len(args), len(placeholders))
	}

	var nextArg int
	var buf []byte

	var index int
	for _, ph := range placeholders {
		buf = append(buf, stmt[index:ph]...)

		arg := args[nextArg]
		nextArg++
		index = ph + 1

		if arg == nil {
			buf = append(buf, "NULL"...)
			continue
		}

		quoted, err := QuoteValue(arg)
		if err != nil {
			return "", err
		}
		buf = append(buf, quoted...)
	}

	// rest of stmt
	buf = append(buf, stmt[index:]...)

	if len(args) > nextArg {
		return "", fmt.Errorf("%d argument(s) not substituted", xmath.AbsInt(len(args)-nextArg))
	} else if len(args) < nextArg {
		return "", fmt.Errorf("%d placeholder(s) not substituted", xmath.AbsInt(len(args)-nextArg))
	}

	return string(buf), nil
}

// PlaceholderIndexes returns the indices of all placeholders within query.
func PlaceholderIndexes(placeholder rune, query string) []int {

	var indexes []int

	var quoted bool
	var quote rune
	for i, r := range bytes.Runes([]byte(query)) {
		// we skip quoted so that we support queries which have placeholder in string literals
		if r == '"' || r == '\'' {
			if quoted && quote == r {
				quoted = false
				quote = 0
				continue
			} else if !quoted {
				quoted = true
				quote = r
				continue
			}
		}

		if quoted {
			continue
		}

		if r == placeholder {
			indexes = append(indexes, i)
		}
	}

	return indexes
}
