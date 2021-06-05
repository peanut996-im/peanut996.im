#!/bin/bash
# author: peanut996
# date: 2021.5.22
# description: 一键运行项目

powershell -File stop.ps1&&rm sso/log/sso.log;rm logic/log/logic.log;rm gate/log/gate.log
cd sso && bash go.sh &
sleep 30s && cd logic && bash go.sh &
sleep 60s && cd gate && bash go.sh
