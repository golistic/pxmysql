// Copyright (c) 2022, Geert JM Vanderkelen

package decimal

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/golistic/xgo/xt"
)

func TestNewDecimal(t *testing.T) {
	t.Run("valid decimals", func(t *testing.T) {
		var cases = map[string]string{
			"123456789.000001": "", // empty means output is equal to input
			"123.000019":       "",
			"-453.00200":       "",
			"1":                "",
			"0001.0001":        "1.0001",
		}

		for c, exp := range cases {
			if exp == "" {
				exp = c
			}
			t.Run(c, func(t *testing.T) {
				d, err := New(c)
				xt.OK(t, err)
				xt.Eq(t, exp, d.String())
				if strings.Contains(c, ".") {
					xt.Eq(t, len(strings.Split(c, ".")[1]), d.maxScale)
				}
			})
		}
	})

	t.Run("invalid integral part", func(t *testing.T) {
		var cases = []string{
			".007",
			"A.007",
			"0x0a.007",
		}

		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				_, err := New(c)
				xt.KO(t, err)
				xt.Eq(t, "bad integral part", errors.Unwrap(err).Error())
			})
		}
	})

	t.Run("invalid fractional part", func(t *testing.T) {
		var cases = []string{
			"0.0x07",
			"0.00L",
			"1.",
		}

		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				_, err := New(c)
				xt.KO(t, err)
				xt.Eq(t, "bad fractional part", errors.Unwrap(err).Error())
			})
		}
	})

	t.Run("bad fractional part (too many separators)", func(t *testing.T) {
		_, err := New("1.2.3")
		xt.KO(t, err)
		xt.Eq(t, "too many separators", errors.Unwrap(err).Error())
	})
}

func TestDecimal_Encode(t *testing.T) {
	var cases = []struct {
		d   *Decimal
		exp []byte
	}{
		{
			d:   MustNew("-123456789.000001000"),
			exp: []byte{9, 18, 52, 86, 120, 144, 0, 0, 16, 0, 0xd0},
		},
		{
			d:   MustNew("123456789.000001000"),
			exp: []byte{0x09, 18, 52, 86, 120, 0x90, 0x00, 0x00, 0x10, 0x00, 0xc0},
		},
		{
			d:   MustNew("123456789.0100"),
			exp: []byte{0x04, 18, 52, 86, 120, 0x90, 0x10, 0x0c},
		},
		{
			d:   MustNew("-123456789.0100"),
			exp: []byte{0x04, 18, 52, 86, 120, 0x90, 0x10, 0x0d},
		},
		{
			d:   MustNew("3.140000000000000000000000000000"),
			exp: []byte{30, 49, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12},
		},
		{
			d: MustNew("9999999999999999999999999999999999999999999999999999999999991234.9"),
			exp: []byte{1, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153,
				153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 18, 52, 156},
		},
	}

	for _, c := range cases {
		t.Run(c.d.String(), func(t *testing.T) {
			dec, err := c.d.Encode()
			xt.OK(t, err)
			xt.Eq(t, c.exp, dec)
		})
	}
}

