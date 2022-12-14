// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"bytes"
	"fmt"

	"github.com/geertjanvdk/xkit/xmath"
)

const stmtPlaceholder = '?'

// substitutePlaceholders replaces the placeholders within stmt with respective element of args.
func substitutePlaceholders(stmt string, args ...any) (string, error) {
	placeholders := placeholderIndexes(stmtPlaceholder, stmt)
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

		quoted, err := quoteSQLValue(arg)
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

func placeholderIndexes(placeholder rune, query string) []int {
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

func quoteSQLValue(p any) (string, error) {
	switch v := p.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case []byte:
		return fmt.Sprintf("_binary'%x'", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case string:
		// handled as default (last return)
	default:
		return "", fmt.Errorf("cannot quote parameter with value type %T", p)
	}

	return "'" + p.(string) + "'", nil
}
