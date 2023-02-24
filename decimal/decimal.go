// Copyright (c) 2022, Geert JM Vanderkelen

package decimal

import (
	"fmt"
	"math/big"
	"strings"
)

// Decimal represents a fixed-point decimal with an integer and fraction part.
// MySQL allows a DECIMAL-type to have 1 to 65 digits, with the scale or fraction
// having 0 to 30 digits. Both the integral and fraction part are therefor stored
// as big.Int.
// This Decimal-type is just storage. The caller can access the integral and
// fraction part through the methods Integral() and Fraction(). There are no
// methods to do calculations.
type Decimal struct {
	integral *big.Int
	fraction *big.Int
	maxScale int
}

// Zero is 0 as a Decimal.
var Zero = MustNew("0")

// New takes string s as decimal (for example "3.14"), and stores the
// integer and fraction in a newly instantiated Decimal object.
// Only base 10 is supported, so no `0x` prefixes for example. The call is responsible
// to remove any thousand separators.
func New(s string) (*Decimal, error) {
	var errBad = "invalid decimal string (%w)"

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf(errBad, fmt.Errorf("too many separators"))
	}

	d := &Decimal{
		integral: &big.Int{},
		fraction: &big.Int{},
	}

	_, ok := d.integral.SetString(parts[0], 10)
	if !ok {
		return nil, fmt.Errorf(errBad, fmt.Errorf("bad integral part"))
	}

	if len(parts) == 2 {
		d.maxScale = len(parts[1])
		if d.maxScale == 0 {
			return nil, fmt.Errorf(errBad, fmt.Errorf("bad fractional part"))
		}
		if parts[1] != strings.Repeat("0", d.maxScale) {
			_, ok := d.fraction.SetString(strings.TrimLeft(parts[1], "0"), 10)
			if !ok {
				return nil, fmt.Errorf(errBad, fmt.Errorf("bad fractional part"))
			}
		}
	}

	return d, nil
}

// MustNew is similar to New but instead of return the error, it panics.
func MustNew(s string) *Decimal {
	d, err := New(s)
	if err != nil {
		panic(err)
	}
	return d
}

func (d *Decimal) fractionAsString() string {
	f := d.fraction.String()
	diff := d.maxScale - len(f)
	if diff > 0 {
		f = strings.Repeat("0", diff) + f
	}
	return f
}

// String returns the textual representation of Decimal as decimal number
// with a dot '.' as separator.
func (d *Decimal) String() string {
	s := d.integral.String()
	if d.fraction != nil && d.fraction.BitLen() > 0 {
		s += "." + d.fractionAsString()
	}
	return s
}

// Equal returns whether d is equal to the other.
func (d *Decimal) Equal(other *Decimal) bool {
	if d == nil || other == nil {
		return false
	}

	return d.integral.Cmp(other.integral) == 0 && d.fraction.Cmp(other.fraction) == 0
}

// Integral returns the integer part of d.
func (d *Decimal) Integral() *big.Int {
	return d.integral
}

// Sign returns the sign information of the integer part of d:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
func (d *Decimal) Sign() int {
	return d.integral.Sign()
}

// Fraction returns the fraction part of d.
func (d *Decimal) Fraction() *big.Int {
	return d.fraction
}

// Encode will encode d as Binary-Coded Decimal (BCD), making it ready
// to send it to MySQL.
func (d *Decimal) Encode() ([]byte, error) {
	bs := []byte(d.integral.String() + d.fractionAsString())

	negative := d.Sign() < 0
	if negative {
		bs = bs[1:]
	}

	bcd := make([]byte, 0, len(bs)/2+2)
	bcd = append(bcd, byte(d.maxScale))

	var signed bool
	for i := 0; i < len(bs); i += 2 {
		b := (bs[i] - 48) << 4
		if i+1 < len(bs) {
			b |= (bs[i+1] - 48) & 0x0f
		} else {
			switch {
			case negative:
				b |= 0x0d
			default:
				b |= 0x0c
			}
			signed = true
		}
		bcd = append(bcd, b)
	}

	if !signed {
		switch {
		case negative:
			bcd = append(bcd, 0xd0)
		default:
			bcd = append(bcd, 0xc0)
		}
	}

	return bcd, nil
}

// NewDecimalFromBCD decodes bcd or Binary-Coded Decimal and returns
// a new instance of Decimal.
func NewDecimalFromBCD(bcd []byte) (*Decimal, error) {
	const digits = "0123456789"

	var errDecode = "cannot decode binary-coded decimal (%w)"

	if len(bcd) < 2 {
		return nil, fmt.Errorf(errDecode, fmt.Errorf("not enough data"))
	}

	d := &Decimal{
		integral: &big.Int{},
		fraction: &big.Int{},
		maxScale: int(bcd[0]),
	}

	last := bcd[len(bcd)-1]

	var s string
	for _, b := range bcd[1 : len(bcd)-1] {
		hi := (b >> 4) & 0x0f
		lo := b & 0x0f
		if hi < 0x0a {
			s += string(digits[hi]) + string(digits[lo])
		} else {
			return nil, fmt.Errorf(errDecode, fmt.Errorf("found bits > 9"))
		}
	}

	if last != 0xc0 && last != 0xd0 {
		hi := (last >> 4) & 0x0f
		last <<= 4
		if hi > 9 {
			return nil, fmt.Errorf(errDecode, fmt.Errorf("last byte value was > 9"))
		} else if last != 0xc0 && last != 0xd0 {
			return nil, fmt.Errorf(errDecode, fmt.Errorf("bad signing bit"))
		}
		s += string(digits[hi])
	}

	var ok bool
	if d.maxScale != 0 {
		cut := len(s) - d.maxScale

		if cut < 0 {
			return nil, fmt.Errorf(errDecode, fmt.Errorf("not enough data with fraction"))
		}

		// ignoring ok; cannot find way to fail
		d.integral, _ = (&big.Int{}).SetString(s[:cut], 10)
		d.fraction, _ = (&big.Int{}).SetString(s[cut:], 10)
	} else {
		d.integral, ok = (&big.Int{}).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf(errDecode, fmt.Errorf("bad integral part"))
		}
	}

	if last == 0xd0 { // 0xc0 would mean negative
		d.integral.Neg(d.integral)
	}

	return d, nil
}
