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
BuildTime=$(date)
BuildTime=${BuildTime// /_}
BuildUser=$(whoami)
BuildUser=${BuildUser// /_}
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_VERSION=$(git log --pretty=format:"%h" -1)
BuildVersion="${GIT_BRANCH}_${GIT_VERSION}"
BuildTimeStamp=`date +%s`
check_ifconfig=``
if [ "${targetos}"=="windows" ]; then
    BuildMachine=$(ipconfig |grep "IPv4" |grep -v "192.168.2.1"| grep -v "127.0.0.1"| grep -v "localhost" | awk -F': ' '{print $2}')
elif [ "$(ls /sbin/ | grep ifconfig)"/sbin/ifconfig | grep "inet" | grep -v "127.0.0.1" | grep -v "inet6" | awk '{print $2}'| tr "\n" " " = "ifconfig" ] ;then
    BuildMachine=$(/sbin/ifconfig | grep "inet" | grep -v "127.0.0.1" | grep -v "inet6" | awk '{print $2}'| tr "\n" " ")
else
    BuildMachine=$(ip addr | grep "inet" | grep -v "127.0.0.1" | grep -v "inet6" | awk '{print $2}' | awk -F'/' '{print $1}' | head -n1)
fi


BuildMachine=${BuildMachine// /_}
BuildMachine="$(uname -n)@${BuildMachine}"
rm -f ./bin/*${appName}*
cd ./src

## 编译参数
go build  -ldflags "-X main.BuildVersion=${BuildVersion} -X main.BuildUser=${BuildUser} -X main.BuildTime=${BuildTime} -X main.BuildMachine=${BuildMachine} " -o ../bin/${appName} .

if [ ${targetos} = "windows" ];then
    cd ../bin
    mv ${appName} ${appName}.exe
    cd ../
fi