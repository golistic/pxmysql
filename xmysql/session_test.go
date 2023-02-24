// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golistic/xt"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/null"
)

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestSession_ExecuteStatement(t *testing.T) {
	config := &ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: userNative,
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	t.Run("numeric data types", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "data_types_numeric"))

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.OK(t, ses.SetCurrentSchema(context.Background(), testSchema))

		exp := map[int64]struct {
			numBit               int64            // BIT
			numBool              bool             // BOOL
			numTinyIntSigned     int64            // TINYINT SIGNED
			numTinyIntUnsigned   uint64           // TINYINT UNSIGNED
			numSmallIntSigned    int64            // SMALLINT SIGNED
			numSmallIntUnsigned  uint64           // SMALLINT UNSIGNED
			numMediumIntSigned   int64            // MEDIUMINT SIGNED
			numMediumIntUnsigned uint64           // MEDIUMINT UNSIGNED
			numIntSigned         int64            // INT SIGNED
			numIntUnsigned       uint64           // INT UNSIGNED
			numBigIntSigned      int64            // BIGINT SIGNED
			numBigIntUnsigned    uint64           // BIGINT UNSIGNED
			numDecimal           *decimal.Decimal // DECIMAL(65,30)
			numDecimal2          *decimal.Decimal // DECIMAL(65,30)
			numDecimal3          *decimal.Decimal // DECIMAL(18,9)
		}{
			1: {
				numBit:               38,
				numBool:              false,
				numTinyIntSigned:     127,
				numTinyIntUnsigned:   0,
				numSmallIntSigned:    32767,
				numSmallIntUnsigned:  0,
				numMediumIntSigned:   8388607,
				numMediumIntUnsigned: 0,
				numIntSigned:         2147483647,
				numIntUnsigned:       0,
				numBigIntSigned:      int64(9223372036854775807),
				numBigIntUnsigned:    uint64(0),
				numDecimal:           decimal.MustNew("3.140000000000000000000000000000"),
				numDecimal2:          decimal.MustNew("9999999999999999999999999999999999999999999999999999999999991234.9"),
				numDecimal3:          decimal.MustNew("123456789.000001000"),
			},
			2: {
				numBit:               6,
				numBool:              true,
				numTinyIntSigned:     -128,
				numTinyIntUnsigned:   255,
				numSmallIntSigned:    -32768,
				numSmallIntUnsigned:  65535,
				numMediumIntSigned:   -8388608,
				numMediumIntUnsigned: 16777215,
				numIntSigned:         -2147483648,
				numIntUnsigned:       4294967295,
				numBigIntSigned:      int64(-9223372036854775808),
				numBigIntUnsigned:    uint64(18446744073709551615),
				numDecimal:           decimal.MustNew("-3.140000000000000000000000000000"),
				numDecimal2:          decimal.MustNew("-9999999999999999999999999999999999999999999999999999999999991234.5"),
				numDecimal3:          decimal.MustNew("-123456789.000001000"),
			},
		}

		res, err := ses.ExecuteStatement(context.Background(),
			"SELECT * FROM data_types_numeric ORDER BY id")
		xt.OK(t, err)

		for _, row := range res.Rows {
			id := row.Values[0].(int64)

			t.Run(fmt.Sprintf("row=%d", id), func(t *testing.T) {
				xt.Eq(t, exp[id].numBit, row.Values[1].(null.Uint64).Uint64)

				resBool := row.Values[2].(null.Int64)
				xt.Assert(t, resBool.Valid)
				xt.Eq(t, exp[id].numBool, resBool.Int64 == 1)

				resTinyIntSigned := row.Values[3].(null.Int64)
				xt.Eq(t, exp[id].numTinyIntSigned, resTinyIntSigned.Int64)

				resTinyIntUnsigned := row.Values[4].(null.Uint64)
				xt.Eq(t, exp[id].numTinyIntUnsigned, resTinyIntUnsigned.Uint64)

				xt.Eq(t, exp[id].numSmallIntSigned, row.Values[5].(int64))
				xt.Eq(t, exp[id].numSmallIntUnsigned, row.Values[6].(uint64))
				xt.Eq(t, exp[id].numMediumIntSigned, row.Values[7].(int64))
				xt.Eq(t, exp[id].numMediumIntUnsigned, row.Values[8].(uint64))
				xt.Eq(t, exp[id].numIntSigned, row.Values[9].(int64))
				xt.Eq(t, exp[id].numIntUnsigned, row.Values[10].(uint64))
				xt.Eq(t, exp[id].numBigIntSigned, row.Values[11].(int64))
				xt.Eq(t, exp[id].numBigIntUnsigned, row.Values[12].(uint64))
				xt.Assert(t, exp[id].numDecimal.Equal(row.Values[13].(*decimal.Decimal)),
					fmt.Sprintf("%s <> %s", exp[id].numDecimal, row.Values[13].(*decimal.Decimal)))
				xt.Assert(t, exp[id].numDecimal2.Equal(row.Values[14].(*decimal.Decimal)))
				xt.Assert(t, exp[id].numDecimal3.Equal(row.Values[15].(*decimal.Decimal)),
					fmt.Sprintf("expected %s == %s", exp[id].numDecimal3, row.Values[15].(*decimal.Decimal)))
			})
		}
	})

	t.Run("datetime data types", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "data_types_datetime"))

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.OK(t, ses.SetCurrentSchema(context.Background(), testSchema))

		locCET, err := time.LoadLocation("CET")
		xt.OK(t, err)

		exp := map[int64]struct {
			dtDATE       time.Time
			dtTIME       time.Duration
			dtDATETIME   time.Time
			dtTIMESTAMP  time.Time
			dtYear       int
			timeLocation *time.Location
		}{
			1: {
				dtDATE:       time.Date(2005, 3, 1, 0, 0, 0, 0, ses.timeLocation),
				dtTIME:       mustParseDuration("8h0m1.123456s"),
				dtDATETIME:   time.Date(2005, 3, 1, 7, 0, 1, 0, ses.timeLocation),
				dtTIMESTAMP:  time.Date(2005, 3, 1, 8, 0, 1, 0, locCET),
				dtYear:       2005,
				timeLocation: locCET,
			},
			2: {
				dtDATE:       time.Date(9999, 12, 31, 0, 0, 0, 0, ses.timeLocation),
				dtTIME:       mustParseDuration("838h59m59s"),
				dtDATETIME:   time.Date(9999, 12, 31, 23, 59, 59, 999999000, ses.timeLocation),
				dtTIMESTAMP:  time.Date(2038, 1, 19, 3, 14, 7, 0, ses.timeLocation),
				dtYear:       1901,
				timeLocation: time.UTC,
			},
			3: {
				dtDATE:       time.Date(1000, 1, 1, 0, 0, 0, 0, ses.timeLocation),
				dtTIME:       mustParseDuration("-838h59m59s"),
				dtDATETIME:   time.Date(1000, 1, 1, 0, 0, 0, 0, ses.timeLocation),
				dtTIMESTAMP:  time.Date(1970, 1, 1, 0, 0, 1, 0, ses.timeLocation),
				dtYear:       1901,
				timeLocation: time.UTC,
			},
		}

		res, err := ses.ExecuteStatement(context.Background(),
			"SELECT * FROM data_types_datetime ORDER BY id")
		xt.OK(t, err)

		for _, row := range res.Rows {
			id := row.Values[0].(int64)

			t.Run(fmt.Sprintf("row=%d", id), func(t *testing.T) {
				xt.Eq(t, exp[id].dtDATE, row.Values[1].(time.Time))
				xt.Eq(t, exp[id].dtTIME, row.Values[2].(time.Duration))
				xt.Eq(t, exp[id].dtDATETIME, row.Values[3].(time.Time))
				xt.Eq(t, exp[id].dtTIMESTAMP, row.Values[4].(time.Time).In(exp[id].timeLocation))
				xt.Eq(t, exp[id].dtYear, row.Values[5].(uint64))
			})
		}
	})

	t.Run("string data types", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "data_types_string"))

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.OK(t, ses.SetCurrentSchema(context.Background(), testSchema))

		res, err := ses.ExecuteStatement(context.Background(), "SELECT * FROM data_types_string ORDER BY id")
		xt.OK(t, err)

		exp := map[int64]struct {
			sChar      string
			sVarchar   string
			sBinary    []byte
			sVarBinary []byte
			sLongText  string
			sTinyBlob  []byte
			sEnum      string
			sSet       []string
		}{
			1: {
				sChar:      "CHAR" + strings.Repeat("a", 251),
				sVarchar:   "VARCHAR" + strings.Repeat("b", 393),
				sBinary:    []byte{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				sVarBinary: []byte{8, 9, 10, 11, 12, 13, 14, 15, 16},
				sLongText:  "LONGTEXT" + strings.Repeat("l", testMySQLMaxAllowedPacket-10),
				sTinyBlob:  []byte("I am a tiny blob"),
				sEnum:      "Go",
				sSet:       []string{"Go", "Python"},
			},
		}

		for _, row := range res.Rows {
			id := row.Values[0].(int64)

			t.Run(fmt.Sprintf("row=%d", id), func(t *testing.T) {
				xt.Assert(t, IsSupportedCollation(res.Columns[1].GetCollation()))
				xt.Eq(t, exp[id].sChar, row.Values[1].(string))
				xt.Eq(t, exp[id].sVarchar, row.Values[2].(string))
				xt.Eq(t, exp[id].sBinary, row.Values[3].([]byte))
				xt.Eq(t, exp[id].sVarBinary, row.Values[4].([]byte))
				xt.Eq(t, exp[id].sLongText, row.Values[5].(string))
				xt.Eq(t, exp[id].sTinyBlob, row.Values[6].([]byte))
				xt.Eq(t, exp[id].sEnum, row.Values[7].(string))
				xt.Eq(t, exp[id].sSet, row.Values[8].([]string))
			})
		}
	})

	t.Run("execute INSERT", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "inserting.sql"))

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.OK(t, ses.SetCurrentSchema(context.Background(), testSchema))

		res, err := ses.ExecuteStatement(context.Background(),
			"INSERT INTO inserts01 (c1) VALUES ('1'),('2')")
		xt.OK(t, err)

		xt.Eq(t, "Records: 2  Duplicates: 0  Warnings: 0", res.notices.producedMessage)
		xt.Eq(t, 2, res.RowsAffected())
		xt.Eq(t, 1, res.LastInsertID()) // first generated is returned with multiple values

		res, err = ses.ExecuteStatement(context.Background(),
			"INSERT INTO inserts01 (c1) VALUES ('3')")
		xt.OK(t, err)

		xt.Eq(t, 3, res.LastInsertID())
	})

	t.Run("use arguments and placeholders", func(t *testing.T) {
		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		res, err := ses.ExecuteStatement(context.Background(), "SELECT ?, ?, '?', ? from dual", 1, "one", 3)
		xt.OK(t, err)

		xt.Assert(t, len(res.Rows) == 1, "expected 1 row")
		xt.Eq(t, 1, res.Rows[0].Values[0].(int64))
		xt.Eq(t, "one", res.Rows[0].Values[1].(string))
		xt.Eq(t, "?", res.Rows[0].Values[2].(string))
		xt.Eq(t, 3, res.Rows[0].Values[3].(int64))
	})
}

