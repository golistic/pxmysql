// Copyright (c) 2022, 2023, Geert JM Vanderkelen

package xmysql

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/golistic/pxmysql/decimal"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlx"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxresultset"
	"github.com/golistic/pxmysql/mysqlerrors"
	"github.com/golistic/pxmysql/null"
)

const flagNotNull = 0x0010

type doneWhenFunc = func(r *Result) bool

type Row struct {
	Values []any
}

// NewRow returns a new instance of Row with values slice set to
// a length of nrColumns.
func NewRow(nrColumns int) *Row {
	if nrColumns < 1 {
		panic("impossible number of columns for row")
	}

	r := &Row{
		Values: make([]any, nrColumns),
	}
	return r
}

type Result struct {
	ok                     bool
	authOK                 bool
	stmtOK                 bool
	fetchDone              bool
	fetchDoneMoreResults   bool
	fetchDoneMoreOutParams bool
	authChallenge          []byte
	notices                notices
	serverCapabilities     *ServerCapabilities
	session                *Session

	Row             *Row
	Rows            []*Row
	Columns         []*mysqlxresultset.ColumnMetaData
	ProducedMessage string
	stmtID          uint32
}

func (rs *Result) Warnings() []error {
	if len(rs.notices.warnings) == 0 {
		return nil
	}

	errors := make([]error, len(rs.notices.warnings))
	for i, n := range rs.notices.warnings {
		errors[i] = mysqlerrors.NewFromWarning(n)
	}

	return errors
}

func (rs *Result) LastInsertID() uint64 {
	return rs.notices.stateChanges.GeneratedInsertID
}

func (rs *Result) RowsAffected() uint64 {
	return rs.notices.stateChanges.RowsAffected
}

func (rs *Result) PreparedStatementID() uint32 {
	return rs.stmtID
}

// StateChanges returns the object contain eventual state changes.
func (n *Result) StateChanges() stateChanges {
	return n.notices.stateChanges
}

// FetchRow fetches the next row for unbuffered results and stores it in Result.Row.
// When no more row is available nil is returned as well as a nil error (not the mysqlerrors.ErrResultNoMore error).
func (rs *Result) FetchRow(ctx context.Context) error {
	if rs.session == nil {
		return fmt.Errorf("failed fetching row (%w)", fmt.Errorf("no session"))
	}

	if rs.fetchDone {
		return nil
	}

	msg, err := read(ctx, rs.session.conn)
	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	}

	switch msg.ServerMessageType() {
	case mysqlx.ServerMessages_RESULTSET_ROW:
		// force time zone
		if rs.session.timeLocation != nil {
			ctx = SetContextTimeLocation(ctx, rs.session.timeLocation)
		} else {
			ctx = SetContextTimeLocation(ctx, defaultTimeLocation)
		}

		if err := rs.readRow(ctx, msg); err != nil {
			return fmt.Errorf("failed fetching row (%w)", err)
		}
		return nil
	case mysqlx.ServerMessages_RESULTSET_FETCH_DONE:
		rs.Row = nil
		return nil
	default:
		trace("unhandled", msg)
		return nil
	}
}

func (rs *Result) readRow(ctx context.Context, msg *serverMessage) error {
	if msg == nil {
		panic("serverMessage cannot be nil")
	}

	if len(rs.Columns) == 0 {
		return fmt.Errorf("no column metadata")
	}

	r := &mysqlxresultset.Row{}
	if err := msg.Unmarshall(r); err != nil {
		return fmt.Errorf("failed unmarshalling '%s' (%w)", msg.ServerMessageType().String(), err)
	}

	row := NewRow(len(rs.Columns))
	for i, value := range r.Field {
		d, err := rs.decodeValue(ctx, value, rs.Columns[i])
		if err != nil {
			return err
		}
		row.Values[i] = d
	}

	rs.Rows = append(rs.Rows, row)
	return nil
}

