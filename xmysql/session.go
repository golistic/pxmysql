// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxconnection"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxprepare"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxresultset"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsession"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsql"
	"github.com/golistic/pxmysql/mysqlerrors"
	"github.com/golistic/pxmysql/null"
)

// Session uses the Connection configuration to set up a session with
// the MySQL server through the X Plugin. When a session is instantiated
// the authentication start, connection is switched to TLS (if needed).
// All interaction with the server is through this the session.
type Session struct {
	id                 int
	mysqlVersion       string
	cnx                *Connection
	conn               net.Conn
	serverCapabilities *ServerCapabilities
	usedAuthMethod     AuthMethodType
	timeLocation       *time.Location

	preparedStmtCount uint32
}

func newSession(ctx context.Context, cnx *Connection) (*Session, error) {
	ses := &Session{
		cnx:          cnx,
		conn:         nil,
		timeLocation: cnx.timeLocation,
	}

	if err := ses.open(ctx); err != nil {
		return nil, err
	}

	return ses, nil
}

func (ses *Session) String() string {
	var state = "closed"
	if id, err := ses.SessionID(context.Background()); err == nil {
		state = fmt.Sprintf("id=%d", id)
	}
	return fmt.Sprintf("<Session:%s>", state)
}

// Close closes this session.
// It sends the Close-message to the MySQL X Plugin for both session and connection.
func (ses *Session) Close() error {
	if ses == nil {
		return nil
	}
	if err := write(context.Background(), ses, &mysqlxsession.Close{}); err != nil {
		return err
	}

	if err := write(context.Background(), ses, &mysqlxconnection.Close{}); err != nil {
		return fmt.Errorf("failed writing closing message (%w)", err)
	}

	if err := ses.conn.Close(); err != nil {
		return fmt.Errorf("failed closing session (%w)", err)
	}
	return nil
}

// UsesTLS returns whether this session uses TLS.
func (ses *Session) UsesTLS() bool {
	_, ok := ses.conn.(*tls.Conn)
	return ok
}

// AuthMethod returns the used authentication method, which might
// differ from what was configured for the connection.
func (ses *Session) AuthMethod() AuthMethodType {
	return ses.usedAuthMethod
}

func (ses *Session) ExecuteStatement(ctx context.Context, stmt string, args ...any) (*Result, error) {
	if len(args) > 0 {
		var err error
		stmt, err = substitutePlaceholders(stmt, args...)
		if err != nil {
			return nil, err
		}
	}

	if err := write(ctx, ses, &mysqlxsql.StmtExecute{
		Stmt: []byte(stmt),
	}); err != nil {
		return nil, fmt.Errorf("failed writing statement execution (%w)", err)
	}

	res, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.stmtOK
	})
	if err != nil {
		return nil, err
	}
	res.session = ses

	return res, nil
}

// PrepareStatement prepares the statement and returns an instance of Prepared which
// contains the Result instance.
// The ID of the prepared statement can be retrieved using Result.PreparedStatementID().
func (ses *Session) PrepareStatement(ctx context.Context, statement string) (*Prepared, error) {
	stmtID := atomic.AddUint32(&ses.preparedStmtCount, 1)

	if err := write(ctx, ses, &mysqlxprepare.Prepare{
		StmtId: &stmtID,
		Stmt: &mysqlxprepare.Prepare_OneOfMessage{
			Type: mysqlxprepare.Prepare_OneOfMessage_STMT.Enum(),
			StmtExecute: &mysqlxsql.StmtExecute{
				Stmt: []byte(statement),
			},
		},
	}); err != nil {
		return nil, err
	}

	res, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.ok
	})
	if err != nil {
		return nil, err
	}
	res.stmtID = stmtID

	return &Prepared{
		session:         ses,
		result:          res,
		numPlaceholders: len(placeholderIndexes(stmtPlaceholder, statement)),
	}, nil
}

func (ses *Session) deallocatePrepareStatement(ctx context.Context, stmtID uint32) error {
	if err := write(ctx, ses, &mysqlxprepare.Deallocate{
		StmtId: &stmtID,
	}); err != nil {
		return err
	}

	return nil
}

// SetCurrentSchema sets the current schema (database) of this session.
func (ses *Session) SetCurrentSchema(ctx context.Context, name string) error {
	if _, err := ses.ExecuteStatement(ctx, "USE `"+name+"`"); err != nil {
		return err
	}
	return nil
}