func TestSession_CurrentSchema(t *testing.T) {
	t.Run("configure schema and get current schema name", func(t *testing.T) {
		config := &ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: userNative,
			Schema:   testSchema,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		schema, err := ses.CurrentSchema(context.Background())
		xt.OK(t, err)
		xt.Eq(t, testSchema, schema)
	})

	t.Run("no current schema in configuration means empty schema name", func(t *testing.T) {
		config := &ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		schema, err := ses.CurrentSchema(context.Background())
		xt.OK(t, err)
		xt.Eq(t, "", schema)
	})
}

func TestSession_SetTimeZone(t *testing.T) {
	config := &ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: userNative,
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	locName := "America/Los_Angeles"
	locUSALA, err := time.LoadLocation(locName)
	xt.OK(t, err)

	t.Run("default UTC when no time zone is set", func(t *testing.T) {
		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		// without time location set, timestamps returned by MySQL are UTC
		nowUSALA := time.Now().In(locUSALA)
		time.Sleep(1 * time.Second)
		res, err := ses.ExecuteStatement(context.Background(), "SELECT NOW(6)")
		xt.OK(t, err)

		myUTCNow := res.Rows[0].Values[0].(time.Time)
		xt.Eq(t, time.UTC, myUTCNow.Location())
		time.Sleep(time.Second)
		xt.Assert(t, myUTCNow.After(nowUSALA), fmt.Sprintf("expected %s > %s", myUTCNow, nowUSALA.UTC()))
	})

	t.Run("set time zone for session", func(t *testing.T) {
		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		// set time location for session
		xt.OK(t, ses.SetTimeZone(context.Background(), locName))
		myLoc, err := ses.TimeZone(context.Background())
		xt.OK(t, err)
		xt.Eq(t, locName, myLoc.String())

		res, err := ses.ExecuteStatement(context.Background(), "SELECT NOW(6)")
		xt.OK(t, err)
		myNow := res.Rows[0].Values[0].(time.Time)
		xt.Eq(t, locUSALA, myNow.Location())

		time.Sleep(time.Second)
		nowUTC := time.Now().UTC()
		xt.Assert(t, nowUTC.After(myNow), fmt.Sprintf("expected %s > %s", nowUTC, myNow.UTC()))
	})
}

