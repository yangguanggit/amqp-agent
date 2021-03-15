#! /bin/bash

#env: debug test release
#eg: bash strat.sh release main

env=$1
if [ -z "${env}" ]; then
    env=release
fi

command=$2
if [ -z ${command} ]; then
    command=main
fi

web/${command} --env=${env} >> ./runtime/stdout.log 2>&1 &
