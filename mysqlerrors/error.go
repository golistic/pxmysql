// Copyright (c) 2022, Geert JM Vanderkelen

package mysqlerrors

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/golistic/pxmysql/interfaces"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxnotice"
)

var (
	ErrContextDeadlineExceeded = errors.New("context deadline exceeded")
)

// Error holds information of MySQL returned error and is implementing
// the Go error interface.
type Error struct {
	Inner        error
	Message      string
	Code         int
	SQLState     string
	Severity     int
	Parameters   []any
	ExtraMessage string
}

// New instantiates an Error object with code being the MySQL
// client or server error code. When the message contains placeholders (is
// a format specifier), params are used to interpolate them.
// When one of the params is an error, it is used as a wrapped error.
// Panics when params contains more than one value which is error-type.
func New(code int, params ...any) *Error {
	e, have := mysqlClientErrors[code]
	if !have {
		panic(fmt.Sprintf("error code %d not registered in mysqlClientErrors", code))
	}

	for _, param := range params {
		if p, ok := param.(error); ok {
			if e.Inner != nil {
				panic("only one parameter can be of error type")
			}
			e.Inner = p
		}
	}

	e.Parameters = params
	return &e
}

func (e *Error) Unwrap() error {
	return e.Inner
}

// Error is the string representation of the error. Messages look the same, but
// they are differently formatted than the MySQL Client message to conform Go
// best practices.
func (e *Error) Error() string {
	msg := e.Message

	if len(e.Parameters) > 0 {
		if e.Inner != nil {
			msg = fmt.Errorf(e.Message, e.Parameters...).Error()
		} else {
			msg = fmt.Sprintf(e.Message, e.Parameters...)
		}
	}

	if e.ExtraMessage != "" {
		msg += " (" + e.ExtraMessage + ")"
	}

	if len(msg) > 2 {
		r := []rune(msg)
		if !strings.HasPrefix(msg, "MySQL") {
			r[0] = unicode.ToLower(r[0])
		}
		msg = string(r)
	}

	return fmt.Sprintf("%s [%d:%s]", msg, e.Code, e.SQLState)
}

// NewFromServerMessage takes msg and transforms it into an Error.
func NewFromServerMessage(msg interfaces.ServerMessager) error {
	myErr := &mysqlx.Error{}

	if err := msg.Unmarshall(myErr); err != nil {
		return err
	}

	return &Error{
		Message:  myErr.GetMsg(),
		Code:     int(myErr.GetCode()),
		SQLState: myErr.GetSqlState(),
		Severity: int(myErr.GetSeverity()),
	}
}

// MySQLWarning holds information of MySQL returned warning and is implementing
// the Go error interface.
type MySQLWarning struct {
	Message string
	Code    int
	Level   string
}

// Error is the string representation of the warning, mimicking how MySQL would
// show them.
func (w *MySQLWarning) Error() string {
	return fmt.Sprintf("%s %d: %s", w.Level, w.Code, w.Message)
}

// NewFromWarning takes msg and transforms it into a MySQLWarning.
func NewFromWarning(msg *mysqlxnotice.Warning) error {
	return &MySQLWarning{
		Level:   msg.Level.String(),
		Message: msg.GetMsg(),
		Code:    int(msg.GetCode()),
	}
}
