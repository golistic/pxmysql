// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golistic/xgo/xnet"
	"github.com/golistic/xgo/xstrings"
	"github.com/golistic/xgo/xt"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/xxt"
	"github.com/golistic/pxmysql/null"
	"github.com/golistic/pxmysql/xmysql"
)

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestNewSession(t *testing.T) {
	t.Run("password is not stored in stored config", func(t *testing.T) {
		expConfig := &xmysql.ConnectConfig{
			Address:  "127.0.0.1",
			Username: "toor",
			Password: xstrings.Pointer("secret"),
			UseTLS:   false,
		}
		expConfig.SetPassword("toor")

		ses, err := xmysql.NewSession(expConfig)
		xt.OK(t, err)

		xt.Assert(t, ses.Config().Password == nil)
		xt.Eq(t, ses.Config().AuthMethod, xmysql.DefaultConnectConfig.AuthMethod)
	})

	t.Run("if not given, AuthMethod is default", func(t *testing.T) {
		expConfig := &xmysql.ConnectConfig{
			Address:  "127.0.0.1",
			Username: "toor",
		}
		expConfig.SetPassword("toor")

		ses, err := xmysql.NewSession(expConfig)
		xt.OK(t, err)

		xt.Eq(t, ses.Config().AuthMethod, xmysql.DefaultConnectConfig.AuthMethod)
	})

	t.Run("default configuration", func(t *testing.T) {
		ses, err := xmysql.NewSession(nil)
		xt.OK(t, err)

		xt.Assert(t, ses.Config().Password == nil)
		xt.Eq(t, ses.Config().Address, xmysql.DefaultConnectConfig.Address)
		xt.Eq(t, ses.Config().Username, xmysql.DefaultConnectConfig.Username)
		xt.Eq(t, ses.Config().AuthMethod, xmysql.DefaultConnectConfig.AuthMethod)
	})

	t.Run("MySQL X Plugin is reachable", func(t *testing.T) {
		expConfig := &xmysql.ConnectConfig{
			Address: testContext.XPluginAddr,
		}

		ses, err := xmysql.NewSession(expConfig)
		xt.OK(t, err)

		xt.Assert(t, ses.IsReachable())
	})

	t.Run("MySQL X Plugin is not reachable", func(t *testing.T) {
		expConfig := &xmysql.ConnectConfig{
			Address: "127.0.0.40",
		}

		ses, err := xmysql.NewSession(expConfig)
		xt.OK(t, err)

		xt.Assert(t, !ses.IsReachable())
	})

	t.Run("address getting defaults if needed", func(t *testing.T) {
		var cases = []struct {
			have string
			exp  string
		}{
			{have: "127.0.0.1", exp: "127.0.0.1:" + xmysql.DefaultPort},
			{have: ":" + xmysql.DefaultPort, exp: xmysql.DefaultHost + ":" + "33060"},
			{have: ":12453", exp: xmysql.DefaultHost + ":" + "12453"},
			{have: "0.0.0.0:" + xmysql.DefaultPort, exp: "0.0.0.0:" + xmysql.DefaultPort},
			{have: "", exp: xmysql.DefaultHost + ":" + xmysql.DefaultPort},
		}

		for _, c := range cases {
			t.Run(c.exp, func(t *testing.T) {
				ses, err := xmysql.NewSession(&xmysql.ConnectConfig{
					Address: c.have,
				})
				xt.OK(t, err)
				xt.Eq(t, c.exp, ses.Config().Address)
			})
		}
	})

	t.Run("valid time zone", func(t *testing.T) {
		locName := "Europe/Berlin"
		exp, err := time.LoadLocation(locName)
		xt.OK(t, err)

		config := &xmysql.ConnectConfig{
			TimeZoneName: locName,
		}

		ses, err := xmysql.NewSession(config)
		xt.OK(t, err)

		xt.Eq(t, exp, ses.TimeLocation())
	})

	t.Run("invalid time zone", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			TimeZoneName: "Foo/Bar",
		}

		_, err := xmysql.NewSession(config)
		xt.KO(t, err)
		xt.Eq(t, "failed loading time location (unknown time zone Foo/Bar)", err.Error())
	})
}

