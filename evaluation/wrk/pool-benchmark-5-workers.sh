#!/bin/bash

# set -e
echo "benchmark 5 workers"
LOGFILE="/home/bohdan/workspace/uni/thesis/execution-system-playground/log.txt"
BIN="/home/bohdan/workspace/uni/thesis/execution-system-playground/main"

"$BIN" --path /home/bohdan/workspace/uni/thesis/execution-system-playground/evaluation/wrk/configs/pool-baseline-5.json >> "$LOGFILE"  &
sleep 2

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

