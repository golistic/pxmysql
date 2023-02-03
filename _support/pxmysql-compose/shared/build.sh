#!/bin/sh

out="/shared/builds/$1"
main="/shared/goapps/$1"

export GOPRIVATE="github.com/golistic/pxmysql"

cd "${main}" || exit 1
go mod tidy
go build -o "${out}" .
