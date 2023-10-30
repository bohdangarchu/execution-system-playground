#!/bin/bash

set -e
# run with sudo
LOGFILE="/home/bohdan/workspace/uni/thesis/execution-system-playground/log.txt"
BIN="/home/bohdan/workspace/uni/thesis/execution-system-playground/main"
echo "benchmark single worker"
echo "docker benchmark"
echo "DOCKER BENCHMARK" >> "$LOGFILE"

"$BIN" --path /home/bohdan/workspace/uni/thesis/execution-system-playground/evaluation/wrk/configs/pool-docker-1.json >> "$LOGFILE"  &

sleep 5

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

sleep 3

echo "firecracker benchmark"
echo "FIRECRACKER BENCHMARK" >> "$LOGFILE"

"$BIN" --path /home/bohdan/workspace/uni/thesis/execution-system-playground/evaluation/wrk/configs/pool-firecracker-1.json >> "$LOGFILE" &

sleep 10

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

sleep 3

echo "baseline benchmark"
echo "BASELINE BENCHMARK" >> "$LOGFILE"

"$BIN" --path /home/bohdan/workspace/uni/thesis/execution-system-playground/evaluation/wrk/configs/pool-baseline-1.json >> "$LOGFILE" &

sleep 5

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill