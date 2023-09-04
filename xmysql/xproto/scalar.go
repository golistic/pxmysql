// Copyright (c) 2023, Geert JM Vanderkelen

package xproto

import (
	"fmt"
	"reflect"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxdatatypes"
)

func Scalar[T reflect.Value | any](value T) *mysqlxdatatypes.Scalar {

	var rv reflect.Value

	switch v := any(value).(type) {
	case nil:
		return NilScalar()
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	switch {
	case rv.Kind() == reflect.Slice:
		switch {
		case rv.Type().Elem().Kind() == reflect.Uint8 && rv.Len() == 0: // empty []byte
			return NilScalar()
		case rv.Type().Elem().Kind() == reflect.Uint8: // []byte
			return BytesScalar(rv.Bytes())
		default:
			panic("unsupported scalar slice")
		}
	case rv.Kind() == reflect.Pointer && rv.IsNil():
		return NilScalar()
	}

	rv = reflect.Indirect(rv)

	switch rv.Kind() {
	case reflect.Bool:
		return BoolScalar(rv.Bool())
	case reflect.String:
		return StringScalar(rv.String())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return UnsignedIntScalar(rv.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return SignedIntScalar(rv.Int())
	case reflect.Float32:
		return Float32Scalar(float32(rv.Float()))
	case reflect.Float64:
		return Float64Scalar(rv.Float())
	default:
		panic(fmt.Sprintf("unsupported scalar value; was %s", rv.Kind()))
	}
}

func Bool(value bool) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: BoolScalar(value),
	}
}

func BoolScalar(value bool) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type:  mysqlxdatatypes.Scalar_V_BOOL.Enum(),
		VBool: proto.Bool(value),
	}
}

func Nil() *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: NilScalar(),
	}
}

func NilScalar() *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type: mysqlxdatatypes.Scalar_V_NULL.Enum(),
	}
}

func SignedInt[T constraints.Signed](value T) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: SignedIntScalar(value),
	}
}

func SignedIntScalar[T constraints.Signed](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type:       mysqlxdatatypes.Scalar_V_SINT.Enum(),
		VSignedInt: proto.Int64(int64(value)),
	}
}

func UnsignedInt[T constraints.Unsigned](value T) *mysqlxdatatypes.Any {
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: UnsignedIntScalar(value),
	}
}

func UnsignedIntScalar[T constraints.Unsigned](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type:         mysqlxdatatypes.Scalar_V_UINT.Enum(),
		VUnsignedInt: proto.Uint64(uint64(value)),
	}
}

func String(value any) *mysqlxdatatypes.Any {
	var v string
	if value != nil {
		switch sv := value.(type) {
		case string:
			v = sv
		case *string:
			if sv != nil {
				v = *sv
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
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: StringScalar(v),
	}
}

func StringScalar[T ~string](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type: mysqlxdatatypes.Scalar_V_STRING.Enum(),
		VString: &mysqlxdatatypes.Scalar_String{
			Value: []byte(value),
		},
	}
}

func Bytes[T ~[]byte](value T) *mysqlxdatatypes.Any {
	v := []byte(value)
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: BytesScalar(v),
	}
}

func BytesScalar[T ~[]byte](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type: mysqlxdatatypes.Scalar_V_OCTETS.Enum(),
		VOctets: &mysqlxdatatypes.Scalar_Octets{
			Value: value,
		},
	}
}

func Float32[T ~float32](value T) *mysqlxdatatypes.Any {
	v := float32(value)
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: Float32Scalar(v),
	}
}

func Float32Scalar[T ~float32](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type:   mysqlxdatatypes.Scalar_V_FLOAT.Enum(),
		VFloat: proto.Float32(float32(value)),
	}
}

func Float64[T ~float64](value T) *mysqlxdatatypes.Any {
	v := float64(value)
	return &mysqlxdatatypes.Any{
		Type:   mysqlxdatatypes.Any_SCALAR.Enum(),
		Scalar: Float64Scalar(v),
	}
}

func Float64Scalar[T ~float64](value T) *mysqlxdatatypes.Scalar {
	return &mysqlxdatatypes.Scalar{
		Type:    mysqlxdatatypes.Scalar_V_DOUBLE.Enum(),
		VDouble: proto.Float64(float64(value)),
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
