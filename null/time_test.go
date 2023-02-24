// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"
	"time"

	"github.com/golistic/xt"
	"github.com/golistic/xtime"
)

func TestTime_Compare(t *testing.T) {
	now := time.Now()
	yesterday := xtime.Yesterday()

	t.Run("Uint64", func(t *testing.T) {
		var cases = []struct {
			n     Nullable
			value *time.Time
			exp   bool
		}{
			{
				n:     Time{Time: now, Valid: true},
				value: &now,
				exp:   true,
			},
			{
				n:     Time{Time: now, Valid: true},
				value: &yesterday,
				exp:   false,
			},
			{
				n:     Time{Time: now, Valid: false},
				value: nil,
				exp:   false,
			},
		}

		for _, c := range cases {
			t.Run("", func(t *testing.T) {
				xt.Eq(t, c.exp, c.n.Compare(c.value))
			})
		}
	})

	t.Run("non-pointer", func(t *testing.T) {
		xt.Assert(t, !Time{Time: now, Valid: false}.Compare(now))
		xt.Assert(t, Time{Time: now, Valid: true}.Compare(now))
	})

	t.Run("panics if value type is not supported", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = Time{Time: now, Valid: true}.Compare("str")
		})
	})

	t.Run("value is explicitly nil", func(t *testing.T) {
		xt.Assert(t, Time{Time: now, Valid: false}.Compare(nil))
	})
}

func TestTime_Value(t *testing.T) {
	now := time.Now()

	t.Run("valid", func(t *testing.T) {
		nt := Time{Time: now, Valid: true}
		v, _ := nt.Value()
		d, ok := v.(time.Time)
		xt.Assert(t, ok, "expected time.Time")
		xt.Eq(t, now, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nd := Time{Time: now, Valid: false}
		v, _ := nd.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
