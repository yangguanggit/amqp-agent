#! /bin/bash

#eg: bash reload.sh main

command=$1
if [ -z ${command} ]; then
    command=main
fi

function Echo() {
    local str=$1
    local color=$2
    if [ -z ${color} ]; then color=green; fi

    case ${color} in
        red) echo -e "\033[31m ${str} \033[0m" ;;
        green) echo -e "\033[32m ${str} \033[0m" ;;
        white) echo -e "\033[37m ${str} \033[0m" ;;
        yellow) echo -e "\033[33m ${str} \033[0m" ;;
        *) echo ${str} ;;
    esac
}

cate=web
port=80

pid=`/usr/sbin/lsof -i tcp:${port} | grep ${command} | awk '{print $2}'`
if [ -z ${pid} ]; then
    Echo "Noting listen the port ${port}."
    exit
fi

for i in `ls -l ${cate} | sed -n '2,$ p' | awk '{print $9}'`
do
    if [ ${i} == ${command} ]; then
        kill -USR2 ${pid}
        Echo "The ${command} listen port ${port} again."
        exit
    fi
done

Echo "The current directory no such file ${command}." red
