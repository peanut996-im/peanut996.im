#!/bin/bash
# author: peanut996
# date: 2021.4.26
# description: 更新远程可执行文件


git submodule update --init --recursive --remote

git add .

git commit -m"update and deploy"

git push origin master