// CurrentSchema retrieves the current schema (database) of this session.
func (ses *Session) CurrentSchema(ctx context.Context) (string, error) {
	res, err := ses.ExecuteStatement(ctx, "SELECT SCHEMA()")
	if err != nil {
		return "", err
	}

	if len(res.Rows) > 1 {
		return "", fmt.Errorf("failed getting current schema (too many rows)")
	}

	if res.Rows[0].Values[0] != nil {
		if v := res.Rows[0].Values[0].(null.String); v.Valid {
			return res.Rows[0].Values[0].(null.String).String, err
		}
	}

	return "", nil // no current schema
}

func (ses *Session) SetCollation(ctx context.Context, name string) error {
	c, ok := collations[name]
	if !ok {
		return fmt.Errorf("failed setting collation ('%s' unsupported)", name)
	}

	_, err := ses.ExecuteStatement(ctx, "SET @@collation_connection = ?", c.Name)
	if err != nil {
		return err
	}

	return nil
}

// Collation retrieves this session's collation information.
func (ses *Session) Collation(ctx context.Context) (*Collation, error) {
	res, err := ses.ExecuteStatement(ctx, "SELECT @@collation_connection, VERSION()")
	if err != nil {
		return nil, err
	}

	if len(res.Rows) != 1 || res.Rows[0].Values[0] == nil {
		return nil, fmt.Errorf("failed getting collation of connection (no data)")
	}

	name := res.Rows[0].Values[0].(null.String)
	if !name.Valid {
		return nil, fmt.Errorf("failed getting collation of connection (name was invalid)")
	}
	c, ok := collations[name.String]
	if !ok {
		v := res.Rows[0].Values[1].(null.String)
		return nil, fmt.Errorf("failed getting collation of connection (unsupported '%s'; MySQL v%s)",
			name.String, v.String)
	}
	return &c, err
}

func (ses *Session) SetTimeZone(ctx context.Context, name string) error {
	l, err := time.LoadLocation(name)
	if err != nil {
		return fmt.Errorf("failed loading location (%w)", err)
	}

	if _, err := ses.ExecuteStatement(ctx, "SET @@time_zone=?", name); err != nil {
		return err
	}

	ses.timeLocation = l
	return nil
}

// TimeZone retrieves the session's time zone as Go time.Location.
func (ses *Session) TimeZone(ctx context.Context) (*time.Location, error) {
	res, err := ses.ExecuteStatement(ctx,
		"SELECT IF(@@time_zone='SYSTEM', @@global.system_time_zone, @@session.time_zone)")
	if err != nil {
		return nil, err
	}

	if len(res.Rows) != 1 {
		return nil, fmt.Errorf("failed getting time zone information (too many rows)")
	}

	if res.Rows[0].Values[0] != nil {
		s := res.Rows[0].Values[0].(null.String)
		if s.Valid {
			return time.LoadLocation(s.String)
		}
	}

	return nil, fmt.Errorf("failed getting time zone information (no data)")
}

// SessionID retrieves the MySQL server connection (session) ID.
func (ses *Session) SessionID(ctx context.Context) (int, error) {
	// CurrentSchema retrieves the current schema (database) of this session.
	res, err := ses.ExecuteStatement(ctx, "SELECT CONNECTION_ID()")
	if err != nil {
		return 0, err
	}

	if len(res.Rows) != 1 || res.Rows[0].Values[0] == nil {
		return 0, fmt.Errorf("failed getting session (no data)")
	}

	return int(res.Rows[0].Values[0].(uint64)), nil
}

func (ses *Session) open(ctx context.Context) error {
	var err error

	network := "tcp"
	address := ses.cnx.config.Address
	errCode := mysqlerrors.ClientBadTCPSocket
	if ses.cnx.config.UnixSockAddr != "" {
		network = "unix"
		address = ses.cnx.config.UnixSockAddr
		errCode = mysqlerrors.ClientBadUnixSocket
	}

	ses.conn, err = new(net.Dialer).DialContext(ctx, network, address)
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return mysqlerrors.New(errCode, opErr.Addr,
			fmt.Errorf(strings.Replace(opErr.Err.Error(), "connect: ", "", -1)))
	} else if err != nil {
		return fmt.Errorf("%s (%w)", err.Error(), driver.ErrBadConn)
	}

	defer func() {
		if err != nil {
			_ = ses.conn.Close()
			ses.conn = nil
		}
	}()

	switch v := ses.conn.(type) {
	case *net.TCPConn:
		if err := v.SetKeepAlive(true); err != nil {
			return fmt.Errorf("failed setting keep-alive (%w)", err)
		}
	case *net.UnixConn:
	default:
		return fmt.Errorf("invalid connection implementation (%T)", ses.conn)
	}

	// we do not write anything, but we want the 'hello' from the server
	res, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.notices.serverHello != nil
	})
	switch {
	case res != nil && res.notices.serverHello == nil:
		return fmt.Errorf("did not get hello from server")
	case err != nil:
		return err
	}

	if err := ses.negotiate(ctx); err != nil {
		return err
	}

	if err = ses.authenticate(ctx); err != nil {
		return err
	}

	if err := ses.metaInformation(ctx); err != nil {
		return err
	}

	if err := ses.SetTimeZone(context.Background(), ses.timeLocation.String()); err != nil {
		return err
	}

	return err
}

