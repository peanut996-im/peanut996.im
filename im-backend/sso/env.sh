#!/bin/bash
# author: peanut996
# date: 2021.4.25
# des: 自动构建

olddir=${PWD}
repository=${PWD##*/}
echo "env setting"

packages=("framework")

for package in ${packages[@]}; do
        cd ../

        checkdir=${PWD}"/"${package}
        echo "package -> "${checkdir}
        if [ -d "$checkdir" ] ;then
            echo "Got it!"
            cd ./${package}
            bash mod.sh
            cd ../ 
        fi

        cd $olddir
done


echo "env setting over"
echo