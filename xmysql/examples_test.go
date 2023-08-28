// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golistic/pxmysql/null"
	"github.com/golistic/pxmysql/xmysql"
)

func ExampleConnection_NewSession_auto_notls() {
	config := &xmysql.ConnectConfig{
		Address:  "127.0.0.1:53360", // see _support/pxmysql-compose/docker-compose.yml
		Username: "user_native",
	}
	config.SetPassword("pwd_user_native")

	session, err := xmysql.CreateSession(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TLS:", session.UsesTLS())
	// Output: TLS: false
}

func ExampleConnection_NewSession_plain_withtls() {
	config := &xmysql.ConnectConfig{
		Address:    "127.0.0.1:53360", // see _support/pxmysql-compose/docker-compose.yml
		AuthMethod: xmysql.AuthMethodPlain,
		UseTLS:     true,
		Username:   "user_native",
	}
	config.SetPassword("pwd_user_native")

	session, err := xmysql.CreateSession(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TLS:", session.UsesTLS())
	fmt.Println("Auth Method:", config.AuthMethod)
	// Output:
	// TLS: true
	// Auth Method: PLAIN
}

func ExampleSession_ExecuteStatement() {
	config := &xmysql.ConnectConfig{
		Address:    "127.0.0.1:53360", // see _support/pxmysql-compose/docker-compose.yml
		AuthMethod: xmysql.AuthMethodPlain,
		UseTLS:     true,
		Username:   "user_native",
	}
	config.SetPassword("pwd_user_native")

	session, err := xmysql.CreateSession(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	q := "SELECT ?, STR_TO_DATE('2005-03-01 07:00:01', '%Y-%m-%d %H:%i:%s')"
	res, err := session.ExecuteStatement(context.Background(), q, "started")
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range res.Rows {
		fmt.Printf("%s at %s\n", row.Values[0].(string), row.Values[1].(null.Time).Time)
	}

	// Output:
	// started at 2005-03-01 07:00:01 +0000 UTC
}
