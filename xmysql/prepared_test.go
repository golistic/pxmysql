// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xutil"
	"github.com/golistic/xstrings"
	"github.com/golistic/xt"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/mysqlerrors"
	"github.com/golistic/pxmysql/null"
)

func TestPrepared_Execute(t *testing.T) {
	xt.OK(t, testContext.Server.LoadSQLScript("prepared_stmt"))

	config := &ConnectConfig{
		Address:      testContext.XPluginAddr,
		Username:     userNative,
		Schema:       testSchema,
		TimeZoneName: "UTC",
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	ses, err := cnx.NewSession(context.Background())
	xt.OK(t, err)

	t.Run("all supported Go types", func(t *testing.T) {
		stmt := "SELECT ?"
		prep, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		float12Ptr := float32(12)
		float14Ptr := float64(14)
		boolPtr := true

		var cases = []any{
			nil,
			false,
			&boolPtr,
			uint(1),
			uint8(2),
			uint16(3),
			uint32(4),
			uint64(5),
			6,
			int8(7),
			int16(7),
			int32(7),
			int64(7),
			"8",
			xutil.StringPtr("9"),
			[]byte("10"),
			float32(11),
			&float12Ptr,
			float32(13),
			&float14Ptr,
			*decimal.MustNew("15.15"),
			decimal.MustNew("16.16"),
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("%T", c), func(t *testing.T) {
				_, err := prep.Execute(context.Background(), c)
				xt.OK(t, err)
			})
		}
	})

	numCols := []string{"bit_", "bool_", "tinyint_", "tinyint_unsigned", "smallint_", "smallint_unsigned",
		"mediumint_", "mediumint_unsigned", "int_", "int_unsigned",
		"bigint_", "bigint_unsigned",
		"decimal_",
		"float_", "float_unsigned", "double_", "double_unsigned",
	}

	t.Run("NOT NULL numeric data types", func(t *testing.T) {
		stmt := "INSERT INTO numeric_not_null (" + strings.Join(numCols, ",") +
			") VALUES(" + xstrings.RepeatJoin("?", len(numCols), ",") + ")"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT " + strings.Join(numCols, ",") + " FROM numeric_not_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vBits := uint64(38)
		vInt8 := int8(math.MinInt8)
		vUint8 := uint8(math.MaxInt8)
		vInt16 := int16(math.MinInt16)
		vUint16 := uint16(math.MaxUint16)
		vInt24 := int32(-1 << 23)
		vUint24 := uint32(1<<24 - 1)
		vInt := math.MinInt32
		vUint := math.MaxUint32
		vInt64 := int64(math.MinInt64)
		vUint64 := uint64(math.MaxUint64)
		vDecimal, err := decimal.New("6.180000000000000000000000000000")
		xt.OK(t, err)
		vFloat32Pos := float32(2838.233)
		vFloat32Neg := -1 * vFloat32Pos
		vFloat64Neg := -7382.34
		vFloat64Pos := 7382.34

		res, err := prepInsert.Execute(context.Background(),
			vBits, true, vInt8, vUint8, vInt16, vUint16,
			vInt24, vUint24, vInt, vUint, vInt64, vUint64,
			vDecimal, vFloat32Neg, vFloat32Pos, vFloat64Neg, vFloat64Pos,
		)
		xt.OK(t, err)
		xt.Eq(t, 1, res.RowsAffected())
		lastInsertID := res.LastInsertID()

		res, err = prepSelect.Execute(context.Background(), lastInsertID)
		xt.OK(t, err)

		xt.Eq(t, 1, len(res.Rows))
		row := res.Rows[0].Values

		// BIT
		xt.Eq(t, vBits, row[0].(uint64))
		// BOOL
		xt.Assert(t, row[1].(int64) > 0)
		// TINYINT
		xt.Eq(t, vInt8, int8(row[2].(int64)))
		xt.Eq(t, vUint8, uint8(row[3].(uint64)))
		// SMALLINT
		xt.Eq(t, vInt16, int16(row[4].(int64)))
		xt.Eq(t, vUint16, uint16(row[5].(uint64)))
		// MEDIUMINT
		xt.Eq(t, vInt24, int32(row[6].(int64)))
		xt.Eq(t, vUint24, uint32(row[7].(uint64)))
		// INT
		xt.Eq(t, vInt, int(row[8].(int64)))
		xt.Eq(t, vUint, uint(row[9].(uint64)))
		// BIGINT
		xt.Eq(t, vInt64, int(row[10].(int64)))
		xt.Eq(t, vUint64, uint(row[11].(uint64)))
		// DECIMAL
		xt.Assert(t, vDecimal.Equal(row[12].(*decimal.Decimal)))
		// FLOAT
		xt.Eq(t, vFloat32Neg, row[13].(float32))
		xt.Eq(t, vFloat32Pos, row[14].(float32))
		// DOUBLE
		xt.Eq(t, vFloat64Neg, row[15].(float64))
		xt.Eq(t, vFloat64Pos, row[16].(float64))
	})

	t.Run("NOT NULL numeric data types", func(t *testing.T) {
		stmt := "INSERT INTO numeric_null (" + strings.Join(numCols, ",") +
			") VALUES(" + xstrings.RepeatJoin("?", len(numCols), ",") + ")"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT " + strings.Join(numCols, ",") + " FROM numeric_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		t.Run("non-nil values", func(t *testing.T) {
			vBits := uint64(38)
			vInt8 := int8(math.MinInt8)
			vUint8 := uint8(math.MaxInt8)
			vInt16 := int16(math.MinInt16)
			vUint16 := uint16(math.MaxUint16)
			vInt24 := int32(-1 << 23)
			vUint24 := uint32(1<<24 - 1)
			vInt := math.MinInt32
			vUint := uint(math.MaxUint32)
			vInt64 := int64(math.MinInt64)
			vUint64 := uint64(math.MaxUint64)
			vDecimal, err := decimal.New("6.180000000000000000000000000000")
			xt.OK(t, err)
			vFloat32Pos := float32(2838.233)
			vFloat32Neg := -1 * vFloat32Pos
			vFloat64Neg := -7382.34
			vFloat64Pos := 7382.34

			res, err := prepInsert.Execute(context.Background(),
				vBits, true, vInt8, vUint8, vInt16, vUint16,
				vInt24, vUint24, vInt, vUint, vInt64, vUint64,
				vDecimal, vFloat32Neg, vFloat32Pos, vFloat64Neg, vFloat64Pos,
			)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// BIT
			xt.Assert(t, null.Compare(row[0].(null.Uint64), vBits))
			// BOOL
			xt.Assert(t, null.Compare(row[1].(null.Int64), 1),
				fmt.Sprintf("expected %d but got %v", 1, row[1].(null.Int64).Int64))
			// TINYINT
			xt.Assert(t, null.Compare(row[2].(null.Int64), vInt8))
			xt.Assert(t, null.Compare(row[3].(null.Uint64), vUint8))
			// SMALLINT
			xt.Assert(t, null.Compare(row[4].(null.Int64), vInt16))
			xt.Assert(t, null.Compare(row[5].(null.Uint64), vUint16))
			// MEDIUMINT
			xt.Assert(t, null.Compare(row[6].(null.Int64), vInt24))
			xt.Assert(t, null.Compare(row[7].(null.Uint64), vUint24))
			// INT
			xt.Assert(t, null.Compare(row[8].(null.Int64), vInt))
			xt.Assert(t, null.Compare(row[9].(null.Uint64), vUint))
			// BIGINT
			xt.Assert(t, null.Compare(row[10].(null.Int64), vInt64))
			xt.Assert(t, null.Compare(row[11].(null.Uint64), vUint64))
			// DECIMAL
			xt.Assert(t, vDecimal.Equal(row[12].(*decimal.Decimal)))
			// FLOAT
			xt.Assert(t, null.Compare(row[13].(null.Float32), vFloat32Neg))
			xt.Assert(t, null.Compare(row[14].(null.Float32), vFloat32Pos))
			// DOUBLE
			xt.Assert(t, null.Compare(row[15].(null.Float64), vFloat64Neg))
			xt.Assert(t, null.Compare(row[16].(null.Float64), vFloat64Pos))
		})
	})

	t.Run("select string and bytes", func(t *testing.T) {
		stmt := "SELECT ? AS char_, ? AS varchar_, ? AS binary_, ? AS longtext_"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vString := "this is a string"
		vBytes := []byte("these are bytes, duh")
		vStringLong := strings.Repeat("a", testMySQLMaxAllowedPacket-1000) // -1000 is compensation for whole packet

		_, err = prepSelect.Execute(context.Background(), vString, vString, vBytes, vStringLong)
		xt.OK(t, err)
	})

	t.Run("too big message", func(t *testing.T) {
		stmt := "SELECT ? AS longtext_"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vStringLong := strings.Repeat("a", testMySQLMaxAllowedPacket+1000)

		_, err = prepSelect.Execute(context.Background(), vStringLong)
		xt.KO(t, err)
		var errMySQL *mysqlerrors.Error
		xt.Assert(t, errors.As(err, &errMySQL), fmt.Sprintf("got: %s", err))
		xt.Eq(t, mysqlerrors.ClientNetPacketTooLarge, errMySQL.Code)

		_, err = prepSelect.Execute(context.Background(), "ok length")
		// connection was re-opened
		xt.KO(t, err)
		xt.Assert(t, errors.As(err, &errMySQL))
		xt.Eq(t, 5110, errMySQL.Code)
	})

	t.Run("big but ok length message", func(t *testing.T) {
		stmt := "SELECT ? AS char_, ? AS varchar_, ? AS binary_, ? AS longtext_"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vString := "this is a string"
		vBytes := []byte("these are bytes, duh")
		vStringLong := strings.Repeat("a", testMySQLMaxAllowedPacket-1000) // -1000 is compensation for whole packet

		_, err = prepSelect.Execute(context.Background(), vString, vString, vBytes, vStringLong)
		xt.OK(t, err)
	})

	t.Run("NOT NULL temporal data types", func(t *testing.T) {
		stmt := "INSERT INTO temporal_not_null (datetime_, date_, timestamp_, year_) VALUES (?,?,?,?)"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT datetime_, date_, timestamp_, year_ FROM temporal_not_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vTime := time.Now()
		vDate := time.Date(2022, 4, 30, 0, 0, 0, 0, time.UTC)
		vYear := 2022

		res, err := prepInsert.Execute(context.Background(),
			vTime, // DATETIME
			vDate, // DATE
			vTime, // TIMESTAMP
			vYear, // YEAR
		)
		xt.OK(t, err)
		xt.Eq(t, 1, res.RowsAffected())
		lastInsertID := res.LastInsertID()

		res, err = prepSelect.Execute(context.Background(), lastInsertID)
		xt.OK(t, err)

		vTimeLoc := vTime.In(ses.timeLocation)

		xt.Eq(t, 1, len(res.Rows))
		row := res.Rows[0].Values

		// DATETIME
		xt.Eq(t, vTimeLoc, row[0].(time.Time))

		// DATE
		xt.Eq(t, vDate, row[1].(time.Time))

		// TIMESTAMP
		xt.Eq(t, vTimeLoc, row[2].(time.Time))

		// YEAR
		xt.Eq(t, vYear, row[3].(uint64))
	})

	t.Run("NULL temporal data types", func(t *testing.T) {
		stmt := "INSERT INTO temporal_null (datetime_, date_, timestamp_, year_) VALUES (?,?,?,?)"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT datetime_, date_, timestamp_, year_ FROM temporal_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		t.Run("non-nil values", func(t *testing.T) {
			vTime := time.Now()
			vDate := time.Date(2022, 4, 30, 0, 0, 0, 0, time.UTC)
			vYear := uint64(2022)

			res, err := prepInsert.Execute(context.Background(),
				&vTime, // DATETIME
				&vDate, // DATE
				&vTime, // TIMESTAMP
				&vYear, // YEAR
			)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			vTimeLoc := vTime.In(ses.timeLocation)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// DATETIME
			xt.Assert(t, null.Compare(row[0].(null.Time), vTimeLoc))

			// DATE
			xt.Assert(t, null.Compare(row[1].(null.Time), vDate))

			// TIMESTAMP
			xt.Assert(t, null.Compare(row[2].(null.Time), vTimeLoc))

			// YEAR
			xt.Assert(t, null.Compare(row[3].(null.Uint64), vYear))
		})

		t.Run("nil values", func(t *testing.T) {
			res, err := prepInsert.Execute(context.Background(),
				nil, // DATETIME
				nil, // DATE
				nil, // TIMESTAMP
				nil, // YEAR
			)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// DATETIME
			xt.Assert(t, null.Compare(row[0].(null.Time), nil))

			// DATE
			xt.Assert(t, null.Compare(row[1].(null.Time), nil))

			// TIMESTAMP
			xt.Assert(t, null.Compare(row[2].(null.Time), nil))

			// YEAR
			xt.Assert(t, null.Compare(row[3].(null.Uint64), nil))
		})
	})

	strCols := []string{"char_", "binary_", "varchar_", "varbinary_",
		"tinyblob_", "tinytext_", "blob_", "text_", "mediumblob_", "mediumtext_", "longblob_", "longtext_",
		"enum_", "set_"}

	t.Run("NOT NULL string types", func(t *testing.T) {
		stmt := "INSERT INTO strings_not_null (" + strings.Join(strCols, ",") +
			") VALUES (" + xstrings.RepeatJoin("?", len(strCols), ",") + ")"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT " + strings.Join(strCols, ",") + " FROM strings_not_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		vChar := strings.Repeat("a", 200)
		vBinary := make([]byte, 255)
		copy(vBinary, vChar)
		vEnum := "Moon"
		vSet := []string{"Earth", "Mars"}

		res, err := prepInsert.Execute(context.Background(),
			vChar,   // CHAR
			vBinary, // BINARY
			vChar,   // VARCHAR
			vBinary, // VARBINARY
			vBinary, // TINYBLOB
			vChar,   // TINYTEXT
			vBinary, // BLOB
			vChar,   // TEXT
			vBinary, // MEDIUMBLOB
			vChar,   // MEDIUMTEXT
			vBinary, // LONGBLOB
			vChar,   // LONGTEXT
			vEnum,   // ENUM
			vSet,    // SET
		)
		xt.OK(t, err)
		xt.Eq(t, 1, res.RowsAffected())
		lastInsertID := res.LastInsertID()

		res, err = prepSelect.Execute(context.Background(), lastInsertID)
		xt.OK(t, err)

		xt.Eq(t, 1, len(res.Rows))
		row := res.Rows[0].Values

		// CHAR-types
		for colType, col := range map[string]int{
			"CHAR": 0, "VARCHAR": 2, "TINYTEXT": 5, "TEXT": 7, "MEDIUMTEXT": 9, "LONGTEXT": 11} {
			t.Run(colType, func(t *testing.T) {
				xt.Eq(t, vChar, row[col].(string))
			})
		}

		// BINARY-types
		for colType, col := range map[string]int{
			"BINARY": 1, "VARBINARY": 3, "TINYBLOB": 4, "BLOB": 6, "MEDIUMBLOB": 8, "LONGBLOB": 10} {
			t.Run(colType, func(t *testing.T) {
				xt.Eq(t, vBinary, row[col].([]byte))
			})
		}

		// ENUM
		xt.Eq(t, vEnum, row[12].(string))

		// SET
		xt.Eq(t, vSet, row[13].([]string))
	})

	t.Run("NULL string types", func(t *testing.T) {
		stmt := "INSERT INTO strings_null (" + strings.Join(strCols, ",") +
			") VALUES (" + xstrings.RepeatJoin("?", len(strCols), ",") + ")"
		prepInsert, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		stmt = "SELECT " + strings.Join(strCols, ",") + " FROM strings_null WHERE id = ?"
		prepSelect, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		t.Run("non-nil values", func(t *testing.T) {
			vChar := strings.Repeat("a", 200)
			vBinary := make([]byte, 255)
			copy(vBinary, vChar)
			vEnum := "Moon"
			vSet := []string{"Earth", "Mars"}

			res, err := prepInsert.Execute(context.Background(),
				vChar,   // CHAR
				vBinary, // BINARY
				vChar,   // VARCHAR
				vBinary, // VARBINARY
				vBinary, // TINYBLOB
				vChar,   // TINYTEXT
				vBinary, // BLOB
				vChar,   // TEXT
				vBinary, // MEDIUMBLOB
				vChar,   // MEDIUMTEXT
				vBinary, // LONGBLOB
				vChar,   // LONGTEXT
				vEnum,   // ENUM
				vSet,    // SET
			)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// CHAR-types
			for colType, col := range map[string]int{
				"CHAR": 0, "VARCHAR": 2, "TINYTEXT": 5, "TEXT": 7, "MEDIUMTEXT": 9, "LONGTEXT": 11} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.String), vChar))
				})
			}

			// BINARY-types
			for colType, col := range map[string]int{
				"BINARY": 1, "VARBINARY": 3, "TINYBLOB": 4, "BLOB": 6, "MEDIUMBLOB": 8, "LONGBLOB": 10} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.Bytes), vBinary))
				})
			}

			// ENUM
			t.Run("ENUM", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[12].(null.String), vEnum))
			})

			// SET
			t.Run("SET", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[13].(null.Strings), vSet))
			})
		})

		t.Run("nil values", func(t *testing.T) {
			// we explicitly set to nil, so we can reuse previously prepared statement
			nils := make([]any, 14)

			res, err := prepInsert.Execute(context.Background(), nils...)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// CHAR-types
			for colType, col := range map[string]int{
				"CHAR": 0, "VARCHAR": 2, "TINYTEXT": 5, "TEXT": 7, "MEDIUMTEXT": 9, "LONGTEXT": 11} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.String), nil))
				})
			}

			// BINARY-types
			for colType, col := range map[string]int{
				"BINARY": 1, "VARBINARY": 3, "TINYBLOB": 4, "BLOB": 6, "MEDIUMBLOB": 8, "LONGBLOB": 10} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.Bytes), nil))
				})
			}

			// ENUM
			t.Run("ENUM", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[12].(null.String), nil))
			})

			// SET
			t.Run("SET", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[13].(null.Strings), nil))
			})
		})

		t.Run("empty values", func(t *testing.T) {
			vChar := ""
			vBinary := make([]byte, 255)
			var vEnum *string
			var vSet []string

			res, err := prepInsert.Execute(context.Background(),
				vChar,   // CHAR
				vBinary, // BINARY
				vChar,   // VARCHAR
				vBinary, // VARBINARY
				vBinary, // TINYBLOB
				vChar,   // TINYTEXT
				vBinary, // BLOB
				vChar,   // TEXT
				vBinary, // MEDIUMBLOB
				vChar,   // MEDIUMTEXT
				vBinary, // LONGBLOB
				vChar,   // LONGTEXT
				vEnum,   // ENUM
				vSet,    // SET
			)
			xt.OK(t, err)
			xt.Eq(t, 1, res.RowsAffected())
			lastInsertID := res.LastInsertID()

			res, err = prepSelect.Execute(context.Background(), lastInsertID)
			xt.OK(t, err)

			xt.Eq(t, 1, len(res.Rows))
			row := res.Rows[0].Values

			// CHAR-types
			for colType, col := range map[string]int{
				"CHAR": 0, "VARCHAR": 2, "TINYTEXT": 5, "TEXT": 7, "MEDIUMTEXT": 9, "LONGTEXT": 11} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.String), vChar))
				})
			}

			// BINARY-types
			for colType, col := range map[string]int{
				"BINARY": 1, "VARBINARY": 3, "TINYBLOB": 4, "BLOB": 6, "MEDIUMBLOB": 8, "LONGBLOB": 10} {
				t.Run(colType, func(t *testing.T) {
					xt.Assert(t, null.Compare(row[col].(null.Bytes), vBinary))
				})
			}

			// ENUM
			t.Run("ENUM", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[12].(null.String), nil))
			})

			// SET
			t.Run("SET", func(t *testing.T) {
				xt.Assert(t, null.Compare(row[13].(null.Strings), vSet))
			})
		})
	})
}
