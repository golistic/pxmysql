// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
)

func Bool(value bool) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:  mysqlxdatatypes.Scalar_V_BOOL.Enum(),
			VBool: proto.Bool(value),
		},
	}
}

func Nil() *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
		},
	}
}

func SignedInt[T constraints.Signed](value T) *mysqlxdatatypes.Any {
	v := int64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:       mysqlxdatatypes.Scalar_V_SINT.Enum(),
			VSignedInt: proto.Int64(v),
		},
	}
}

func UnsignedInt[T constraints.Unsigned](value T) *mysqlxdatatypes.Any {
	v := uint64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:         mysqlxdatatypes.Scalar_V_UINT.Enum(),
			VUnsignedInt: proto.Uint64(v),
		},
	}
}

func String(value any) *mysqlxdatatypes.Any {
	var v []byte
	if value != nil {
		switch sv := value.(type) {
		case string:
			v = []byte(sv)
		case *string:
			if sv != nil {
				v = []byte(*sv)
			} else {
				return Nil()
			}
		default:
			panic(fmt.Sprintf("String accepts string or *string; not %T", value))
		}
	} else {
		return Nil()
	}

	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type: mysqlxdatatypes.Scalar_V_STRING.Enum(),
			VString: &mysqlxdatatypes.Scalar_String{
				Value: v,
			},
		},
	}
}

func Bytes[T []byte](value T) *mysqlxdatatypes.Any {
	v := []byte(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type: mysqlxdatatypes.Scalar_V_OCTETS.Enum(),
			VOctets: &mysqlxdatatypes.Scalar_Octets{
				Value: v,
			},
		},
	}
}

func Float32[T ~float32](value T) *mysqlxdatatypes.Any {
	v := float32(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:   mysqlxdatatypes.Scalar_V_FLOAT.Enum(),
			VFloat: proto.Float32(v),
		},
	}
}

func Float64[T ~float64](value T) *mysqlxdatatypes.Any {
	v := float64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:    mysqlxdatatypes.Scalar_V_DOUBLE.Enum(),
			VDouble: proto.Float64(v),
		},
	}
}

func Decimal(value decimal.Decimal) *mysqlxdatatypes.Any {
	// MySQL X Protocol does not support sending the encoded BCD.
	return String(value.String())
}

func Time(value time.Time, timeZoneName string) (*mysqlxdatatypes.Any, error) {
	tz, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return nil, err
	}
	return String(value.In(tz).Format("2006-01-02 15:04:05.999999")), nil
}
