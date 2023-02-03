// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xnet"
	"github.com/geertjanvdk/xkit/xt"
	"github.com/geertjanvdk/xkit/xutil"
)

const (
	userNative    = "user_native"
	userNativePwd = "pwd_user_native"
)

const (
	userCachedSHA256    = "user_sha256"
	userCachedSHA256Pwd = "pwd_user_sha256"
)

func TestNewConnection(t *testing.T) {
	t.Run("set configuration", func(t *testing.T) {
		expConfig := &ConnectConfig{
			Address:  "127.0.0.5",
			Username: "toor",
			UseTLS:   false,
		}
		expConfig.SetPassword("toor")

		cnx, err := NewConnection(expConfig)
		xt.OK(t, err)

		xt.Assert(t, !cnx.config.Password.Valid)
		xt.Eq(t, cnx.config.AuthMethod, DefaultConnectConfig.AuthMethod)
	})

	t.Run("default configuration", func(t *testing.T) {
		cnx, err := NewConnection(nil)
		xt.OK(t, err)

		xt.Assert(t, !cnx.config.Password.Valid)
		xt.Eq(t, cnx.config.Address, DefaultConnectConfig.Address)
		xt.Eq(t, cnx.config.Username, DefaultConnectConfig.Username)
		xt.Eq(t, cnx.password, DefaultConnectConfig.Password.String)
		xt.Eq(t, cnx.config.AuthMethod, DefaultConnectConfig.AuthMethod)
	})

	t.Run("MySQL X Plugin is reachable", func(t *testing.T) {
		expConfig := &ConnectConfig{
			Address: testContext.XPluginAddr,
		}
		cnx, err := NewConnection(expConfig)
		xt.OK(t, err)

		xt.Assert(t, cnx.isReachable())
	})

	t.Run("MySQL X Plugin is not reachable", func(t *testing.T) {
		expConfig := &ConnectConfig{
			Address: "127.0.0.40",
		}
		cnx, err := NewConnection(expConfig)
		xt.OK(t, err)
		xt.Eq(t, "127.0.0.40:"+defaultXPluginPort, cnx.config.Address)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err = cnx.NewSession(ctx)
		xt.KO(t, err)
		xt.KO(t, errors.Unwrap(err))
		xt.Eq(t, "i/o timeout", errors.Unwrap(err).Error())
	})

	t.Run("address getting defaults if needed", func(t *testing.T) {
		var cases = []struct {
			have string
			exp  string
		}{
			{have: "127.0.0.1", exp: "127.0.0.1:" + defaultXPluginPort},
			{have: ":" + defaultXPluginPort, exp: defaultXPluginHost + ":" + "33060"},
			{have: ":12453", exp: defaultXPluginHost + ":" + "12453"},
			{have: "0.0.0.0:" + defaultXPluginPort, exp: "0.0.0.0:" + defaultXPluginPort},
			{have: "", exp: defaultXPluginHost + ":" + defaultXPluginPort},
		}

		for _, c := range cases {
			t.Run(c.exp, func(t *testing.T) {
				cnx, err := NewConnection(&ConnectConfig{
					Address: c.have,
				})
				xt.OK(t, err)
				xt.Eq(t, c.exp, cnx.config.Address)
			})
		}
	})

	t.Run("valid time zone", func(t *testing.T) {
		locName := "Europe/Berlin"
		exp, err := time.LoadLocation(locName)
		xt.OK(t, err)

		config := &ConnectConfig{
			TimeZoneName: locName,
		}

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		xt.Eq(t, exp, cnx.timeLocation)
	})

	t.Run("invalid time zone", func(t *testing.T) {
		config := &ConnectConfig{
			TimeZoneName: "Foo/Bar",
		}

		_, err := NewConnection(config)
		xt.KO(t, err)
		xt.Eq(t, "failed loading time location (unknown time zone Foo/Bar)", err.Error())
	})
}