func (ses *Session) negotiate(ctx context.Context) error {
	if ses.cnx.config.UseTLS {
		if err := write(ctx, ses, &mysqlxconnection.CapabilitiesSet{
			Capabilities: &mysqlxconnection.Capabilities{
				Capabilities: []*mysqlxconnection.Capability{{
					Name:  proto.String("tls"),
					Value: boolAsScalar(true),
				}},
			},
		}); err != nil {
			return fmt.Errorf("failed setting capabilities (%w)", err)
		}

		_, err := ses.handleResult(ctx, func(r *Result) bool {
			return r.ok
		})
		if err != nil {
			return err
		}

		// if TLS is supported, we got an OK, and we do the TLS handshake
		var tlsConfig *tls.Config
		if ses.cnx.config.TLSServerCACertPath != "" {
			if err := addServerCACertFromFile(ses.cnx.config.TLSServerCACertPath); err != nil {
				return err
			}
			tlsConfig = &tls.Config{
				InsecureSkipVerify: false, // explicit to make clear
				RootCAs:            serverCAPool,
				ServerName:         ses.cnx.serverHostname(),
			}
		} else {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		ses.conn = tls.Client(ses.conn, tlsConfig)
		if err := ses.conn.(*tls.Conn).Handshake(); err != nil {
			_ = ses.conn.Close()
			return fmt.Errorf("failed handshake (%w)", err)
		}
	}

	if err := ses.getServerCapabilities(ctx); err != nil {
		return err
	}

	return nil
}

func (ses *Session) authenticate(ctx context.Context) error {
	_, tlsOK := ses.conn.(*tls.Conn)

	var authMethods []AuthMethodType
	if ses.cnx.config.AuthMethod == AuthMethodAuto {
		if ses.serverCapabilities != nil && ses.serverCapabilities.TLS {
			authMethods = append(authMethods, AuthMethodPlain)
		}
		authMethods = append(authMethods, defaultAuthMethods...)
	} else if ses.usedAuthMethod == AuthMethodPlain && !tlsOK {
		return fmt.Errorf("plain text authentication only supported over TLS")
	} else {
		authMethods = []AuthMethodType{ses.cnx.config.AuthMethod}
	}

	var authErr error
	var canRetry bool

	for _, method := range authMethods {
		canRetry, authErr = ses.authenticateWith(ctx, method)
		switch {
		case authErr == nil:
			return nil
		case !canRetry && authErr != nil:
			return authErr
		}
	}

	return authErr
}

func (ses *Session) authenticateWith(ctx context.Context, method AuthMethodType) (bool, error) {
	ses.usedAuthMethod = method

	if method == AuthMethodPlain {
		return ses.authenticatePlain(ctx)
	}

	// send AuthenticateStart
	if err := write(ctx, ses, &mysqlxsession.AuthenticateStart{
		MechName: proto.String(string(method)),
	}); err != nil {
		return false, fmt.Errorf("failed starting authentication (%w)", err)
	}

	res, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.authChallenge != nil
	})
	if err != nil {
		return true, err
	}

	// create authentication data
	var authDataFunc func(authn) ([]byte, error)
	switch method {
	case AuthMethodMySQL41:
		authDataFunc = authMySQL41Data
	case AuthMethodSHA256Memory:
		authDataFunc = authSHA256Data
	default:
		return true, fmt.Errorf("unsupported authentication method '%s'", method)
	}

	authData, err := authDataFunc(authn{
		username:  ses.cnx.config.Username,
		password:  ses.cnx.password,
		challenge: res.authChallenge,
		schema:    ses.cnx.config.Schema,
	})
	if err != nil {
		return false, err
	}

	// send AuthenticateContinue
	if err := write(ctx, ses, &mysqlxsession.AuthenticateContinue{
		AuthData: authData,
	}); err != nil {
		return false, fmt.Errorf("failed continuing authentication (%w)", err)
	}

	return ses.finishAuthentication(ctx)
}