func (rs *Result) decodeValue(ctx context.Context, value []byte, column *mysqlxresultset.ColumnMetaData) (any, error) {
	var goValue any
	valid := len(value) > 0

	if traceValues {
		var notNull string
		if column.GetFlags()&flagNotNull > 0 {
			notNull = " NOT NULL"
		}
		fmt.Printf("[%s%s]\n", column.Type.String(), notNull)
		fmt.Print(hex.Dump(value))
	}

	switch *column.Type {
	case mysqlxresultset.ColumnMetaData_SINT:
		var v int64
		if len(value) > 0 {
			var n int
			v, n = binary.Varint(value)
			if n != len(value) {
				return nil, fmt.Errorf("failed decoding %#v as SINT", value)
			}
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Int64{
				Int64: v,
				Valid: valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_UINT:
		var v uint64
		if len(value) > 0 {
			var n int
			v, n = binary.Uvarint(value)
			if n != len(value) {
				return nil, fmt.Errorf("failed decoding %#v as UINT", value)
			}
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Uint64{
				Uint64: v,
				Valid:  valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_DOUBLE:
		var v float64
		if len(value) > 0 {
			v = math.Float64frombits(binary.LittleEndian.Uint64(value))
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Float64{
				Float64: v,
				Valid:   valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_FLOAT:
		var v float32
		if len(value) > 0 {
			v = math.Float32frombits(binary.LittleEndian.Uint32(value))
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Float32{
				Float32: v,
				Valid:   len(value) > 0,
			}
		}

	case mysqlxresultset.ColumnMetaData_BYTES, mysqlxresultset.ColumnMetaData_ENUM:
		_, isString := collationIDs[column.GetCollation()]

		if isString {
			var v string
			if valid {
				v = string(value[:len(value)-1])
			}
			if column.GetFlags()&flagNotNull > 0 {
				goValue = v
			} else {
				goValue = null.String{
					String: v,
					Valid:  valid,
				}
			}
		} else {
			var v []byte
			if valid {
				v = value[:len(value)-1]
			}
			if column.GetFlags()&flagNotNull > 0 {
				goValue = v
			} else {
				goValue = null.Bytes{
					Bytes: v,
					Valid: valid,
				}
			}
		}

	case mysqlxresultset.ColumnMetaData_TIME:
		var v time.Duration

		if len(value) > 0 {
			negative := value[0] == 1
			h, n := binary.Uvarint(value[1:])

			m := int64(value[n+1])
			s := int64(value[n+2])
			var us int64

			// decode microseconds as nanoseconds if available
			if len(value) > 1+n+2 { // 1+ == is first byte denoting sign
				us = func() int64 {
					v, _ := binary.Uvarint(value[n+3:])
					return int64(v) * 1000
				}()
			}

			t := (((int64(h) * 3600) + (m * 60) + s) * int64(1000000000)) + us // duration is in nanoseconds
			if negative {
				t *= -1
			}

			v = time.Duration(t)
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Duration{
				Duration: v,
				Valid:    valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_DATETIME:
		var v time.Time
		if len(value) > 0 {
			parts := [7]int{}

			// more verbose, but avoiding use of bytes.NewReader and binary.ReadUvarint
			parts[0] = func() int {
				v, _ := binary.Uvarint(value[0:2]) // always 2 bytes
				return int(v)
			}()
			parts[1] = int(value[2])
			parts[2] = int(value[3])

			// decode hour if available
			if len(value) > 4 {
				parts[3] = int(value[4])
			}

			// decode minutes if available
			if len(value) > 5 {
				parts[4] = int(value[5])
			}

			// decode seconds if available
			if len(value) > 6 {
				parts[5] = int(value[6])
			}

			// decode microseconds as nanoseconds if available
			if len(value) > 7 {
				parts[6] = func() int {
					v, _ := binary.Uvarint(value[7:])
					return int(v) * 1000
				}()
			}

			v = time.Date(
				parts[0], time.Month(parts[1]), parts[2],
				parts[3], parts[4], parts[5],
				parts[6], ContextTimeLocation(ctx))
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Time{
				Time:  v,
				Valid: valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_SET:
		var v []string
		emptySet := valid && (value[0] == 0x01 && len(value) == 1)

		if valid && !emptySet {
			for n := 0; n < len(value); {
				l, s := binary.Uvarint(value[n:])
				v = append(v, string(value[n+s:n+s+int(l)]))
				n += s + int(l)
			}
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Strings{
				Strings: v,
				Valid:   valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_BIT:
		var v uint64

		if valid {
			var n int
			if v, n = binary.Uvarint(value); n < 0 {
				return nil, fmt.Errorf("failed decoding %#v as BIT", value)
			}
		}

		if column.GetFlags()&flagNotNull > 0 {
			goValue = v
		} else {
			goValue = null.Uint64{
				Uint64: v,
				Valid:  valid,
			}
		}

	case mysqlxresultset.ColumnMetaData_DECIMAL:
		var err error
		if goValue, err = decimal.NewDecimalFromBCD(value); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported column type %s; value %#v", column.Type.String(), value)
	}

	if goValue == nil {
		return nil, nil
	}

	return goValue, nil
}
