pxmysql - Driver using the MySQL X Protocol
===========================================

Copyright (c) 2022, 2023, Geert JM Vanderkelen

<div>
  <img alt="Go: 1.21" src="_badges/go-version.svg">
  <img alt="license: MIT" src="_badges/license.svg">
</div>

The Go pxmysql package implements the MySQL X Protocol and provides a Go `sql/driver`
which uses it. The X Protocol communicates with MySQL using TCP port 33060 (default)
using structured data serialized using Protocol Buffers.

Note that the MySQL X Protocol is an alternative, an extension of the
conventional well known text-based MySQL protocol.  
If you are looking for a driver using the MySQL port 3306, please use the
excellent [github.com/go-sql-driver/mysql][3].

Installation
------------

The `pxmysql` package supports Go 1.19 and greater.

```go get -u github.com/golistic/pxmysql```

Quick Start
-----------

The below code connects to the MySQL server and gets the current time.

Note! This is MySQL Protocol X; it uses TCP port `33060` (not `3306`!).

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golistic/pxmysql/register"
	
	// or import the following to use the driver name "mysql"
	// _ "github.com/golistic/pxmysql/register/mysql"
)

func main() {
	db, err := sql.Open("pxmysql", "scott:tiger@tcp(127.0.0.1:33060)/somedb?useTLS=true")
	if err != nil {
		log.Fatalln(err)
	}

	var n string
	if err := db.QueryRow("SELECT NOW()").Scan(&n); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server time: %s", n)
}
```

The `useTLS` option is required when you use a user set to use the authentication
method `caching_sha2_password`.

Features
--------

We implement the X Protocol coming with MySQL 8.0. Earlier versions of MySQL
are not supported.

The following is a to-do list which shows also which features we are
implementing. The order might vary, but it should give an idea where we are
heading (bit like a roadmap).

* [x] Set up project using Git with README and LICENSE
* [x] Generate Go code from MySQL's `.protoc` files.
* [x] Setup testing environment with Docker image of MySQL
* [x] Connect and authenticate
    - [x] Using PLAIN (implies TLS support)
    - [x] Using MYSQL41
    - [x] SHA256_MEMORY
    - [x] Schema when connecting
    - [x] Unix socket connection
* [x] Certificate Authority of the MySQL server
* [ ] Query using SQL statements
    - [x] Consume all rows
    - [x] Last inserted ID & Rows Affected
    - [x] Test MySQL numeric data types
    - [x] Test MySQL Date and Time data types
    - [x] Test MySQL String data types
    - [x] Test MySQL Decimal data type
    - [ ] MySQL to Go types
* [ ] Collections
    - [ ] Create & Ensure
    - [ ] Drop
* [ ] Use Prepared Statement
    - [x] Type `statement` (the conventional SQL PREPARE)
    - [ ] CRUD operations: types `INSERT`, `FIND`, and `DELETE`
* [ ] Public APIs and documentation
* [x] Add Go `sql` driver

Requirements
------------

* Go 1.21 or greater
* MySQL 8.0 or greater (built using 8.0.34)

### Dependencies

We strive to have as little dependencies as possible using only the Go standard
library, extension or experimental packages.

We do, however, generate Go code from MySQL Server Protocol Buffer files.


Usage
-----

### Without Using Go's SQL Driver

1. Create a `Connection` object which holds configuration with which sessions
   are opened.
2. Using the `Connection.NewSession()` method, we create or more `Session`
   objects which open the actual connection and authenticate using the MySQL
   Plugin.

The following example uses a user which has its authentication plugin set as
`mysql_native_password`:

```go
package main

import (
  "context"
  "fmt"
  "log"

  "github.com/golistic/pxmysql/xmysql"
)

func main() {
  pwd := "tiger"
  config := &xmysql.ConnectConfig{
    Address:    "127.0.0.1:33060", // default X Plugin port
    AuthMethod: xmysql.AuthMethodMySQL41,
    Username:   "scott",
    Password:   &pwd,
  }

  session, err := xmysql.CreateSession(context.Background(), config)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Session", session)
}
```

In the following snippet we select a string and a timestamp:

```go
package main

import (
  "time"

  "github.com/golistic/pxmysql/xmysql"
)

// result of query: "SELECT 'now', NOW()" 
func handleResult(res *xmysql.ResultResult) {
  for _, row := range res.Rows {
    _ = row[0].(string)
    _ = row[1].(time.Time)
  }
}

```

MySQL Types to Go
-----------------

The below table gives an overview which type assertion needs to be used to get
the Go variant of the MySQL value.

The above shows how for each row of the result, there is a slice of any-values
which need to be type asserted.

| MySQL Types                                      | Go Type            | .. can be NULL  |
|--------------------------------------------------|--------------------|-----------------|
| `BINARY/VARBINARY`                               | `[]byte`           | `null.Bytes`    |
| `BIT`                                            | `uint64`           | `null.Uint64`   |
| `BLOB/TINYBLOB/MEDIUMBLOB/LONGBLOB`              | `[]byte`           | `null.Bytes`    |
| `CHAR/VARCHAR`                                   | `string`           | `null.String`   |
| `DATETIME/TIMESTAMP`                             | `time.Time`        | `null.Time`     |
| `DATE`                                           | `time.Time`        | `null.Time`     |
| `DECIMAL/NUMERIC`                                | `*decimal.Decimal` | `null.Decimal`  |
| `DOUBLE`                                         | `float64`          | `null.Float64`  |
| `ENUM`                                           | `string`           | `null.Strings`  |
| `FLOAT`                                          | `float32`          | `null.Float32`  |
| `SET`                                            | `[]string`         | `null.Strings`  |
| `SIGNED TINYINT/SMALLINT/MEDIUMINT/INT/BIGINT`   | `int64`            | `null.Int64`    |
| `TEXT/TINYTEXT/MEDIUMTEXT/LONGTEXT`              | `string`           | `null.String`   |
| `TIME`                                           | `time`             | `null.Duration` |
| `UNSIGNED TINYINT/SMALLINT/MEDIUMINT/INT/BIGINT` | `uint64`           | `null.Uint64`   |
| `YEAR`                                           | `int`              | `null.Int64`    |

### MySQL DECIMAL type

The MySQL DECIMAL-type is decoded into `decimal.Decimal` which stores the
integral and fractional (scale) parts as `*big.Int`. This way we support for
example MySQL `DECIMAL(65,1)` or `DECIMAL(32,30)`.

The `decimal.Decimal` struct can only be used for storage and retrieval. It is
not possible to do calculations. Users that need to manipulate the values
should use different types retrieving integer part using `Integral()` and
the fraction using `Fractional()` method.

The above example might not be useful for to most, but it is possible. Users
that need to calculate further will need to get the integer and precision
parts out of the `decimal.Decimal` type (see `Integral()` and `Fractional()`
methods) (or do it in MySQL).

The string representation of `decimal.Decimal` will add zero-padding to the
fractional part. When MySQL returns, for example, `82.003400` then the zero
on the right are not trimmed.


Configuration
-------------

When creating a new connection using `pxmysql.NewConnection`, it is possible to
provide a `ConnectConfig` instance. If not provided (is nil), then the default
configuration is used.

The `ConnectConfig`-type has the following attributes:

* `Address`: the host and TCP port on which the MySQL X Plugin listens
  (default: `127.0.0.1:33060`)
* `UseTLS`: when true, switches to TLS when possible (default: `false`)
* `Username`: username used when authenticating (default: `root`)
* `Password`: password used when authenticating (default: ``, empty)
* `Schema`: schema (database) to use after authenticating (default: ``, empty)
* `AuthMethod`: a support authentication mechanism (default: `AUTO`)
* `TLSServerCACertPath`: path to the file containing the CA certificate used by
  the MySQL server. If not provided, the system's Certificate Authorities are
  used.
* `TimeZoneName`: set time location for decoding DATETIME and TIMESTAMP MySQL
  data types to Go `time.Time` (see [MySQL Manual to support this][2])
  (default: UTC)

### Driver name

We use the driver name "pxmysql", which needs to be used with Go's `sql.Open`.
To register the driver, you use an anonymous import as follows:

```go
package yourstuff

import (
    _ "github.com/golistic/pxmysql/register"
)
```

Note the extra sub-package, which is different from other drivers. We do like to
explicitly load, and not register whenever the driver is imported.

Some projects require the "mysql" name to be registered, using the driver name
as a SQL dialect. This can be achieved by using the sub-package `register/mysql`:

```go
package yourstuff

