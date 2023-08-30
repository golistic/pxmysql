module example.com/unix_socket

go 1.21

replace github.com/golistic/pxmysql => /go/src/github.com/golistic/pxmysql

require github.com/golistic/pxmysql v1.0.0

require (
	github.com/golistic/xgo v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20230817173708-d852ddb80c63 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
