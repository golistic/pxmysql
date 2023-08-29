// Copyright (c) 2023, Geert JM Vanderkelen

package statements

import (
	"fmt"
	"strings"
)

// QuoteValue quotes p so that it can be safely used to substituted placeholders
// within a SQL query.
func QuoteValue(p any) (string, error) {

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

func QuoteIdentifier(p string) (string, error) {
	return "`" + strings.Replace(p, "`", "``", -1) + "`", nil
}