func TestCreateSession(t *testing.T) {
	t.Run("connection times out", func(t *testing.T) {
		expConfig := &xmysql.ConnectConfig{
			Address: "127.0.0.40",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		_, err := xmysql.GetSession(ctx, expConfig)
		xt.KO(t, errors.Unwrap(err))
		xt.Eq(t, "i/o timeout", errors.Unwrap(err).Error())
	})

	t.Run("hello from server and server capabilities", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Username: "user_native",
			Address:  testContext.XPluginAddr,
		}
		config.SetPassword("pwd_user_native")

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.Assert(t, ses.ServerCapabilities() != nil)
		xt.Assert(t, len(ses.ServerCapabilities().AuthMechanisms) > 0)
		// no TLS means no PLAIN
		xt.Assert(t, !xstrings.SliceHas(ses.ServerCapabilities().AuthMechanisms, string(xmysql.AuthMethodPlain)))
		xt.Assert(t, !ses.ServerCapabilities().TLS)
	})

	t.Run("incorrectly connect to conventional MySQL Protocol", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address: testContext.MySQLAddr,
		}

		_, err := xmysql.GetSession(context.Background(), config)
		xt.KO(t, err)
		xt.Eq(t, "wrong protocol [2005:HY000]", err.Error())
	})

	t.Run("connect to something which is not MySQL", func(t *testing.T) {
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(xnet.MustGetLocalhostTCPPort()))
		something, err := net.Listen("tcp", addr)
		xt.OK(t, err)
		defer func() { _ = something.Close() }()

		go func() {
			for {
				conn, err := something.Accept()
				xt.OK(t, err, "could not accept")

				_, _ = conn.Write([]byte{0xff, 0x00, 0x00, 0x00, 0x11, 0x34})
				_ = conn.Close()
				break
			}
		}()

		config := &xmysql.ConnectConfig{
			Address: addr,
		}

		_, err = xmysql.GetSession(context.Background(), config)
		xt.KO(t, err)
		xt.Eq(t, "failed reading message payload (unexpected EOF)", err.Error())
	})

	t.Run("use TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			UseTLS:   true,
			Username: xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.Assert(t, ses.ServerCapabilities() != nil)
		xt.Assert(t, len(ses.ServerCapabilities().AuthMechanisms) > 0)
		xt.Assert(t, xstrings.SliceHas(ses.ServerCapabilities().AuthMechanisms, string(xmysql.AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("use TLS with MySQL Server CA Certificate", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:             testContext.XPluginAddr,
			UseTLS:              true,
			Username:            xxt.UserNative,
			TLSServerCACertPath: "_testdata/mysql_ca.pem",
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.Assert(t, ses.ServerCapabilities() != nil)
		xt.Assert(t, len(ses.ServerCapabilities().AuthMechanisms) > 0)
		xt.Assert(t, xstrings.SliceHas(ses.ServerCapabilities().AuthMechanisms, string(xmysql.AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("use TLS with MySQL Server CA Certificate using environment", func(t *testing.T) {
		xt.OK(t, os.Setenv("xmysql_CA_CERT", "_testdata/mysql_ca.pem"))
		defer func() {
			_ = os.Unsetenv("xmysql_CA_CERT")
		}()

		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			UseTLS:   true,
			Username: xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.Assert(t, ses.ServerCapabilities() != nil)
		xt.Assert(t, len(ses.ServerCapabilities().AuthMechanisms) > 0)
		xt.Assert(t, xstrings.SliceHas(ses.ServerCapabilities().AuthMechanisms, string(xmysql.AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("cannot use PLAIN authn method without TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     false,
			AuthMethod: xmysql.AuthMethodPlain,
			Username:   xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		_, err := xmysql.GetSession(context.Background(), config)
		xt.KO(t, err)
		xt.Assert(t, strings.Contains(err.Error(), "plain text authentication only supported over TLS"))
	})

	t.Run("can use MYSQL41 authn method without TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     false,
			AuthMethod: xmysql.AuthMethodMySQL41,
			Username:   xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		xt.Eq(t, xmysql.AuthMethodMySQL41, ses.AuthMethod())
	})

	t.Run("can use MYSQL41 authn method with TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: xmysql.AuthMethodMySQL41,
			Username:   xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		xt.Eq(t, xmysql.AuthMethodMySQL41, ses.AuthMethod())
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("can use PLAIN authn method with TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: xmysql.AuthMethodPlain,
			Username:   xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		xt.Eq(t, xmysql.AuthMethodPlain, ses.AuthMethod())
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("can use AUTO authn method with TLS", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: xmysql.AuthMethodAuto,
			Username:   xxt.UserCachedSHA256,
		}
		config.SetPassword(xxt.UserCachedSHA256Pwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		defer func() { xt.OK(t, ses.Close()) }()

		xt.Eq(t, xmysql.AuthMethodPlain, ses.AuthMethod()) // AUTO will try plain first, if TLS is enabled
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.ServerCapabilities().TLS)
	})

	t.Run("SHA256 caching after plain authentication using TLS", func(t *testing.T) {
		xt.OK(t, testContext.Server.FlushPrivileges())
		password := xxt.UserCachedSHA256Pwd
		username := xxt.UserCachedSHA256

		config := (&xmysql.ConnectConfig{
			Address: testContext.XPluginAddr,
			UseTLS:  true,
			// AUTO AuthMethod is default which uses PLAIN
			Username: username,
		}).SetPassword(password)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)
		defer func() { xt.OK(t, ses.Close()) }()

		t.Run("use SHA256_MEMORY now works without TLS", func(t *testing.T) {
			config := &xmysql.ConnectConfig{
				Address:    testContext.XPluginAddr,
				UseTLS:     false,
				AuthMethod: xmysql.AuthMethodSHA256Memory, // force
				Username:   username,
			}
			config.SetPassword(password)

			ses, err := xmysql.GetSession(context.Background(), config)
			xt.OK(t, err)
			defer func() { xt.OK(t, ses.Close()) }()

			xt.Eq(t, xmysql.AuthMethodSHA256Memory, ses.AuthMethod())
			xt.Assert(t, !ses.UsesTLS(), "expected no tls.Conn")

			t.Run("again, but now with AUTO", func(t *testing.T) {
				config := &xmysql.ConnectConfig{
					Address:    testContext.XPluginAddr,
					AuthMethod: xmysql.AuthMethodAuto,
					Username:   username,
				}
				config.SetPassword(password)

				ses, err := xmysql.GetSession(context.Background(), config)
				xt.OK(t, err)
				defer func() { xt.OK(t, ses.Close()) }()

				xt.Eq(t, xmysql.AuthMethodSHA256Memory, ses.AuthMethod())
				xt.Assert(t, !ses.UsesTLS(), "expected no tls.Conn")
			})
		})
	})

	t.Run("session uses connection time location", func(t *testing.T) {
		locName := "Europe/Berlin"
		exp, err := time.LoadLocation(locName)
		xt.OK(t, err)

		config := &xmysql.ConnectConfig{
			Address:      testContext.XPluginAddr,
			UseTLS:       true,
			Username:     xxt.UserNative,
			TimeZoneName: "Europe/Berlin",
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		sesZone, err := ses.TimeZone(context.Background())
		xt.OK(t, err)

		xt.Eq(t, exp, sesZone)
	})
}

func TestSession_ExecuteStatement(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	t.Run("numeric data types", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "data_types_numeric"))

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

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
				xt.Assert(t, exp[id].numDecimal.Equal(row.Values[13].(decimal.Decimal)),
					fmt.Sprintf("%s <> %v", exp[id].numDecimal, row.Values[13].(decimal.Decimal)))
				xt.Assert(t, exp[id].numDecimal2.Equal(row.Values[14].(decimal.Decimal)))
				xt.Assert(t, exp[id].numDecimal3.Equal(row.Values[15].(decimal.Decimal)),
					fmt.Sprintf("expected %s == %v", exp[id].numDecimal3, row.Values[15].(decimal.Decimal)))
			})
		}
	})

	t.Run("datetime data types", func(t *testing.T) {
		xt.OK(t, testContext.Server.LoadSQLScript("base", "data_types_datetime"))

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

		locCET, err := time.LoadLocation("CET")
		xt.OK(t, err)

		loc := ses.TimeLocation()

		exp := map[int64]struct {
			dtDATE       time.Time
			dtTIME       time.Duration
			dtDATETIME   time.Time
			dtTIMESTAMP  time.Time
			dtYear       int
			timeLocation *time.Location
		}{
			1: {
				dtDATE:       time.Date(2005, 3, 1, 0, 0, 0, 0, loc),
				dtTIME:       mustParseDuration("8h0m1.123456s"),
				dtDATETIME:   time.Date(2005, 3, 1, 7, 0, 1, 0, loc),
				dtTIMESTAMP:  time.Date(2005, 3, 1, 8, 0, 1, 0, locCET),
				dtYear:       2005,
				timeLocation: locCET,
			},
			2: {
				dtDATE:       time.Date(9999, 12, 31, 0, 0, 0, 0, loc),
				dtTIME:       mustParseDuration("838h59m59s"),
				dtDATETIME:   time.Date(9999, 12, 31, 23, 59, 59, 999999000, loc),
				dtTIMESTAMP:  time.Date(2038, 1, 19, 3, 14, 7, 0, loc),
				dtYear:       1901,
				timeLocation: time.UTC,
			},
			3: {
				dtDATE:       time.Date(1000, 1, 1, 0, 0, 0, 0, loc),
				dtTIME:       mustParseDuration("-838h59m59s"),
				dtDATETIME:   time.Date(1000, 1, 1, 0, 0, 0, 0, loc),
				dtTIMESTAMP:  time.Date(1970, 1, 1, 0, 0, 1, 0, loc),
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

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

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
				xt.Assert(t, xmysql.IsSupportedCollation(res.Columns[1].GetCollation()))
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

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.OK(t, ses.SetActiveSchema(context.Background(), testSchema))

		res, err := ses.ExecuteStatement(context.Background(),
			"INSERT INTO inserts01 (c1) VALUES ('1'),('2')")
		xt.OK(t, err)

		xt.Eq(t, "Records: 2  Duplicates: 0  Warnings: 0", res.StateChanges().ProducedMessage)
		xt.Eq(t, 2, res.RowsAffected())
		xt.Eq(t, 1, res.LastInsertID()) // first generated is returned with multiple values

		res, err = ses.ExecuteStatement(context.Background(),
			"INSERT INTO inserts01 (c1) VALUES ('3')")
		xt.OK(t, err)

		xt.Eq(t, 3, res.LastInsertID())
	})

	t.Run("use arguments and placeholders", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		res, err := ses.ExecuteStatement(context.Background(), "SELECT ?, ?, '?', ? from dual", 1, "one", 3)
		xt.OK(t, err)

		xt.Assert(t, len(res.Rows) == 1, "expected 1 row")
		xt.Eq(t, 1, res.Rows[0].Values[0].(int64))
		xt.Eq(t, "one", res.Rows[0].Values[1].(string))
		xt.Eq(t, "?", res.Rows[0].Values[2].(string))
		xt.Eq(t, 3, res.Rows[0].Values[3].(int64))
	})

	t.Run("zero hour timestamp", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		res, err := ses.ExecuteStatement(context.Background(),
			fmt.Sprintf("SELECT TIMESTAMP('2023-01-01 01:00:00') as ts from dual"))
		xt.OK(t, err)
		xt.Eq(t, 1, len(res.Rows))

		have := res.Rows[0].Values[0].(null.Time)
		have.Time.Equal(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	})
}

func TestSession_CurrentSchema(t *testing.T) {
	t.Run("configure schema and get current schema name", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: xxt.UserNative,
			Schema:   testSchema,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		schema := ses.ActiveSchemaName()
		xt.OK(t, err)
		xt.Eq(t, testSchema, schema)
	})

	t.Run("no current schema in configuration means empty schema name", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		schema := ses.ActiveSchemaName()
		xt.OK(t, err)
		xt.Eq(t, "", schema)
	})
}

func TestSession_SetTimeZone(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	locName := "America/Los_Angeles"
	locUSALA, err := time.LoadLocation(locName)
	xt.OK(t, err)

	t.Run("default UTC when no time zone is set", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
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
		ses, err := xmysql.GetSession(context.Background(), config)
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
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	t.Run("set collation and retrieve", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		exp := xmysql.Collations["utf8mb4_sinhala_ci"]
		xt.OK(t, ses.SetCollation(context.Background(), "utf8mb4_sinhala_ci"))
		have, err := ses.Collation(context.Background())
		xt.OK(t, err)
		xt.Eq(t, exp.ID, have.ID)
		xt.Assert(t, &exp != have)
	})

	t.Run("set invalid collation", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		c := "big5_chinese_ci"
		err = ses.SetCollation(context.Background(), c)
		xt.KO(t, err)
		xt.Eq(t, fmt.Sprintf("failed setting collation ('%s' unsupported)", c), err.Error())
	})
}

func TestSession_PrepareStatement(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	ses, err := xmysql.GetSession(context.Background(), config)
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
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	ses, err := xmysql.GetSession(context.Background(), config)
	xt.OK(t, err)

	t.Run("successfully deallocate", func(t *testing.T) {
		stmt := "SELECT ?"
		prep, err := ses.PrepareStatement(context.Background(), stmt)
		xt.OK(t, err)

		res, err := prep.Execute(context.Background(), 3)
		xt.OK(t, err)

		xt.OK(t, ses.DeallocatePrepareStatement(context.Background(), res.PreparedStatementID()))

		_, err = prep.Execute(context.Background(), 3)
		xxt.AssertMySQLError(t, err, 5110)
	})
}

func TestSession_ActiveSchemaName(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
		Schema:   "pxmysql_tests",
	}
	config.SetPassword(xxt.UserNativePwd)

	t.Run("current database name as configured", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		xt.Eq(t, config.Schema, ses.ActiveSchemaName())
	})

	t.Run("change active schema", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		exp := "pxmysql_tests_a"
		xt.OK(t, ses.SetActiveSchema(context.Background(), exp))
		xt.Eq(t, exp, ses.ActiveSchemaName())
		xt.Eq(t, config.Schema, ses.DefaultSchemaName())
	})
}

