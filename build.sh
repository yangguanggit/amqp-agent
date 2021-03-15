#! /bin/bash

#env: debug test release
#cate: web command
#eg: bash build.sh release web main
#eg: bash build.sh release command main

env=$1
if [ -z "${env}" ]; then
    env=release
fi

cate=$2
if [ -z "${cate}" ]; then
    cate=web
fi

output=$3
if [ -z ${output} ]; then
    output=main
fi

if [ "$env" = debug ]; then
    go build -ldflags "-s -w" -o ${cate}/${output} --tags "${env}" ${cate}/main.go
else
    CGO_ENABLED=0 go build -ldflags "-s -w" -o ${cate}/${output} --tags "${env}" ${cate}/main.go
fi
