#!/bin/bash
# author: peanut996
# date: 2021.4.25
# description: 一键运行项目

appName=${PWD##*/}
targetos=`uname | tr "[A-Z]" "[a-z]"`
if [[ $targetos == "mingw"* ]];then
targetos="windows"
fi


bash ./build.sh $targetos

echo "run $appName: ./bin/$appName -c ./etc/config.yaml $1 $2 $3 $4"
echo ""
./bin/$appName -c ./etc/config.yaml $1 $2 $3 $4