func TestNewDecimalFromBCD(t *testing.T) {
	t.Run("good cases", func(t *testing.T) {
		var cases = []struct {
			b   []byte
			exp *Decimal
		}{
			{
				exp: MustNew("0"),
				b:   []byte{0x00, 0x0c},
			},
			{
				exp: MustNew("1"),
				b:   []byte{0x00, 0x1c},
			},
			{
				exp: MustNew("-1.00000"),
				b:   []byte{0x05, 0x10, 0x00, 0x00, 0xd0},
			},
			{
				exp: MustNew("-123456789.000001000"),
				b:   []byte{9, 18, 52, 86, 120, 144, 0, 0, 16, 0, 208},
			},
			{
				exp: MustNew("123456789.000001000"),
				b:   []byte{9, 18, 52, 86, 120, 144, 0, 0, 16, 0, 192},
			},
			{
				exp: MustNew("3.140000000000000000000000000000"),
				b:   []byte{30, 49, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12},
			},
			{
				exp: MustNew("9999999999999999999999999999999999999999999999999999999999991234.9"),
				b: []byte{1, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153,
					153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 18, 52, 156},
			},
		}

		for _, c := range cases {
			t.Run(c.exp.String(), func(t *testing.T) {
				dec, err := NewDecimalFromBCD(c.b)
				xt.OK(t, err)
				xt.Assert(t, c.exp.Equal(dec))
			})
		}
	})

	t.Run("bad cases", func(t *testing.T) {
		var cases = []struct {
			b   []byte
			exp error
		}{
			{
				exp: fmt.Errorf("not enough data"),
				b:   []byte{0x09},
			},
			{
				exp: fmt.Errorf("not enough data with fraction"),
				b:   []byte{0x09, 0x0d},
			},
			{
				exp: fmt.Errorf("bad signing bit"),
				// last byte incorrect
				b: []byte{0x09, 0x12, 0x34, 0x56, 0x78, 0x90, 0x00, 0x00, 0x10, 0x00, 0x8e},
			},
			{
				exp: fmt.Errorf("last byte value was > 9"),
				// last byte incorrect
				b: []byte{0x09, 0x12, 0x34, 0x56, 0x78, 0x90, 0x00, 0x00, 0x10, 0x00, 0xad},
			},
			{
				exp: fmt.Errorf("found bits > 9"),
				// 0xff is not valid
				b: []byte{0x04, 0x12, 0x34, 0xff, 0xff, 0xd0},
			},
		}

		for _, c := range cases {
			t.Run(c.exp.Error(), func(t *testing.T) {
				_, err := NewDecimalFromBCD(c.b)
				xt.KO(t, err)
				xt.Eq(t, c.exp, errors.Unwrap(err))
			})
		}
	})
}

func TestZero(t *testing.T) {
	xt.Assert(t, Zero.Equal(MustNew("0")))
}

func TestDecimal_Fraction(t *testing.T) {
	t.Run("> 0", func(t *testing.T) {
		dec := MustNew("1234.98774")
		xt.Eq(t, int64(98774), dec.Fraction().Int64())
	})

	t.Run("beyond 64-bit", func(t *testing.T) {
		expStr := "99999999999999999999999"
		exp, ok := (&big.Int{}).SetString(expStr, 10)
		xt.Assert(t, ok)
		dec := MustNew("1234." + expStr)
		xt.Assert(t, exp.Cmp(dec.Fraction()) == 0)
	})

	t.Run("== 0", func(t *testing.T) {
		dec := MustNew("1234.0")
		xt.Eq(t, int64(0), dec.Fraction().Int64())
	})
}

func TestDecimal_Integral(t *testing.T) {
	t.Run("> 0", func(t *testing.T) {
		dec := MustNew("1234.98774")
		xt.Eq(t, int64(1234), dec.Integral().Int64())
	})

	t.Run("beyond 64-bit", func(t *testing.T) {
		expStr := "99999999999999999999999"
		exp, ok := (&big.Int{}).SetString(expStr, 10)
		xt.Assert(t, ok)
		dec := MustNew(expStr + ".98774")
		xt.Assert(t, exp.Cmp(dec.Integral()) == 0)
	})

	t.Run("== 0", func(t *testing.T) {
		dec := MustNew("0.98774")
		xt.Eq(t, int64(0), dec.Integral().Int64())
	})
}

func TestDecimal_Equal(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		var left *Decimal
		var right *Decimal

		xt.Assert(t, !left.Equal(right))
	})

	t.Run("integral part not equal", func(t *testing.T) {
		left := MustNew("1234.9876")
		right := MustNew("5566.9876")

		xt.Assert(t, !left.Equal(right))
	})

	t.Run("fractional part not equal", func(t *testing.T) {
		left := MustNew("1234.9876")
		right := MustNew("1234.5566")

		xt.Assert(t, !left.Equal(right))
	})

	t.Run("not equal", func(t *testing.T) {
		left := MustNew("5566.6655")
		right := MustNew("1234.1234")

		xt.Assert(t, !left.Equal(right))
	})

	t.Run("equal", func(t *testing.T) {
		left := MustNew("1234.9876")
		right := MustNew("1234.9876")

		xt.Assert(t, left.Equal(right))
	})
}

func TestMustNew(t *testing.T) {
	xt.Panics(t, func() {
		MustNew("abc")
	})
}
