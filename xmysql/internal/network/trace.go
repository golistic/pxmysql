// Copyright (c) 2023, Geert JM Vanderkelen

package network

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxnotice"
	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxsql"
)

var (
	traceReadWrites bool
	TraceValues     bool
)

func init() {
	_, traceReadWrites = os.LookupEnv("PXMYSQL_TRACE")
	_, TraceValues = os.LookupEnv("PXMYSQL_TRACE_VALUES")
	if !traceReadWrites {
		TraceValues = false
	}
}

// Trace is used for debugging and is enabled by setting the PYMYSQL_TRACE
// environment variable.
func Trace(action string, msg any, a ...any) {
	if !traceReadWrites || msg == nil {
		return
	}

	var indicator string

	switch action {
	case "w", "write":
		indicator = "\n> write:"
	case "r", "read":
		indicator = "< read"
	case "un", "unhandled":
		indicator = "< unhandled "
	case "error":
		indicator = "< ERROR     "
	case "state":
		indicator = "\t< STATE     "
	default:
		indicator = "< unknown"
	}

	prefix := "\t"

	var s string
	var topic string
	switch v := msg.(type) {
	case *ServerMessage:
		topic = v.ServerMessageType().String()
	case *mysqlxnotice.SessionStateChanged:
		topic = v.GetParam().String()
		doc, err := json.MarshalIndent(v.Value, prefix, "  ")
		if err != nil {
			panic(err)
		}
		if doc[1] != '}' {
			s = fmt.Sprintf("    %s\n", string(doc))
		}
	case *mysqlxsql.StmtExecute:
		s = "  SQL Statement: " + string(v.Stmt) + "\n"
	case proto.Message:
		topic = string(v.ProtoReflect().Descriptor().Name())
		doc, err := json.MarshalIndent(v, prefix, "  ")
		if err != nil {
			panic(err)
		}
		if doc[1] != '}' {
			s = fmt.Sprintf(prefix+"%s\n", string(doc))
		}
	case string:
		topic = v
	default:
		topic = fmt.Sprintf("unhandled %T", msg)
	}

	_, err := fmt.Fprintf(os.Stderr, indicator+" "+topic+"\n"+s)
	if err != nil {
		panic(err)
	}

	if len(a) > 0 {
		_, err := fmt.Fprintf(os.Stderr, prefix+fmt.Sprint(a...)+"\n")
		if err != nil {
			panic(err)
		}
	}
}
