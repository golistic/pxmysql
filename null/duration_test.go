// Copyright (c) 2022, Geert JM Vanderkelen

package null

import (
	"testing"
	"time"

	"github.com/golistic/xt"
)

func TestDuration_Compare(t *testing.T) {
	dur2d3h, _ := time.ParseDuration("2d3h")
	dur5h6m, _ := time.ParseDuration("5h6m")

	var cases = []struct {
		n     Nullable
		value any
		exp   bool
	}{
		{
			n:     Duration{Duration: dur2d3h, Valid: true},
			value: &dur2d3h,
			exp:   true,
		},
		{
			n:     Duration{Duration: dur2d3h, Valid: true},
			value: &dur5h6m, // supposed to be different
			exp:   false,
		},
		{
			n:     Duration{Duration: dur2d3h, Valid: true},
			value: dur2d3h,
			exp:   true,
		},
		{
			n:     Duration{Duration: dur2d3h, Valid: true},
			value: dur5h6m, // supposed to be different
			exp:   false,
		},
		{
			n:     Duration{Valid: false},
			value: nil,
			exp:   true,
		},

		{
			n:     Duration{Valid: false},
			value: 0,
			exp:   false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			xt.Eq(t, c.exp, c.n.Compare(c.value))
		})
	}

	t.Run("panics if value type is not supported", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = Duration{Duration: 0, Valid: true}.Compare("str")
		})
	})
}

func TestDuration_Value(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dur2d3h, _ := time.ParseDuration("2d3h")
		nd := Duration{Duration: dur2d3h, Valid: true}
		v, _ := nd.Value()
		d, ok := v.(time.Duration)
		xt.Assert(t, ok, "expected time.Duration")
		xt.Eq(t, dur2d3h, d)
	})

	t.Run("not valid", func(t *testing.T) {
		nd := Duration{Duration: 0, Valid: false}
		v, _ := nd.Value()
		xt.Eq(t, nil, v, "expected nil")
	})
}
