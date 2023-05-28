// Copyright (c) 2022, Geert JM Vanderkelen

package mysqlerrors

// MySQL Client errors as found in the MySQL manual under
// https://dev.mysql.com/doc/mysql-errors/8.0/en/client-error-reference.html.
// Names have been altered so that they make more sense. For example,
// CR_CONNECTION_ERROR became ClientBadUnixSocket.
const (
	ClientUnknown           = 2000
	ClientBadUnixSocket     = 2002
	ClientBadTCPSocket      = 2005
	ClientWrongProtocol     = 2007
	ClientNetPacketTooLarge = 2020
)

var mysqlClientErrors = map[int]Error{
	ClientUnknown: { // 2000
		Message:  "unknown error",
		Code:     ClientUnknown,
		SQLState: "HY000",
	},
	ClientBadUnixSocket: { // 2002
		Message:  "cannot connect to local MySQL server through socket '%s' (%w)",
		Code:     ClientBadUnixSocket,
		SQLState: "HY000",
	},
	ClientBadTCPSocket: { // 2005
		Message:  "unknown MySQL server host '%s' (%w)",
		Code:     ClientBadTCPSocket,
		SQLState: "HY000",
	},
	ClientWrongProtocol: { // 2007
		Message:  "wrong protocol",
		Code:     ClientBadTCPSocket,
		SQLState: "HY000",
	},
	ClientNetPacketTooLarge: { // 2020
		Message:  "got packet bigger than 'mysqlx_max_allowed_packet' bytes",
		Code:     ClientNetPacketTooLarge,
		SQLState: "HY000",
	},
}
