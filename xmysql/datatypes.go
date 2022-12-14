// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
)

func boolAsScalar(value bool) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:  mysqlxdatatypes.Scalar_V_BOOL.Enum(),
			VBool: proto.Bool(value),
		},
	}
}

func nilAsScalar() *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
		},
	}
}

func signedIntAsScalar[T constraints.Signed](value T) *mysqlxdatatypes.Any {
	v := int64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:       mysqlxdatatypes.Scalar_V_SINT.Enum(),
			VSignedInt: proto.Int64(v),
		},
	}
}

func unsignedIntAsScalar[T constraints.Unsigned](value T) *mysqlxdatatypes.Any {
	v := uint64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:         mysqlxdatatypes.Scalar_V_UINT.Enum(),
			VUnsignedInt: proto.Uint64(v),
		},
	}
}

func stringAsScalar(value any) *mysqlxdatatypes.Any {
	var v []byte
	if value != nil {
		switch sv := value.(type) {
		case string:
			v = []byte(sv)
		case *string:
			if sv != nil {
				v = []byte(*sv)
			} else {
				return nilAsScalar()
			}
		default:
			panic(fmt.Sprintf("stringAsScalar accepts string or *string; not %T", value))
		}
	} else {
		return nilAsScalar()
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

func byteSliceAsScalar[T []byte](value T) *mysqlxdatatypes.Any {
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

func float32IntAsScalar[T ~float32](value T) *mysqlxdatatypes.Any {
	v := float32(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:   mysqlxdatatypes.Scalar_V_FLOAT.Enum(),
			VFloat: proto.Float32(v),
		},
	}
}

func float64IntAsScalar[T ~float64](value T) *mysqlxdatatypes.Any {
	v := float64(value)
	return &mysqlxdatatypes.Any{
		Type: mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: &mysqlxdatatypes.Scalar{
			Type:    mysqlxdatatypes.Scalar_V_DOUBLE.Enum(),
			VDouble: proto.Float64(v),
		},
	}
}

func decimalAsScalar(value decimal.Decimal) *mysqlxdatatypes.Any {
	// MySQL X Protocol does not support sending the encoded BCD.
	return stringAsScalar(value.String())
}

func timeAsScalar(value time.Time, timeZoneName string) (*mysqlxdatatypes.Any, error) {
	tz, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return nil, err
	}
	return stringAsScalar(value.In(tz).Format("2006-01-02 15:04:05.999999")), nil
}
