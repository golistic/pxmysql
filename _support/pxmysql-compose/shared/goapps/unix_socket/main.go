// Copyright (c) 2022, Geert JM Vanderkelen

/*
This application is executed within the MySQL container by the pxmysql
Go tests.

The MySQL sock-file within the Docker container cannot be accessed through
a Docker volume. Copying in an application which uses pxmysql is therefor
the only way to automate testing of UNIX socket file support.
*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golistic/pxmysql/register"
)

func main() {
	// credentials are the once used when running pxmysql tests
	db, err := sql.Open("pxmysql", "root:rootpwd@unix(/var/lib/mysql/mysqlx.sock)")
	if err != nil {
		log.Fatalln("open:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("ping:", err)
	}

	var version string
	if err := db.QueryRowContext(context.Background(), "SELECT VERSION()").Scan(&version); err != nil {
		log.Fatalln("query row:", err)
	}

	fmt.Println(version)
}