func (ses *Session) authenticatePlain(ctx context.Context) (bool, error) {
	_, tlsOK := ses.conn.(*tls.Conn)

	if !tlsOK {
		return true, fmt.Errorf("plain text authentication only supported over TLS")
	}

	authData, err := authMySQLPlain(authn{
		username: ses.cnx.config.Username,
		password: ses.cnx.password,
		schema:   ses.cnx.config.Schema,
	})
	if err != nil {
		return false, err
	}

	if err := write(ctx, ses, &mysqlxsession.AuthenticateStart{
		MechName: proto.String(string(AuthMethodPlain)),
		AuthData: authData,
	}); err != nil {
		return false, err
	}

	return ses.finishAuthentication(ctx)
}

func (ses *Session) finishAuthentication(ctx context.Context) (bool, error) {
	result, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.authOK
	})
	if err != nil {
		return true, err // try other methods
	}

	return !result.authOK, nil // try other methods when needed
}

func (ses *Session) getServerCapabilities(ctx context.Context) error {
	if err := write(ctx, ses, &mysqlxconnection.CapabilitiesGet{}); err != nil {
		return fmt.Errorf("failed getting capabilities (%w)", err)
	}

	res, err := ses.handleResult(ctx, func(r *Result) bool {
		return r.serverCapabilities != nil
	})
	if err != nil {
		return err
	}

	ses.serverCapabilities = res.serverCapabilities
	return nil
}

func (ses *Session) metaInformation(ctx context.Context) error {
	// CurrentSchema retrieves the current schema (database) of this session.
	res, err := ses.ExecuteStatement(ctx, "SELECT VERSION(), CONNECTION_ID()")
	if err != nil {
		return err
	}

	if len(res.Rows) != 1 {
		return fmt.Errorf("failed getting meta information (no rows)")
	}

	ses.mysqlVersion = res.Rows[0].Values[0].(string)
	ses.id = int(res.Rows[0].Values[1].(uint64))

	return nil
}

func (ses *Session) handleResult(ctx context.Context, doneWhen doneWhenFunc) (*Result, error) {
	result := &Result{session: ses}

	// force time zone
	if ses.timeLocation != nil {
		ctx = SetContextTimeLocation(ctx, ses.timeLocation)
	} else {
		ctx = SetContextTimeLocation(ctx, defaultTimeLocation)
	}

	for done := false; !done; {
		msg, err := read(ctx, ses.conn)
		switch {
		case err == io.EOF:
			done = true
			continue
		case err != nil:
			return nil, err
		}

		msgType := msg.ServerMessageType()
		switch msgType {
		case mysqlx.ServerMessages_OK:
			result.ok = true
		case mysqlx.ServerMessages_ERROR:
			return nil, mysqlerrors.NewFromServerMessage(msg)
		case mysqlx.ServerMessages_CONN_CAPABILITIES:
			result.serverCapabilities, err = NewServerCapabilitiesFromMessage(msg)
			if err != nil {
				return nil, err
			}
		case mysqlx.ServerMessages_SESS_AUTHENTICATE_CONTINUE:
			m := &mysqlxsession.AuthenticateContinue{}
			if err := msg.Unmarshall(m); err != nil {
				return nil, fmt.Errorf("failed unmarshalling %s (%w)", msgType.String(), err)
			}
			result.authChallenge = m.AuthData
		case mysqlx.ServerMessages_SESS_AUTHENTICATE_OK:
			result.authOK = true
		case mysqlx.ServerMessages_NOTICE:
			if err := result.notices.add(msg); err != nil {
				return nil, err
			}
		case mysqlx.ServerMessages_RESULTSET_COLUMN_META_DATA:
			m := &mysqlxresultset.ColumnMetaData{}
			if err := msg.Unmarshall(m); err != nil {
				return nil, fmt.Errorf("failed unmarshalling '%s' (%w)", msgType.String(), err)
			}
			result.Columns = append(result.Columns, m)
		case mysqlx.ServerMessages_RESULTSET_ROW:
			if err := result.readRow(ctx, msg); err != nil {
				return nil, err
			}
		case mysqlx.ServerMessages_RESULTSET_FETCH_DONE:
			result.fetchDone = true
		case mysqlx.ServerMessages_RESULTSET_FETCH_DONE_MORE_RESULTSETS:
			result.fetchDoneMoreResults = true
		case mysqlx.ServerMessages_SQL_STMT_EXECUTE_OK:
			result.stmtOK = true
		case mysqlx.ServerMessages_RESULTSET_FETCH_DONE_MORE_OUT_PARAMS:
			result.fetchDoneMoreOutParams = true
		default:
			trace("unhandled", msg)
		}

		if doneWhen != nil {
			done = doneWhen(result)
		}
	}

	return result, nil
}
