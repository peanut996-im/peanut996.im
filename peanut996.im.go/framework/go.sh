#!/bin/bash
# author: peanut996
# date: 2021.3.31
# description: 一键运行项目

appName=${PWD##*/}
targetos=`uname | tr "[A-Z]" "[a-z]"`
if [[ $targetos == "mingw"* ]];then
targetos="windows"
fi
bash ./build.sh $targetos

echo "./bin/$appName -f ./etc/config.yaml $1 $2 $3 $4"
echo "run $appName..."
echo 
./bin/$appName -f ./etc/config.yaml $1 $2 $3 $4