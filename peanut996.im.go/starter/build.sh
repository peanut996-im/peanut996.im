#!/bin/bash
# author: peanut996
# date: 2021.3.31
# description: 编译项目

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# 编译选项
# windows
# linux
# darwin
if [ ! -n "$1" ] ;then
    echo "you need input target os { windows | linux | darwin }. -'darwin' is mac os"
    exit
else
    echo "target os: $1"
    echo
fi
targetos=$1
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

bash ./mod.sh

appName=${PWD##*/}

rm -f ./bin/*${appName}*
cd ./src

## 编译参数
go build -o ../bin/${appName} .

if [ ${targetos} = "windows" ];then
    cd ./bin
    mv ${appName} ${appName}.exe
    cd ../
fi