import (
    _ "github.com/golistic/pxmysql/register/mysql"
)
```

We do not default to "mysql" as driver name, so it is possible to use
other drivers using MySQL at the same time.

### Authentication methods

The following authentication methods are supported:

* `PLAIN` or `pxmysql.AuthMethodPlain`: can only be used when TLS is available
* `MYSQL41` or `pxmysql.AuthMethodMySQL41`: the older, but more compatible
  mechanism (uses SHA1)
* `SHA256_MEMORY` or `pxmysql.AuthMethodSHA256Memory`: the newer mechanism which
  uses SHA256 and is cached; caveat: user must first authenticate using `PLAIN`
  over TLS to get the password hash cached

By default `AUTO` or `pxmysql.AuthMethodAuto` which tries the above methods in
order they are mentioned.

For `SHA256_MEMORY` or the `caching_sha2_password` plugin, you need to use TLS
for the first time the user connects. After that, it is possible to use non-TLS
with either `AUTO` or `SHA256_MEMORY`.

MySQL Documentation
-------------------

The MySQL X Protocol is documented in the following locations (similar content,
and the combination is more useful):

* https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_xplugin.html
* https://dev.mysql.com/doc/internals/en/x-protocol.html
* https://dev.mysql.com/doc/refman/8.0/en/x-plugin.html

Development
-----------

### Compile Go code from MySQL X Plugin Protocol Buffer definitions

The package `pxmysql` needs Go-code compiled from the MySQL X Plugin Protocol
Buffer definitions. This requires the [protoc][1] compiler to be installed
first. After, all that needs to be done is run the following from within the
root of the repository:

    go run ./cmd/genprotobuf

The above application will download the necessary files from the MySQL Server
GitHub repository eliminating the need to have the MySQL sources locally.
The generated Go code is stored under `internal/xmysql`.

In the same folder, an `info.md` containing the MySQL version that was used
to generate the code, and the timestamp when the command was run.

### Keep Collations Up-To-Date

We include the collation information of MySQL 8.0 by generating a Go source file
called `collations_data.go`. Note that only the `utf8mb3` collations are
considered since the X Protocol only supports this character set.

Use the following command line to (re)generate the file:

```shell
go run ./cmd/gencollations -address localhost:33060
```

See `go run ./cmd/gencollations -help` for more information.

### Run tests

Tests use a MySQL server running within a Docker container. It can be started
using the Docker Compose configuration found under `_support/mysqld`:

    cd _support/pxmysql-compose
    docker compose up -d

The above uses the Docker `compose` plugin. Alternatively, use `docker-compose`.

Run tests using:

    go test ./...

### Environment variables

Environment variables which could be useful when developing and debugging:

* `PXMYSQL_TRACE`: when set, messages read and written will be printed to STDERR
* `PXMYSQL_TRACE_VALUES`: when set, also dumps values of columns for each rows
  (does nothing when PXMYSQL_TRACE is not set)

Tests can be configured using the following environment variables:

* `PXMYSQL_TEST_DOCKER_CONTAINER`: specifies the name of the Docker container to
  use for running tests using MySQL (default `pxmysql.test.db`)
* `PXMYSQL_TEST_DOCKER_CONTAINER_GO`: specifies the name of the Docker container
  which is used to build Go applications (default `pxmysql.test.go`)
* `PXMYSQL_TEST_DOCKER_XPLUGIN_ADDR`: is the X Plugin address the container
  exposes (default `127.0.0.1:53360`)
* `PXMYSQL_TEST_DOCKER_MYSQL_ADDR`: is the conventional/classic MySQL address;
  not using X Protocol (default `127.0.01:53306`)
* `PXMYSQL_TEST_DOCKER_MYSQL_PWD`: password of the root-user, set when starting
  the MySQL container (default `rootpwd`)
* `PXMYSQL_TEST_DOCKER_EXEC`: path of the `docker` executable

About The Author
----------------

Geert Vanderkelen worked for about 13 years at MySQL AB/MySQL Inc/Sun/Oracle
as Support Engineer and Developer. He is the original author of MySQL Connector/Python
(which started as a hobby project), and MySQL Router. Today, Geert implements
backend services and more using primarily Go/GraphQL, and advocates the
goodness of MySQL and its relational kin.

License
-------

Distributed under the MIT license. See `LICENSE.md` for more information.

[1]: https://grpc.io/docs/languages/go/quickstart/

[2]: https://dev.mysql.com/doc/refman/8.0/en/time-zone-support.html

[3]: https://github.com/go-sql-driver/mysql