func TestSession_Schemas(t *testing.T) {
	config := &xmysql.ConnectConfig{
		Address:  testContext.XPluginAddr,
		Username: xxt.UserNative,
	}
	config.SetPassword(xxt.UserNativePwd)

	t.Run("all databases", func(t *testing.T) {
		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		exp := []string{"information_schema", "performance_schema", "pxmysql_tests", "pxmysql_tests_a"}
		sort.Strings(exp)

		schemas, err := ses.Schemas(context.Background())
		xt.OK(t, err)
		xt.Assert(t, len(schemas) >= len(exp), fmt.Sprintf("expected at least %d", len(exp)))

		var got []string
		for _, s := range schemas {
			got = append(got, s.Name())
		}
		sort.Strings(got)

		for _, n := range exp {
			_, ok := slices.BinarySearch(got, n)
			xt.Assert(t, ok, fmt.Sprintf("expected schema %s to be available", n))
		}
	})
}

func TestSession_CreateSchema(t *testing.T) {
	t.Run("create new schema and drop it", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: "root",
			Password: xstrings.Pointer(testContext.MySQLRootPwd),
		}

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		schemaName := "pxmysql_2839cks829dka"
		schema, err := ses.CreateSchema(context.Background(), schemaName)
		xt.OK(t, err)

		xt.Eq(t, schemaName, schema.Name())

		t.Run("drop schema", func(t *testing.T) {
			xt.OK(t, ses.DropSchema(context.Background(), schemaName))
		})
	})

	t.Run("unprivileged users cannot create schema", func(t *testing.T) {
		config := &xmysql.ConnectConfig{
			Address:  testContext.XPluginAddr,
			Username: xxt.UserNative,
		}
		config.SetPassword(xxt.UserNativePwd)

		ses, err := xmysql.GetSession(context.Background(), config)
		xt.OK(t, err)

		schemaName := "pxmysql_sico29d9kpap21"
		_, err = ses.CreateSchema(context.Background(), schemaName)
		xt.KO(t, err)

		exp := fmt.Sprintf("access denied for user '%s'@'%%' to database '%s' [1044:42000]",
			xxt.UserNative, schemaName)
		xt.Eq(t, exp, errors.Unwrap(err).Error())
	})
}