func TestConnection_NewSession(t *testing.T) {
	t.Run("hello from server and server capabilities", func(t *testing.T) {
		config := &ConnectConfig{
			Username: "user_native",
			Address:  testContext.XPluginAddr,
		}
		config.SetPassword("pwd_user_native")
		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.Assert(t, ses.serverCapabilities != nil)
		xt.Assert(t, len(ses.serverCapabilities.AuthMechanisms) > 0)
		// no TLS means no PLAIN
		xt.Assert(t, !xutil.HasString(ses.serverCapabilities.AuthMechanisms, string(AuthMethodPlain)))
		xt.Assert(t, !ses.serverCapabilities.TLS)
	})

	t.Run("incorrectly connect to MySQL Protocol", func(t *testing.T) {
		config := &ConnectConfig{
			Address: testContext.MySQLAddr,
		}
		cnx, err := NewConnection(config)
		xt.OK(t, err)

		_, err = cnx.NewSession(context.Background())
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

		config := &ConnectConfig{
			Address: addr,
		}
		cnx, err := NewConnection(config)
		xt.OK(t, err)

		_, err = cnx.NewSession(context.Background())
		xt.KO(t, err)
		xt.Eq(t, "failed reading message payload (unexpected EOF)", err.Error())
	})

	t.Run("use TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:  testContext.XPluginAddr,
			UseTLS:   true,
			Username: userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.Assert(t, ses.serverCapabilities != nil)
		xt.Assert(t, len(ses.serverCapabilities.AuthMechanisms) > 0)
		xt.Assert(t, xutil.HasString(ses.serverCapabilities.AuthMechanisms, string(AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("use TLS with MySQL Server CA Certificate", func(t *testing.T) {
		config := &ConnectConfig{
			Address:             testContext.XPluginAddr,
			UseTLS:              true,
			Username:            userNative,
			TLSServerCACertPath: "_testdata/mysql_ca.pem",
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.Assert(t, ses.serverCapabilities != nil)
		xt.Assert(t, len(ses.serverCapabilities.AuthMechanisms) > 0)
		xt.Assert(t, xutil.HasString(ses.serverCapabilities.AuthMechanisms, string(AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("use TLS with MySQL Server CA Certificate using environment", func(t *testing.T) {
		xt.OK(t, os.Setenv("xmysql_CA_CERT", "_testdata/mysql_ca.pem"))
		defer func() {
			_ = os.Unsetenv("xmysql_CA_CERT")
		}()

		config := &ConnectConfig{
			Address:  testContext.XPluginAddr,
			UseTLS:   true,
			Username: userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		xt.Assert(t, ses.serverCapabilities != nil)
		xt.Assert(t, len(ses.serverCapabilities.AuthMechanisms) > 0)
		xt.Assert(t, xutil.HasString(ses.serverCapabilities.AuthMechanisms, string(AuthMethodPlain)))
		xt.Assert(t, ses.UsesTLS(), "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("cannot use PLAIN authn method without TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     false,
			AuthMethod: AuthMethodPlain,
			Username:   userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		_, err = cnx.NewSession(context.Background())
		xt.KO(t, err)
		xt.Assert(t, strings.Contains(err.Error(), "plain text authentication only supported over TLS"))
	})

	t.Run("can use MYSQL41 authn method without TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     false,
			AuthMethod: AuthMethodMySQL41,
			Username:   userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.Eq(t, AuthMethodMySQL41, ses.usedAuthMethod)
	})

	t.Run("can use MYSQL41 authn method with TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: AuthMethodMySQL41,
			Username:   userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.Eq(t, AuthMethodMySQL41, ses.usedAuthMethod)
		_, ok := ses.conn.(*tls.Conn)
		xt.Assert(t, ok, "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("can use PLAIN authn method with TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: AuthMethodPlain,
			Username:   userNative,
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.Eq(t, AuthMethodPlain, ses.usedAuthMethod)
		_, ok := ses.conn.(*tls.Conn)
		xt.Assert(t, ok, "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("can use AUTO authn method with TLS", func(t *testing.T) {
		config := &ConnectConfig{
			Address:    testContext.XPluginAddr,
			UseTLS:     true,
			AuthMethod: AuthMethodAuto,
			Username:   userCachedSHA256,
		}
		config.SetPassword(userCachedSHA256Pwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)
		defer func() { xt.OK(t, ses.Close()) }()

		xt.Eq(t, AuthMethodPlain, ses.usedAuthMethod) // AUTO will try plain first, if TLS is enabled
		_, ok := ses.conn.(*tls.Conn)
		xt.Assert(t, ok, "expected tls.Conn")
		xt.Assert(t, ses.serverCapabilities.TLS)
	})

	t.Run("SHA256 caching after plain authentication using TLS", func(t *testing.T) {
		xt.OK(t, testContext.Server.FlushPrivileges())
		password := userCachedSHA256Pwd
		username := userCachedSHA256

		config := (&ConnectConfig{
			Address: testContext.XPluginAddr,
			UseTLS:  true,
			// AUTO AuthMethod is default which uses PLAIN
			Username: username,
		}).SetPassword(password)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)
		defer func() { xt.OK(t, ses.Close()) }()

		t.Run("use SHA256_MEMORY now works without TLS", func(t *testing.T) {
			config := &ConnectConfig{
				Address:    testContext.XPluginAddr,
				UseTLS:     false,
				AuthMethod: AuthMethodSHA256Memory, // force
				Username:   username,
			}
			config.SetPassword(password)

			cnx, err := NewConnection(config)
			xt.OK(t, err)

			ses, err := cnx.NewSession(context.Background())
			xt.OK(t, err)
			defer func() { xt.OK(t, ses.Close()) }()

			xt.Eq(t, AuthMethodSHA256Memory, ses.usedAuthMethod)
			_, ok := ses.conn.(*tls.Conn)
			xt.Assert(t, !ok, "expected no tls.Conn")

			t.Run("again, but now with AUTO", func(t *testing.T) {
				config := &ConnectConfig{
					Address:    testContext.XPluginAddr,
					AuthMethod: AuthMethodAuto,
					Username:   username,
				}
				config.SetPassword(password)

				cnx, err := NewConnection(config)
				xt.OK(t, err)

				ses, err := cnx.NewSession(context.Background())
				xt.OK(t, err)
				defer func() { xt.OK(t, ses.Close()) }()

				xt.Eq(t, AuthMethodSHA256Memory, ses.usedAuthMethod)
				_, ok := ses.conn.(*tls.Conn)
				xt.Assert(t, !ok, "expected no tls.Conn")
			})
		})
	})

	t.Run("session uses connection time location", func(t *testing.T) {
		locName := "Europe/Berlin"
		exp, err := time.LoadLocation(locName)
		xt.OK(t, err)

		config := &ConnectConfig{
			Address:      testContext.XPluginAddr,
			UseTLS:       true,
			Username:     userNative,
			TimeZoneName: "Europe/Berlin",
		}
		config.SetPassword(userNativePwd)

		cnx, err := NewConnection(config)
		xt.OK(t, err)

		ses, err := cnx.NewSession(context.Background())
		xt.OK(t, err)

		sesZone, err := ses.TimeZone(context.Background())
		xt.OK(t, err)

		xt.Eq(t, exp, sesZone)
	})
}