func TestSession_SetCollation(t *testing.T) {
	config := &ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: userNative,
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	t.Run("set collation and retrieve", func(t *testing.T) {
		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		exp := collations["utf8mb4_sinhala_ci"]
		xt.OK(t, ses.SetCollation(context.Background(), "utf8mb4_sinhala_ci"))
		have, err := ses.Collation(context.Background())
		xt.OK(t, err)
		xt.Eq(t, exp.ID, have.ID)
		xt.Assert(t, &exp != have)
	})

	t.Run("set invalid collation", func(t *testing.T) {
		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		c := "big5_chinese_ci"
		err = ses.SetCollation(context.Background(), c)
		xt.KO(t, err)
		xt.Eq(t, fmt.Sprintf("failed setting collation ('%s' unsupported)", c), err.Error())
	})
}

func TestSession_PrepareStatement(t *testing.T) {
	config := &ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: userNative,
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	ses, err := cnx.NewSession(context.Background())
	xt.OK(t, err)

	stmt := "SELECT SQRT(POW(?,2) + POW(?,2)) AS hypotenuse"
	prep, err := ses.PrepareStatement(context.Background(), stmt)
	xt.OK(t, err)

	res, err := prep.Execute(context.Background(), 3, 4)
	xt.OK(t, err)

	xt.Assert(t, len(res.Rows) == 1, "expected 1 row")
	xt.Assert(t, res.Rows[0].Values[0].(null.Float64).Compare(float64(5)))
}

func TestSession_DeallocatePrepareStatement(t *testing.T) {
	config := &ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: userNative,
	}
	config.SetPassword(userNativePwd)

	cnx, err := NewConnection(config)
	xt.OK(t, err)

	ses, err := cnx.NewSession(context.Background())
	xt.OK(t, err)

	t.Run("successfully deallocate", func(t *testing.T) {
		stmt := "SELECT ?"
		prep, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		_, err = prep.Execute(context.Background(), 3)
		xt.OK(t, err)

		xt.OK(t, ses.deallocatePrepareStatement(context.Background(), prep.result.stmtID))

		_, err = prep.Execute(context.Background(), 3)
		xxt.AssertMySQLError(t, err, 5110)
	})

}
